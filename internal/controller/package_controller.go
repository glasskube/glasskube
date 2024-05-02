/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	ctrladapter "github.com/glasskube/glasskube/internal/adapter/controllerruntime"
	"github.com/glasskube/glasskube/internal/controller/conditions"
	"github.com/glasskube/glasskube/internal/controller/labels"
	"github.com/glasskube/glasskube/internal/controller/owners"
	ownerutils "github.com/glasskube/glasskube/internal/controller/owners/utils"
	"github.com/glasskube/glasskube/internal/controller/requeue"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/names"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/pkg/condition"
	"go.uber.org/multierr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// PackageReconciler reconciles a Package object
type PackageReconciler struct {
	client.Client
	record.EventRecorder
	*owners.OwnerManager
	RepoClientset     repoclient.RepoClientset
	ValueResolver     *manifestvalues.Resolver
	Scheme            *runtime.Scheme
	ManifestAdapter   manifest.ManifestAdapter
	HelmAdapter       manifest.ManifestAdapter
	KustomizeAdapter  manifest.ManifestAdapter
	DependencyManager *dependency.DependendcyManager
}

//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packages/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update
//+kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases,verbs=get;list;watch;create;update
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=helmrepositories,verbs=get;list;watch;create;update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Package object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PackageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pkg packagesv1alpha1.Package

	if err := r.Get(ctx, req.NamespacedName, &pkg); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !pkg.DeletionTimestamp.IsZero() {
		changed, err := conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions,
			condition.Pending, "Package is being deleted")
		if changed {
			telemetry.ForOperator().ReportDelete(&pkg)
		}
		return ctrl.Result{}, err
	}

	telemetry.ForOperator().ReconcilePackage(&pkg)

	return r.reconcilePackage(ctx, pkg)
}

func (r *PackageReconciler) reconcilePackage(ctx context.Context, pkg v1alpha1.Package) (ctrl.Result, error) {
	return (&PackageReconcilationContext{PackageReconciler: r, pkg: &pkg}).reconcile(ctx)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PackageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.OwnerManager == nil {
		r.OwnerManager = owners.NewOwnerManager(r.Scheme)
	}
	if r.ValueResolver == nil {
		r.ValueResolver = manifestvalues.NewResolver(
			ctrladapter.NewPackageClientAdapter(r.Client),
			ctrladapter.NewKubernetesClientAdapter(r.Client),
		)
	}
	controllerBuilder := ctrl.NewControllerManagedBy(mgr).For(&packagesv1alpha1.Package{})
	for _, adapter := range []manifest.ManifestAdapter{r.HelmAdapter, r.KustomizeAdapter, r.ManifestAdapter} {
		if adapter != nil {
			if err := adapter.ControllerInit(controllerBuilder, r.Client, r.Scheme); err != nil {
				return err
			}
		}
	}
	return controllerBuilder.
		Owns(&packagesv1alpha1.PackageInfo{}, builder.MatchEveryOwner).
		Owns(&packagesv1alpha1.Package{}, builder.MatchEveryOwner).
		Complete(r)
}

type PackageReconcilationContext struct {
	*PackageReconciler
	pkg                   *v1alpha1.Package
	pi                    *v1alpha1.PackageInfo
	isSuccess             bool
	shouldUpdateStatus    bool
	currentOwnedResources []v1alpha1.OwnedResourceRef
	currentOwnedPackages  []v1alpha1.OwnedResourceRef
}

func (r *PackageReconcilationContext) setShouldUpdate(value bool) {
	r.shouldUpdateStatus = r.shouldUpdateStatus || value
}

func (r *PackageReconcilationContext) reconcile(ctx context.Context) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	if err := r.ensurePackageInfo(ctx); err != nil {
		return requeue.Always(ctx, err)
	}

	if meta.IsStatusConditionTrue(r.pi.Status.Conditions, string(condition.Ready)) {
		log.V(1).Info("PackageInfo is ready", "packageInfo", r.pi.Name)
		return r.reconcilePackageInfoReady(ctx)
	} else if meta.IsStatusConditionFalse(r.pi.Status.Conditions, string(condition.Ready)) {
		packageInfoCondition := meta.FindStatusCondition(r.pi.Status.Conditions, string(condition.Ready))
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.Reason(packageInfoCondition.Reason), packageInfoCondition.Message))
		return r.finalize(ctx)
	} else {
		r.setShouldUpdate(
			conditions.SetUnknown(ctx, &r.pkg.Status.Conditions, condition.Pending, "PackageInfo status is unknown"))
		return r.finalize(ctx)
	}
}

func (r *PackageReconcilationContext) reconcilePackageInfoReady(ctx context.Context) (ctrl.Result, error) {
	piManifest := r.pi.Status.Manifest
	if piManifest == nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.UnsupportedFormat, "manifest must not be nil"))
		return r.finalizeNoRequeue(ctx)
	}

	if !r.ensureDependencies(ctx) {
		return r.finalize(ctx)
	}

	var patches []manifestvalues.TargetPatch
	if resolvedValues, err := r.ValueResolver.Resolve(ctx, r.pkg.Spec.Values); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.ValueConfigurationInvalid, err.Error()))
		return r.finalizeWithError(ctx, err)
	} else if err := manifestvalues.ValidateResolvedValues(*piManifest, resolvedValues); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.ValueConfigurationInvalid, err.Error()))
		return r.finalizeWithError(ctx, err)
	} else if p, err := manifestvalues.GeneratePatches(*piManifest, resolvedValues); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.InstallationFailed, err.Error()))
		return r.finalizeWithError(ctx, err)
	} else {
		patches = p
	}

	// First, collect the adapters for all included manifests and ensure that they are supported.
	// If one manifest type is not supported, no action must be performed!

	var adaptersToRun []manifest.ManifestAdapter
	if len(piManifest.Manifests) > 0 {
		if r.ManifestAdapter == nil {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
					condition.UnsupportedFormat, "manifests not supported"))
			return r.finalizeNoRequeue(ctx)
		}
		adaptersToRun = append(adaptersToRun, r.ManifestAdapter)
	}
	if piManifest.Kustomize != nil {
		if r.KustomizeAdapter == nil {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
					condition.UnsupportedFormat, "kustomize not supported"))
			return r.finalizeNoRequeue(ctx)
		}
		adaptersToRun = append(adaptersToRun, r.KustomizeAdapter)
	}
	if piManifest.Helm != nil {
		if r.HelmAdapter == nil {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
					condition.UnsupportedFormat, "helm not supported"))
			return r.finalizeNoRequeue(ctx)
		}
		adaptersToRun = append(adaptersToRun, r.HelmAdapter)
	}

	results := make([]result.ReconcileResult, 0, len(adaptersToRun))
	var errs error
	for _, adapter := range adaptersToRun {
		if result, err := adapter.Reconcile(ctx, r.pkg, piManifest, patches); err != nil {
			errs = multierr.Append(errs, err)
		} else {
			results = append(results, *result)
			ownerutils.Add(&r.currentOwnedResources, result.OwnedResources...)
		}
	}

	if errs != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.InstallationFailed, errs.Error()))
		return r.finalizeWithError(ctx, errs)
	} else if !r.handleAdapterResults(ctx, results) {
		return r.finalize(ctx)
	} else {
		r.afterSuccess(ctx, results)
		return r.finalize(ctx)
	}
}

func (r *PackageReconcilationContext) ensureDependencies(ctx context.Context) bool {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("ensuring dependencies", "dependencies", r.pi.Status.Manifest.Dependencies)

	var failed []string
	if result, err := r.DependencyManager.Validate(ctx, r.pi.Status.Manifest, r.pkg.Spec.PackageInfo.Version); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.InstallationFailed, fmt.Sprintf("error validating dependencies: %v", err)))
		return false
	} else if result.Status == dependency.ValidationResultStatusResolvable {
		for _, requirement := range result.Requirements {
			if requirement.Transitive {
				// Only direct dependencies should be touched in the context of the reconciliation of a package.
				continue
			}
			newPkg := &packagesv1alpha1.Package{
				ObjectMeta: metav1.ObjectMeta{
					Name:      requirement.Name,
					Namespace: r.pkg.Namespace,
				},
			}

			repositories, err := r.RepoClientset.Aggregate().GetReposForPackage(requirement.Name)
			if err != nil {
				log.Error(err, "could not find repos for package", "required", requirement.Name)
				failed = append(failed, requirement.Name)
				continue
			}

			var repositoryName string
			switch len(repositories) {
			case 0:
				log.Error(err, "could not find package in any repo", "required", requirement.Name)
				failed = append(failed, requirement.Name)
				continue
			case 1:
				repositoryName = repositories[0].Name
			default:
				log.Error(err, "dependency in multiple repos is not supported yet", "required", requirement.Name)
				failed = append(failed, requirement.Name)
				continue
			}

			if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, newPkg, func() error {
				if err := r.SetOwner(r.pkg, newPkg, owners.BlockOwnerDeletion); err != nil {
					return fmt.Errorf("unable to set ownerReference on required package: %w", err)
				}
				newPkg.Spec = packagesv1alpha1.PackageSpec{
					PackageInfo: packagesv1alpha1.PackageInfoTemplate{
						Name:           requirement.Name,
						Version:        requirement.Version,
						RepositoryName: repositoryName,
					},
				}
				return nil
			}); err != nil {
				log.Error(err, "Failed to create required package", "required", requirement.Name)
				failed = append(failed, requirement.Name)
			}
		}
		if len(failed) > 0 {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
					condition.InstallationFailed, fmt.Sprintf("required package(s) not installed: %v", strings.Join(failed, ","))))
			return false
		}
	} else if result.Status == dependency.ValidationResultStatusConflict {
		var parts []string
		for _, c := range result.Conflicts {
			parts = append(parts, fmt.Sprintf("need version %v of %v but found %v", c.Required.Version, c.Actual.Name, c.Actual.Version))
		}
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.InstallationFailed, fmt.Sprintf("conflicting dependencies: %v", strings.Join(parts, ","))))
		return false
	}

	var ownedPackages []packagesv1alpha1.OwnedResourceRef
	var waitingFor []string
	// if all requirements fulfilled, status can be checked
	for _, dep := range r.pi.Status.Manifest.Dependencies {
		var requiredPkg packagesv1alpha1.Package
		if err := r.Get(ctx, types.NamespacedName{
			Namespace: r.pkg.Namespace,
			Name:      dep.Name,
		}, &requiredPkg); err != nil {
			message := fmt.Sprintf("failed to get required package %v: %v", dep.Name, err)
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions, condition.InstallationFailed, message))
			return false
		} else {
			// if the required package already exists, we set the owner reference if
			// * there already exists another package owner reference (i.e. it has been installed as dependency of another package)
			// * but NOT if there exists no other package owner reference (i.e. it has been installed manually)
			if ok, err := r.HasAnyOwnerOfType(r.pkg, &requiredPkg); err != nil || ok {
				if err != nil {
					log.Error(err, "Failed to check for owner references", "owner", r.pkg.Name, "owned", requiredPkg.Name)
				}
				if ok, err := r.HasOwner(r.pkg, &requiredPkg); err != nil || !ok {
					log.Info("Updating existing required package with new owner reference", "owner", r.pkg.Name, "owned", requiredPkg.Name)
					if err := r.SetOwner(r.pkg, &requiredPkg, owners.BlockOwnerDeletion); err != nil {
						log.Error(err, "Failed to set owner reference", "owner", r.pkg.Name, "owned", requiredPkg.Name)
						failed = append(failed, dep.Name)
					}
					if err := r.Update(ctx, &requiredPkg); err != nil {
						log.Error(err, "Failed to updated required package", "owner", r.pkg.Name, "owned", requiredPkg.Name)
						failed = append(failed, dep.Name)
					}
				}
			}

			if owned, err := ownerutils.ToOwnedResourceRef(r.Scheme, &requiredPkg); err != nil {
				log.Error(err, "Failed to create OwnedResourceRef", "package", requiredPkg)
				failed = append(failed, dep.Name)
			} else {
				ownedPackages = append(ownedPackages, owned)
			}
			if meta.IsStatusConditionTrue(requiredPkg.Status.Conditions, string(condition.Failed)) {
				failed = append(failed, requiredPkg.Name)
			} else if !meta.IsStatusConditionTrue(requiredPkg.Status.Conditions, string(condition.Ready)) {
				waitingFor = append(waitingFor, requiredPkg.Name)
			}
		}
	}
	if len(failed) > 0 {
		message := fmt.Sprintf("required package(s) not installed: %v", strings.Join(failed, ","))
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions, condition.InstallationFailed, message))
		return false
	}
	if len(waitingFor) > 0 {
		message := fmt.Sprintf("waiting for required package(s) %v", strings.Join(waitingFor, ","))
		r.setShouldUpdate(
			conditions.SetUnknown(ctx, &r.pkg.Status.Conditions, condition.Pending, message))
		return false
	}

	ownerutils.Add(&r.currentOwnedPackages, ownedPackages...)

	return true
}

func (r *PackageReconcilationContext) ensurePackageInfo(ctx context.Context) error {
	packageInfo := packagesv1alpha1.PackageInfo{
		ObjectMeta: metav1.ObjectMeta{Name: names.PackageInfoName(*r.pkg)},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("PackageInfo", packageInfo.Name)

	log.V(1).Info("ensuring PackageInfo")
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, &packageInfo, func() error {
		if err := r.SetOwner(r.pkg, &packageInfo, owners.BlockOwnerDeletion); err != nil {
			return fmt.Errorf("unable to set ownerReference on PackageInfo: %w", err)
		}
		packageInfo.Spec = packagesv1alpha1.PackageInfoSpec{
			Name:           r.pkg.Spec.PackageInfo.Name,
			Version:        r.pkg.Spec.PackageInfo.Version,
			RepositoryName: r.pkg.Spec.PackageInfo.RepositoryName,
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not create or update PackageInfo: %w", err)
	}

	log.V(1).Info("ensured PackageInfo", "result", result)

	// After CreateOrUpdate, PackageInfo does not always have correct TypeMeta.
	// To fix this, we Get it again here.
	// TODO: Find out why this is even needed.
	if result == controllerutil.OperationResultCreated {
		err := retry.OnError(retry.DefaultBackoff, apierrors.IsNotFound, func() error {
			return r.Get(ctx, client.ObjectKeyFromObject(&packageInfo), &packageInfo)
		})
		if err != nil {
			return err
		}
	}

	if changed, err := ownerutils.AddOwnedResourceRef(r.Scheme, &r.pkg.Status.OwnedPackageInfos, &packageInfo); err != nil {
		log.Error(err, "could not add PackageInfo to owned resources")
	} else {
		r.setShouldUpdate(changed)
	}

	r.pi = &packageInfo
	return nil
}

func (r *PackageReconcilationContext) handleAdapterResults(ctx context.Context, results []result.ReconcileResult) bool {
	var firstFailed *result.ReconcileResult
	var firstWaiting *result.ReconcileResult
	for i, result := range results {
		if result.IsFailed() && firstFailed == nil {
			firstFailed = &results[i]
		} else if result.IsWaiting() && firstWaiting == nil {
			firstWaiting = &results[i]
		}
	}

	if firstFailed != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
				condition.InstallationFailed, firstFailed.Message))
		return false
	} else if firstWaiting != nil {
		r.setShouldUpdate(
			conditions.SetUnknown(ctx, &r.pkg.Status.Conditions, condition.Pending, firstWaiting.Message))
		return false
	} else {
		return true
	}
}

func (r *PackageReconcilationContext) afterSuccess(ctx context.Context, results []result.ReconcileResult) {
	reason := condition.UpToDate
	message := "PackageInfo has nothing to apply (no helm or kustomize manifest present)"
	if len(results) > 0 {
		messages := make([]string, 0, len(results))
		for _, result := range results {
			messages = append(messages, result.Message)
		}
		reason = condition.InstallationSucceeded
		message = strings.Join(messages, "\n")
	}

	r.setShouldUpdate(
		conditions.SetReady(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions, reason, message))
	r.setShouldUpdate(r.pkg.Status.Version != r.pi.Status.Version)
	r.pkg.Status.Version = r.pi.Status.Version
	r.isSuccess = true
}

func (r *PackageReconcilationContext) finalize(ctx context.Context) (ctrl.Result, error) {
	return requeue.Always(ctx, r.actualFinalize(ctx))
}

func (r *PackageReconcilationContext) finalizeNoRequeue(ctx context.Context) (ctrl.Result, error) {
	return requeue.OnError(ctx, r.actualFinalize(ctx))
}

func (r *PackageReconcilationContext) finalizeWithError(ctx context.Context, err error) (ctrl.Result, error) {
	return requeue.Always(ctx, multierr.Append(err, r.actualFinalize(ctx)))
}

func (r *PackageReconcilationContext) actualFinalize(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)

	r.setShouldUpdate(ownerutils.Add(&r.pkg.Status.OwnedResources, r.currentOwnedResources...))
	r.setShouldUpdate(ownerutils.Add(&r.pkg.Status.OwnedPackages, r.currentOwnedPackages...))

	var errs error
	if r.isSuccess {
		if err := r.cleanup(ctx); err != nil {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.Status.Conditions,
					condition.InstallationFailed, fmt.Sprintf("cleanup failed: %v", err)))
			errs = multierr.Append(errs, err)
		}
		log.V(1).Info("cleanup done")
	}

	if r.shouldUpdateStatus {
		err := r.Status().Update(ctx, r.pkg)
		if err != nil {
			log.Error(err, "package status update failed")
		} else {
			log.Info("package status updated")
		}
		errs = multierr.Append(errs, err)
	}
	return errs
}

func (r *PackageReconcilationContext) cleanup(ctx context.Context) error {
	return multierr.Combine(
		r.pruneOwnedResources(ctx),
		r.pruneOwnedPackageInfos(ctx),
		r.pruneOwnedPackages(ctx),
	)
}

func (r *PackageReconcilationContext) pruneOwnedResources(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	var errs error
	ownedResourcesCopy := r.pkg.Status.OwnedResources[:]
OuterLoop:
	for _, ref := range r.pkg.Status.OwnedResources {
		for _, newRef := range r.currentOwnedResources {
			if ownerutils.RefersToSameResource(ref, newRef) {
				// ref is still an owned resource
				continue OuterLoop
			}
		}
		// ref is no longer an owned resource.
		// check if it is managed by the operator and delete it if it is.
		obj := ownerutils.OwnedResourceRefToObject(ref)
		if err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
			if !apierrors.IsNotFound(err) {
				multierr.AppendInto(&errs, fmt.Errorf("could get resource during pruning: %w", err))
			}
		} else if labels.IsManaged(obj) {
			if err := r.Delete(ctx, obj); err != nil && !apierrors.IsNotFound(err) {
				multierr.AppendInto(&errs, fmt.Errorf("could not prune resource: %w", err))
			} else {
				log.V(1).Info("pruned resource", "reference", ref)
				r.setShouldUpdate(ownerutils.Remove(&ownedResourcesCopy, ref))
			}
		} else {
			log.V(1).Info("skipped pruning unmanaged resource", "reference", ref)
			r.setShouldUpdate(ownerutils.Remove(&ownedResourcesCopy, ref))
		}
	}
	r.pkg.Status.OwnedResources = ownedResourcesCopy
	return errs
}

func (r *PackageReconcilationContext) pruneOwnedPackageInfos(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	currentRef, err := ownerutils.ToOwnedResourceRef(r.Scheme, r.pi)
	if err != nil {
		return err
	}
	var compositeErr error
	for _, ref := range r.pkg.Status.OwnedPackageInfos {
		if !ownerutils.RefersToSameResource(ref, currentRef) {
			var packageInfo packagesv1alpha1.PackageInfo
			key := client.ObjectKeyFromObject(ownerutils.OwnedResourceRefToObject(ref))
			if err := r.Get(ctx, key, &packageInfo); apierrors.IsNotFound(err) {
				log.Info("PackageInfo not found", "PackageInfo", ref.Name)
			} else if err != nil {
				compositeErr = multierr.Append(compositeErr, err)
				continue
			} else {
				if owned, err := ownerutils.ObjHasOwner(&packageInfo, r.pkg); err != nil {
					compositeErr = multierr.Append(compositeErr, err)
					continue
				} else if owned {
					// Remove the owner reference for pkg
					if err := controllerutil.RemoveOwnerReference(r.pkg, &packageInfo, r.Scheme); err != nil {
						compositeErr = multierr.Append(compositeErr, err)
						continue
					}
					if len(packageInfo.OwnerReferences) > 0 {
						// If other owner references remain, update the PackageInfo with the owner reference removed
						log.V(1).Info("updating old package info", "PackageInfo", packageInfo.Name)
						if err := r.Update(ctx, &packageInfo); client.IgnoreNotFound(err) != nil {
							compositeErr = multierr.Append(compositeErr, err)
							continue
						}
					} else {
						// If no other owner references remain, delete the PackageInfo
						log.V(1).Info("deleting old package info", "PackageInfo", packageInfo.Name)
						if err := r.Delete(ctx, &packageInfo); client.IgnoreNotFound(err) != nil {
							compositeErr = multierr.Append(compositeErr, err)
							continue
						}
					}
				}
			}

			// Remove the PackageInfo from the owned PackageInfos field of pkg
			ownerutils.RemoveOwnedResourceRef(&r.pkg.Status.OwnedPackageInfos, ref)
			r.setShouldUpdate(true)
		}
	}
	return compositeErr
}

func (r *PackageReconcilationContext) pruneOwnedPackages(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	var errs error
	ownedPackagesCopy := r.pkg.Status.OwnedPackages[:]
OuterLoop:
	for _, ref := range r.pkg.Status.OwnedPackages {
		for _, newRef := range r.currentOwnedPackages {
			if ownerutils.RefersToSameResource(ref, newRef) {
				continue OuterLoop
			}
		}

		var oldReqPkg v1alpha1.Package
		if err := r.Get(ctx, types.NamespacedName{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		}, &oldReqPkg); err != nil {
			if apierrors.IsNotFound(err) {
				r.setShouldUpdate(ownerutils.Remove(&ownedPackagesCopy, ref))
			} else {
				log.Error(err, "Failed to get old required package", "oldPackage", ref.Name)
			}
		} else if owning, err := r.HasOwner(r.pkg, &oldReqPkg); err != nil {
			log.Error(err, "Failed to check owner references of old required package", "oldPackage", ref.Name)
		} else if owning {
			if ref.MarkedForDeletion {
				log.Info(fmt.Sprintf("reference from %v to %v marked for deletion will be removed", r.pkg.Name, ref.Name))
				var cnt int
				cnt, err = r.CountOwnersOfType(r.pkg, &oldReqPkg)
				if err != nil {
					log.Error(err, "Failed to check owner references of old required package", "oldPackage", ref.Name)
					continue
				}

				if err := r.RemoveOwner(r.pkg, &oldReqPkg); err != nil {
					log.Error(err, "Failed to remove owner reference", "oldPackage", ref.Name)
				} else if err := r.Update(ctx, &oldReqPkg); err != nil {
					log.Error(err, "Failed to update old package with removed owner reference", "oldPackage", ref.Name)
				} else {
					log.Info(fmt.Sprintf("Removed owner reference from %v to %v", r.pkg.Name, ref.Name))

					if cnt == 1 {
						// remove the old package if we were the only package owning it
						deletePropagationForeground := metav1.DeletePropagationForeground
						if err := r.Delete(ctx, ownerutils.OwnedResourceRefToObject(ref), &client.DeleteOptions{
							PropagationPolicy: &deletePropagationForeground,
						}); err != nil && !apierrors.IsNotFound(err) {
							errs = multierr.Append(errs, fmt.Errorf("could not prune package: %w", err))
						} else {
							log.V(1).Info("pruned package", "reference", ref)
						}
					}
				}
			} else {
				log.Info(fmt.Sprintf("marking for deletion: reference from %v to %v", r.pkg.Name, ref.Name))
				r.setShouldUpdate(ownerutils.MarkForDeletion(&ownedPackagesCopy, ref))
			}
		} else {
			r.setShouldUpdate(ownerutils.Remove(&ownedPackagesCopy, ref))
		}
	}
	r.pkg.Status.OwnedPackages = ownedPackagesCopy
	return errs
}

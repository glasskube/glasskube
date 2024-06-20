package controller

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	ctrladapter "github.com/glasskube/glasskube/internal/adapter/controllerruntime"
	"github.com/glasskube/glasskube/internal/controller/conditions"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/controller/labels"
	"github.com/glasskube/glasskube/internal/controller/owners"
	ownerutils "github.com/glasskube/glasskube/internal/controller/owners/utils"
	"github.com/glasskube/glasskube/internal/controller/requeue"
	"github.com/glasskube/glasskube/internal/controller/watch"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/names"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/util"
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

const (
	packageDeletionFinalizer = "packages.glasskube.dev/packageDeletion"
)

type PackageReconcilerCommon struct {
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

func (r *PackageReconcilerCommon) baseSetup(
	mgr ctrl.Manager, object ctrlpkg.Package, lister watch.PackageLister) (*builder.Builder, error) {

	if r.OwnerManager == nil {
		r.OwnerManager = owners.NewOwnerManager(r.Scheme)
	}
	if r.ValueResolver == nil {
		r.ValueResolver = manifestvalues.NewResolver(
			ctrladapter.NewPackageClientAdapter(r.Client),
			ctrladapter.NewKubernetesClientAdapter(r.Client),
		)
	}

	controllerBuilder := ctrl.NewControllerManagedBy(mgr).
		For(object).
		Watches(&v1alpha1.PackageInfo{},
			watch.EnqueueRequestsFromOwnedResource(r.Scheme, lister, watch.OwnedPackageInfos)).
		Watches(&v1alpha1.ClusterPackage{},
			watch.EnqueueRequestsFromOwnedResource(r.Scheme, lister, watch.OwnedPackages))

	if err := r.InitAdapters(controllerBuilder); err != nil {
		return nil, err
	}
	return controllerBuilder, nil
}

func (r *PackageReconcilerCommon) InitAdapters(builder *builder.Builder) error {
	for _, adapter := range []manifest.ManifestAdapter{r.HelmAdapter, r.KustomizeAdapter, r.ManifestAdapter} {
		if adapter != nil {
			if err := adapter.ControllerInit(builder, r.Client, r.Scheme); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *PackageReconcilerCommon) reconcile(ctx context.Context, pkg ctrlpkg.Package) (ctrl.Result, error) {
	prc := &PackageReconcilationContext{PackageReconcilerCommon: r, pkg: pkg}
	log := ctrl.LoggerFrom(ctx)

	if !pkg.GetDeletionTimestamp().IsZero() {
		return prc.reconcileAfterDeletion(ctx)
	}

	telemetry.ForOperator().ReconcilePackage(prc.pkg)
	prc.ensureFinalizer()

	if err := prc.ensurePackageInfo(ctx); err != nil {
		return requeue.Always(ctx, err)
	}

	if meta.IsStatusConditionTrue(prc.pi.Status.Conditions, string(condition.Ready)) {
		log.V(1).Info("PackageInfo is ready", "packageInfo", prc.pi.Name)
		return prc.reconcilePackageInfoReady(ctx)
	} else if meta.IsStatusConditionFalse(prc.pi.Status.Conditions, string(condition.Ready)) {
		packageInfoCondition := meta.FindStatusCondition(prc.pi.Status.Conditions, string(condition.Ready))
		prc.setShouldUpdate(
			conditions.SetFailed(ctx, prc.EventRecorder, prc.pkg, &prc.pkg.GetStatus().Conditions,
				condition.Reason(packageInfoCondition.Reason), packageInfoCondition.Message))
		return prc.finalize(ctx)
	} else {
		prc.setShouldUpdate(
			conditions.SetUnknown(ctx, &prc.pkg.GetStatus().Conditions, condition.Pending, "PackageInfo status is unknown"))
		return prc.finalize(ctx)
	}
}

type PackageReconcilationContext struct {
	*PackageReconcilerCommon
	pkg                   ctrlpkg.Package
	pi                    *v1alpha1.PackageInfo
	isSuccess             bool
	shouldUpdateStatus    bool
	shouldUpdateResource  bool
	currentOwnedResources []v1alpha1.OwnedResourceRef
	currentOwnedPackages  []v1alpha1.OwnedResourceRef
}

func (r *PackageReconcilationContext) setShouldUpdate(value bool) {
	r.shouldUpdateStatus = r.shouldUpdateStatus || value
}

func (r *PackageReconcilationContext) reconcilePackageInfoReady(ctx context.Context) (ctrl.Result, error) {
	piManifest := r.pi.Status.Manifest
	if piManifest == nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
				condition.UnsupportedFormat, "manifest must not be nil"))
		return r.finalizeNoRequeue(ctx)
	}

	if !r.ensureDependencies(ctx) {
		return r.finalize(ctx)
	}

	var patches []manifestvalues.TargetPatch
	if resolvedValues, err := r.ValueResolver.Resolve(ctx, r.pkg.GetSpec().Values); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
				condition.ValueConfigurationInvalid, err.Error()))
		return r.finalizeWithError(ctx, err)
	} else if err := manifestvalues.ValidateResolvedValues(*piManifest, resolvedValues); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
				condition.ValueConfigurationInvalid, err.Error()))
		return r.finalizeWithError(ctx, err)
	} else if p, err := manifestvalues.GeneratePatches(*piManifest, resolvedValues); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
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
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
					condition.UnsupportedFormat, "manifests not supported"))
			return r.finalizeNoRequeue(ctx)
		}
		adaptersToRun = append(adaptersToRun, r.ManifestAdapter)
	}
	if piManifest.Kustomize != nil {
		if r.KustomizeAdapter == nil {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
					condition.UnsupportedFormat, "kustomize not supported"))
			return r.finalizeNoRequeue(ctx)
		}
		adaptersToRun = append(adaptersToRun, r.KustomizeAdapter)
	}
	if piManifest.Helm != nil {
		if r.HelmAdapter == nil {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
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
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
				condition.InstallationFailed, errs.Error()))
		return r.finalizeWithError(ctx, errs)
	} else if !r.handleAdapterResults(ctx, results) {
		return r.finalize(ctx)
	} else {
		r.afterSuccess(ctx, results)
		return r.finalize(ctx)
	}
}

func (r *PackageReconcilationContext) reconcileAfterDeletion(ctx context.Context) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	r.setShouldUpdate(conditions.SetUnknown(ctx, &r.pkg.GetStatus().Conditions,
		condition.Pending, "Package is being deleted"))

	if r.shouldUpdateStatus {
		telemetry.ForOperator().ReportDelete(r.pkg)
	}

	if slices.Contains(r.pkg.GetFinalizers(), "packages.glasskube.dev/packageDeletion") {
		var err error
		if len(r.pkg.GetStatus().OwnedPackages) != 0 {
			multierr.AppendInto(&err, r.pruneOwnedPackages(ctx, true))
			log.Info("waiting for deletion of required packages")
		} else if len(r.pkg.GetStatus().OwnedPackageInfos) != 0 {
			multierr.AppendInto(&err, r.pruneOwnedPackageInfos(ctx, true))
			log.Info("waiting for deletion of package infos")
		} else {
			r.pkg.SetFinalizers(util.DeleteAll(r.pkg.GetFinalizers(), packageDeletionFinalizer))
			r.shouldUpdateResource = true
		}

		if err != nil {
			return r.finalizeWithError(ctx, err)
		}
	}

	return r.finalizeNoRequeue(ctx)
}

func (r *PackageReconcilationContext) ensureFinalizer() {
	if !slices.Contains(r.pkg.GetFinalizers(), packageDeletionFinalizer) {
		r.pkg.SetFinalizers(append(r.pkg.GetFinalizers(), packageDeletionFinalizer))
		r.shouldUpdateResource = true
	}
}

func (r *PackageReconcilationContext) ensureDependencies(ctx context.Context) bool {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info("ensuring dependencies", "dependencies", r.pi.Status.Manifest.Dependencies)

	var failed []string
	if result, err := r.DependencyManager.Validate(ctx, r.pi.Status.Manifest,
		r.pkg.GetSpec().PackageInfo.Version); err != nil {
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
				condition.InstallationFailed, fmt.Sprintf("error validating dependencies: %v", err)))
		return false
	} else if result.Status == dependency.ValidationResultStatusResolvable {
		for _, requirement := range result.Requirements {
			if requirement.Transitive {
				// Only direct dependencies should be touched in the context of the reconciliation of a package.
				continue
			}
			newPkg := &packagesv1alpha1.ClusterPackage{
				ObjectMeta: metav1.ObjectMeta{
					Name:      requirement.Name,
					Namespace: r.pkg.GetNamespace(),
				},
			}

			newPkg.SetInstalledAsDependency(true)

			repositories, err := r.RepoClientset.Meta().GetReposForPackage(requirement.Name)
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
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
					condition.InstallationFailed, fmt.Sprintf("required package(s) not installed: %v", strings.Join(failed, ","))))
			return false
		}
	} else if result.Status == dependency.ValidationResultStatusConflict {
		var parts []string
		for _, c := range result.Conflicts {
			parts = append(parts, fmt.Sprintf("need version %v of %v but found %v",
				c.Required.Version, c.Actual.Name, c.Actual.Version))
		}
		r.setShouldUpdate(
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
				condition.InstallationFailed, fmt.Sprintf("conflicting dependencies: %v", strings.Join(parts, ","))))
		return false
	}

	var ownedPackages []packagesv1alpha1.OwnedResourceRef
	var waitingFor []string
	// if all requirements fulfilled, status can be checked
	for _, dep := range r.pi.Status.Manifest.Dependencies {
		var requiredPkg packagesv1alpha1.ClusterPackage
		if err := r.Get(ctx, types.NamespacedName{Name: dep.Name}, &requiredPkg); err != nil {
			if apierrors.IsNotFound(err) {
				waitingFor = append(waitingFor, dep.Name)
			} else {
				message := fmt.Sprintf("failed to get required package %v: %v", dep.Name, err)
				r.setShouldUpdate(
					conditions.SetFailed(ctx, r.EventRecorder, r.pkg,
						&r.pkg.GetStatus().Conditions, condition.InstallationFailed, message))
				return false
			}
		} else {
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
		r.setShouldUpdate(conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
			condition.InstallationFailed, message))
		return false
	}

	if len(waitingFor) > 0 {
		message := fmt.Sprintf("waiting for required package(s) %v", strings.Join(waitingFor, ","))
		r.setShouldUpdate(
			conditions.SetUnknown(ctx, &r.pkg.GetStatus().Conditions, condition.Pending, message))
		return false
	}

	ownerutils.Add(&r.currentOwnedPackages, ownedPackages...)

	return true
}

func (r *PackageReconcilationContext) ensurePackageInfo(ctx context.Context) error {
	packageInfo := packagesv1alpha1.PackageInfo{
		ObjectMeta: metav1.ObjectMeta{Name: names.PackageInfoName(r.pkg)},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("PackageInfo", packageInfo.Name)

	log.V(1).Info("ensuring PackageInfo")
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, &packageInfo, func() error {
		packageInfo.Spec = packagesv1alpha1.PackageInfoSpec{
			Name:           r.pkg.GetSpec().PackageInfo.Name,
			Version:        r.pkg.GetSpec().PackageInfo.Version,
			RepositoryName: r.pkg.GetSpec().PackageInfo.RepositoryName,
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

	if changed, err := ownerutils.AddOwnedResourceRef(
		r.Scheme, &r.pkg.GetStatus().OwnedPackageInfos, &packageInfo); err != nil {
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
			conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
				condition.InstallationFailed, firstFailed.Message))
		return false
	} else if firstWaiting != nil {
		r.setShouldUpdate(
			conditions.SetUnknown(ctx, &r.pkg.GetStatus().Conditions, condition.Pending, firstWaiting.Message))
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
		conditions.SetReady(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions, reason, message))
	r.setShouldUpdate(r.pkg.GetStatus().Version != r.pi.Status.Version)
	r.pkg.GetStatus().Version = r.pi.Status.Version
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

	r.setShouldUpdate(ownerutils.Add(&r.pkg.GetStatus().OwnedResources, r.currentOwnedResources...))
	r.setShouldUpdate(ownerutils.Add(&r.pkg.GetStatus().OwnedPackages, r.currentOwnedPackages...))

	var errs error
	if r.isSuccess {
		if err := r.cleanup(ctx); err != nil {
			r.setShouldUpdate(
				conditions.SetFailed(ctx, r.EventRecorder, r.pkg, &r.pkg.GetStatus().Conditions,
					condition.InstallationFailed, fmt.Sprintf("cleanup failed: %v", err)))
			errs = multierr.Append(errs, err)
		}
		log.V(1).Info("cleanup done")
	}

	if r.shouldUpdateStatus {
		if err := r.Status().Update(ctx, r.pkg); err != nil {
			log.Error(err, "package status update failed")
			errs = multierr.Append(errs, err)
		} else {
			log.Info("package status updated")
		}
	} else if r.shouldUpdateResource {
		if err := r.Update(ctx, r.pkg); err != nil {
			log.Error(err, "package update failed")
			errs = multierr.Append(errs, err)
		} else {
			log.Info("package updated")
		}
	}

	return errs
}

func (r *PackageReconcilationContext) cleanup(ctx context.Context) error {
	return multierr.Combine(
		r.pruneOwnedResources(ctx),
		r.pruneOwnedPackageInfos(ctx, false),
		r.pruneOwnedPackages(ctx, false),
	)
}

func (r *PackageReconcilationContext) pruneOwnedResources(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	var errs error
	ownedResourcesCopy := r.pkg.GetStatus().OwnedResources[:]
OuterLoop:
	for _, ref := range r.pkg.GetStatus().OwnedResources {
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
	r.pkg.GetStatus().OwnedResources = ownedResourcesCopy
	return errs
}

func (r *PackageReconcilationContext) pruneOwnedPackageInfos(ctx context.Context, all bool) error {
	log := ctrl.LoggerFrom(ctx)

	allPackages, err := r.getAllPackagesAndClusterPackages(ctx)
	if err != nil {
		return err
	}

	var compositeErr error
	for _, ref := range r.pkg.GetStatus().OwnedPackageInfos {
		if !all && ref.Name == names.PackageInfoName(r.pkg) {
			continue
		}

		// find other packages that require this
		stillUsed := false
	loopAllPackages:
		for _, pkg := range allPackages {
			if pkg.GetName() == r.pkg.GetName() && pkg.GetNamespace() == r.pkg.GetNamespace() {
				// skip current pkg
				continue
			}
			for _, otherRef := range pkg.GetStatus().OwnedPackageInfos {
				if otherRef.Name == ref.Name {
					stillUsed = true
					break loopAllPackages
				}
			}
		}

		if !stillUsed {
			log.V(1).Info("deleting old package info", "PackageInfo", ref.Name)
			if err := r.Delete(ctx, ownerutils.OwnedResourceRefToObject(ref)); client.IgnoreNotFound(err) != nil {
				compositeErr = multierr.Append(compositeErr, err)
				continue
			}
		}

		// Remove the PackageInfo from the owned PackageInfos field of pkg
		ownerutils.RemoveOwnedResourceRef(&r.pkg.GetStatus().OwnedPackageInfos, ref)
		r.setShouldUpdate(true)
	}
	return compositeErr
}

func (r *PackageReconcilationContext) pruneOwnedPackages(ctx context.Context, all bool) error {
	log := ctrl.LoggerFrom(ctx)
	var errs error
	ownedPackagesCopy := r.pkg.GetStatus().OwnedPackages[:]

	allPackages, err := r.getAllPackagesAndClusterPackages(ctx)
	if err != nil {
		return err
	}

OuterLoop:
	for _, ref := range r.pkg.GetStatus().OwnedPackages {
		if !all {
			for _, newRef := range r.currentOwnedPackages {
				if ownerutils.RefersToSameResource(ref, newRef) {
					continue OuterLoop
				}
			}
		}

		stillUsed := false
	AllPackagesLoop:
		for _, otherPkg := range allPackages {
			if !otherPkg.GetDeletionTimestamp().IsZero() ||
				(otherPkg.GetName() == r.pkg.GetName() && otherPkg.GetNamespace() == r.pkg.GetNamespace()) {
				continue
			}
			for _, otherRef := range otherPkg.GetStatus().OwnedPackages {
				if otherRef.Name == ref.Name {
					stillUsed = true
					break AllPackagesLoop
				}
			}
		}

		if !stillUsed {
			var oldReqPkg v1alpha1.ClusterPackage
			if err := r.Get(ctx,
				types.NamespacedName{Name: ref.Name, Namespace: ref.Namespace}, &oldReqPkg); err != nil {
				if apierrors.IsNotFound(err) {
					r.setShouldUpdate(ownerutils.Remove(&ownedPackagesCopy, ref))
				} else {
					multierr.AppendInto(&errs, fmt.Errorf("failed to get old required package: %w", err))
					continue OuterLoop
				}
			} else if oldReqPkg.InstalledAsDependency() {
				// remove the old package if we were the only package owning it and it was previously installed as a dependency
				if err := r.Delete(ctx, ownerutils.OwnedResourceRefToObject(ref),
					&client.DeleteOptions{PropagationPolicy: util.Pointer(metav1.DeletePropagationForeground)},
				); err != nil && !apierrors.IsNotFound(err) {
					errs = multierr.Append(errs, fmt.Errorf("could not prune package: %w", err))
					continue OuterLoop
				} else {
					log.V(1).Info("pruned package", "reference", ref)
				}
			}
		}

		r.setShouldUpdate(ownerutils.Remove(&ownedPackagesCopy, ref))
	}
	r.pkg.GetStatus().OwnedPackages = ownedPackagesCopy
	return errs
}

func (r *PackageReconcilerCommon) getAllPackagesAndClusterPackages(ctx context.Context) ([]ctrlpkg.Package, error) {
	var pkgList packagesv1alpha1.PackageList
	if err := r.List(ctx, &pkgList, &client.ListOptions{}); err != nil {
		return nil, err
	}

	var cpkgList packagesv1alpha1.ClusterPackageList
	if err := r.List(ctx, &cpkgList, &client.ListOptions{}); err != nil {
		return nil, err
	}

	allPackages := make([]ctrlpkg.Package, len(pkgList.Items)+len(cpkgList.Items))
	for i := range pkgList.Items {
		allPackages[i] = &pkgList.Items[i]
	}
	for i := range cpkgList.Items {
		allPackages[len(pkgList.Items)+i] = &cpkgList.Items[i]
	}

	return allPackages, nil
}

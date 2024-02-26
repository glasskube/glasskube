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
	"slices"

	"go.uber.org/multierr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/controller/conditions"
	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/controller/requeue"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/pkg/condition"
)

// PackageReconciler reconciles a Package object
type PackageReconciler struct {
	client.Client
	record.EventRecorder
	*owners.OwnerManager
	Scheme    *runtime.Scheme
	Helm      manifest.ManifestAdapter
	Kustomize manifest.ManifestAdapter
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
	log := ctrl.LoggerFrom(ctx)
	var pkg packagesv1alpha1.Package

	if err := r.Get(ctx, req.NamespacedName, &pkg); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !pkg.DeletionTimestamp.IsZero() {
		err := conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions,
			condition.Pending, "Package is being deleted")
		return ctrl.Result{}, err
	}

	packageInfo, changed, err := r.ensurePackageInfo(ctx, &pkg)
	if err != nil {
		return requeue.Always(ctx, err)
	}
	if changed {
		if err := r.Status().Update(ctx, &pkg); err != nil {
			return requeue.Always(ctx, err)
		}
	}

	if meta.IsStatusConditionTrue(packageInfo.Status.Conditions, string(condition.Ready)) {
		log.V(1).Info("PackageInfo is ready", "packageInfo", packageInfo.Name)
		piManifest := packageInfo.Status.Manifest
		if piManifest == nil {
			err := conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "manifest must not be nil")
			return requeue.Always(ctx, err)
		}

		if len(piManifest.Manifests) > 0 {
			err := r.reconcilePlainManifests(ctx, &pkg, *packageInfo, piManifest)
			return requeue.Always(ctx, err)
		} else if piManifest.Helm != nil {
			if r.Helm != nil {
				err := r.reconcileManifestWithAdapter(ctx, r.Helm, pkg, *packageInfo, piManifest)
				return requeue.Always(ctx, err)
			} else {
				err := conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "helm manifest not supported")
				return requeue.Always(ctx, err)
			}
		} else if piManifest.Kustomize != nil {
			if r.Kustomize != nil {
				err := r.reconcileManifestWithAdapter(ctx, r.Kustomize, pkg, *packageInfo, piManifest)
				return requeue.Always(ctx, err)
			} else {
				err := conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "kustomize manifest not supported")
				return requeue.Always(ctx, err)
			}
		} else {
			err := r.afterSuccess(ctx, &pkg, *packageInfo, []packagesv1alpha1.OwnedResourceRef{}, condition.UpToDate, "PackageInfo has nothing to apply (no helm or kustomize manifest present)")
			return requeue.Always(ctx, err)
		}
	} else if meta.IsStatusConditionFalse(packageInfo.Status.Conditions, string(condition.Ready)) {
		packageInfoCondition := meta.FindStatusCondition(packageInfo.Status.Conditions, string(condition.Ready))
		err := conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, &pkg, &pkg.Status.Conditions, condition.Reason(packageInfoCondition.Reason), packageInfoCondition.Message)
		return requeue.Always(ctx, err)
	} else {
		err := conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.Pending, "PackageInfo status is unknown")
		return requeue.Always(ctx, err)
	}
}

func (r *PackageReconciler) ensurePackageInfo(ctx context.Context, pkg *packagesv1alpha1.Package) (*packagesv1alpha1.PackageInfo, bool, error) {
	packageInfo := packagesv1alpha1.PackageInfo{
		ObjectMeta: metav1.ObjectMeta{Name: generatePackageInfoName(*pkg)},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("PackageInfo", packageInfo.Name)

	log.V(1).Info("ensuring PackageInfo")
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, &packageInfo, func() error {
		if err := r.SetOwner(pkg, &packageInfo, owners.BlockOwnerDeletion); err != nil {
			return fmt.Errorf("unable to set ownerReference on PackageInfo: %w", err)
		}
		packageInfo.Spec = packagesv1alpha1.PackageInfoSpec{
			Name:          pkg.Spec.PackageInfo.Name,
			Version:       pkg.Spec.PackageInfo.Version,
			RepositoryUrl: pkg.Spec.PackageInfo.RepositoryUrl,
		}
		return nil
	})
	if err != nil {
		return nil, false, fmt.Errorf("could not create or update PackageInfo: %w", err)
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
			return nil, false, err
		}
	}

	changed := addOwnedResourceRef(&pkg.Status.OwnedPackageInfos, &packageInfo, &packageInfo)
	return &packageInfo, changed, nil
}

func generatePackageInfoName(pkg packagesv1alpha1.Package) string {
	if pkg.Spec.PackageInfo.Version != "" {
		return escapeResourceName(fmt.Sprintf("%v--%v", pkg.Spec.PackageInfo.Name, pkg.Spec.PackageInfo.Version))
	} else {
		return escapeResourceName(pkg.Spec.PackageInfo.Name)
	}
}

func (r *PackageReconciler) reconcileManifestWithAdapter(ctx context.Context, adapter manifest.ManifestAdapter, pkg packagesv1alpha1.Package, packageInfo packagesv1alpha1.PackageInfo, piManifest *packagesv1alpha1.PackageManifest) error {
	if result, err := adapter.Reconcile(ctx, &pkg, piManifest); err != nil {
		log := ctrl.LoggerFrom(ctx)
		log.Error(err, "could reconcile manifest")
		return conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "error during manifest reconciliation: "+err.Error())
	} else if result.IsReady() {
		return r.afterSuccess(ctx, &pkg, packageInfo, []packagesv1alpha1.OwnedResourceRef{}, condition.InstallationSucceeded, result.Message)
	} else if result.IsWaiting() {
		return conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.Pending, result.Message)
	} else if result.IsFailed() {
		return conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, &pkg, &pkg.Status.Conditions, condition.InstallationFailed, result.Message)
	} else {
		return conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, result.Message)
	}
}

func (r *PackageReconciler) reconcilePlainManifests(ctx context.Context, pkg *packagesv1alpha1.Package, packageInfo packagesv1alpha1.PackageInfo, manifest *packagesv1alpha1.PackageManifest) error {
	allOwned := make([]packagesv1alpha1.OwnedResourceRef, 0)
	for _, m := range manifest.Manifests {
		if owned, err := r.reconcilePlainManifest(ctx, *pkg, m); err != nil {
			return conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, pkg, &pkg.Status.Conditions, condition.InstallationFailed, err.Error())
		} else {
			allOwned = append(allOwned, owned...)
		}
	}

	return r.afterSuccess(ctx, pkg, packageInfo, allOwned, condition.InstallationSucceeded, "all manifests reconciled")
}

func (r *PackageReconciler) reconcilePlainManifest(ctx context.Context, pkg packagesv1alpha1.Package, manifest packagesv1alpha1.PlainManifest) ([]packagesv1alpha1.OwnedResourceRef, error) {
	log := ctrl.LoggerFrom(ctx)
	objectsToApply, err := clientutils.FetchResources(manifest.Url)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("fetched "+manifest.Url, "objectCount", len(*objectsToApply))

	ownedResources := make([]packagesv1alpha1.OwnedResourceRef, 0, len(*objectsToApply))

	// TODO: check if namespace is terminating before applying

	for _, obj := range *objectsToApply {
		if err := r.SetOwner(&pkg, &obj, owners.BlockOwnerDeletion); err != nil {
			return nil, fmt.Errorf("could set owner reference: %w", err)
		}
		if err := r.Patch(ctx, &obj, client.Apply, client.FieldOwner("packages.glasskube.dev/package-controller"), client.ForceOwnership); err != nil {
			return nil, fmt.Errorf("could not apply resource: %w", err)
		}
		log.V(1).Info("applied resource", "kind", obj.GroupVersionKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
		addOwnedResourceRef(&ownedResources, &obj, &obj)
	}
	return ownedResources, nil
}

func (r *PackageReconciler) afterSuccess(ctx context.Context, pkg *packagesv1alpha1.Package, packageInfo packagesv1alpha1.PackageInfo, newOwnedResources []packagesv1alpha1.OwnedResourceRef, reason condition.Reason, message string) error {
	if err := r.pruneOwnedResources(ctx, pkg, newOwnedResources); err != nil {
		return conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, pkg, &pkg.Status.Conditions, condition.InstallationFailed, err.Error())
	}
	ownedResourcesChanged := !slices.Equal(pkg.Status.OwnedResources, newOwnedResources)
	if ownedResourcesChanged {
		pkg.Status.OwnedResources = newOwnedResources[:]
	}

	ownedPackageInfosChanged, err := r.pruneOwnedPackageInfos(ctx, pkg, packageInfo)
	if err != nil {
		return conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, pkg, &pkg.Status.Conditions, condition.InstallationFailed, err.Error())
	}

	conditionsChanged := conditions.SetReady(ctx, r.EventRecorder, pkg, &pkg.Status.Conditions, reason, message)

	pkg.Status.Version = packageInfo.Status.Version

	if ownedPackageInfosChanged || ownedResourcesChanged || conditionsChanged {
		return r.Status().Update(ctx, pkg)
	}

	return nil
}

func (r *PackageReconciler) pruneOwnedResources(ctx context.Context, pkg *packagesv1alpha1.Package, newOwnedResources []packagesv1alpha1.OwnedResourceRef) error {
	log := ctrl.LoggerFrom(ctx)

OuterLoop:
	for _, ref := range pkg.Status.OwnedResources {
		for _, newRef := range newOwnedResources {
			if refersToSameResource(ref, newRef) {
				continue OuterLoop
			}
		}
		if err := r.Delete(ctx, ownedResourceRefToObject(ref)); err != nil {
			if !apierrors.IsNotFound(err) {
				return fmt.Errorf("could not prune resource: %w", err)
			}
		}
		log.V(1).Info("pruned resource", "reference", ref)
	}
	return nil
}

func (r *PackageReconciler) pruneOwnedPackageInfos(ctx context.Context, pkg *packagesv1alpha1.Package, current packagesv1alpha1.PackageInfo) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	currentRef := toOwnedResourceRef(&current, &current)
	var compositeErr error
	changed := false
	for _, ref := range pkg.Status.OwnedPackageInfos {
		if !refersToSameResource(ref, currentRef) {
			var packageInfo packagesv1alpha1.PackageInfo
			if err := r.Get(ctx, client.ObjectKeyFromObject(ownedResourceRefToObject(ref)), &packageInfo); apierrors.IsNotFound(err) {
				log.Info("PackageInfo not found", "PackageInfo", ref.Name)
			} else if err != nil {
				compositeErr = multierr.Append(compositeErr, err)
				continue
			} else {
				if owned, err := objHasOwner(&packageInfo, pkg); err != nil {
					compositeErr = multierr.Append(compositeErr, err)
					continue
				} else if owned {
					// Remove the owner reference for pkg
					if err := controllerutil.RemoveOwnerReference(pkg, &packageInfo, r.Scheme); err != nil {
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
			removeOwnedResourceRef(&pkg.Status.OwnedPackageInfos, ref)
			changed = true
		}
	}
	return changed, compositeErr
}

// SetupWithManager sets up the controller with the Manager.
func (r *PackageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.OwnerManager == nil {
		r.OwnerManager = owners.NewOwnerManager(r.Scheme)
	}
	controllerBuilder := ctrl.NewControllerManagedBy(mgr).For(&packagesv1alpha1.Package{})
	if r.Helm != nil {
		if err := r.Helm.ControllerInit(controllerBuilder, r.Client, r.Scheme); err != nil {
			return err
		}
	}
	if r.Kustomize != nil {
		if err := r.Kustomize.ControllerInit(controllerBuilder, r.Client, r.Scheme); err != nil {
			return err
		}
	}
	return controllerBuilder.
		Owns(&packagesv1alpha1.PackageInfo{}, builder.MatchEveryOwner).
		Complete(r)
}

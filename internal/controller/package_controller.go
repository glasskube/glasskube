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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/conditions"
	"github.com/glasskube/glasskube/internal/controller/requeue"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/pkg/condition"
)

// PackageReconciler reconciles a Package object
type PackageReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Helm      manifest.ManifestAdapter
	Kustomize manifest.ManifestAdapter
}

//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packages/finalizers,verbs=update

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
	log := log.FromContext(ctx)
	var pkg packagesv1alpha1.Package

	if err := r.Get(ctx, req.NamespacedName, &pkg); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(1).Info("Failed to fetch Package: " + err.Error())
			return requeue.Never(ctx, nil)
		} else {
			return requeue.Always(ctx, fmt.Errorf("failed to fetch package: %w", err))
		}
	}

	if err := conditions.SetInitialAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions); err != nil {
		return requeue.Always(ctx, err)
	}

	actualPackageInfo, err := r.ensurePackageInfo(ctx, pkg)
	if err != nil {
		return requeue.Always(ctx, err)
	}

	if meta.IsStatusConditionTrue(actualPackageInfo.Status.Conditions, string(condition.Ready)) {
		log.V(1).Info("PackageInfo is ready", "packageinfo", actualPackageInfo.Name)
		piManifest := actualPackageInfo.Status.Manifest
		if piManifest == nil {
			err := conditions.SetFailedAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "manifest must not be nil")
			return requeue.Always(ctx, err)
		}

		if piManifest.Helm != nil {
			if r.Helm != nil {
				err := r.reconcileManifestWithAdapter(ctx, r.Helm, pkg, piManifest)
				return requeue.Always(ctx, err)
			} else {
				err := conditions.SetFailedAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "helm manifest not supported")
				return requeue.Always(ctx, err)
			}
		} else if piManifest.Kustomize != nil {
			if r.Kustomize != nil {
				err := r.reconcileManifestWithAdapter(ctx, r.Kustomize, pkg, piManifest)
				return requeue.Always(ctx, err)
			} else {
				err := conditions.SetFailedAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "kustomize manifest not supported")
				return requeue.Always(ctx, err)
			}
		} else {
			err := conditions.SetReadyAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.UpToDate, "PackageInfo has nothing to apply (no helm or kustomize manifest present)")
			return requeue.Always(ctx, err)
		}
	} else if meta.IsStatusConditionFalse(actualPackageInfo.Status.Conditions, string(condition.Ready)) {
		packageInfoCondition := meta.FindStatusCondition(actualPackageInfo.Status.Conditions, string(condition.Ready))
		err := conditions.SetFailedAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.Reason(packageInfoCondition.Reason), packageInfoCondition.Message)
		return requeue.Always(ctx, err)
	} else {
		err := conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.Pending, "PackageInfo status is unknown")
		return requeue.Always(ctx, err)
	}
}

func (r *PackageReconciler) reconcileManifestWithAdapter(ctx context.Context, adapter manifest.ManifestAdapter, pkg packagesv1alpha1.Package, piManifest *packagesv1alpha1.PackageManifest) error {
	log := log.FromContext(ctx)

	if result, err := adapter.Reconcile(ctx, r.Client, &pkg, piManifest); err != nil {
		log.Error(err, "could reconcile manifest")
		err = conditions.SetFailedAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, "error during manifest reconciliation")
		if err != nil {
			return fmt.Errorf("could not set conditions: %w", err)
		}
	} else if result.IsReady() {
		conditions.SetReadyAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.InstallationSucceeded, result.Message)
		if err != nil {
			return fmt.Errorf("could not set conditions: %w", err)
		}
	} else if result.IsWaiting() {
		conditions.SetReadyAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.Pending, result.Message)
		if err != nil {
			return fmt.Errorf("could not set conditions: %w", err)
		}
	} else {
		conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions, condition.UnsupportedFormat, result.Message)
		if err != nil {
			return fmt.Errorf("could not set conditions: %w", err)
		}
	}

	return nil
}

func (r *PackageReconciler) desiredPackageInfo(pkg *packagesv1alpha1.Package) (*packagesv1alpha1.PackageInfo, error) {
	packageInfo := &packagesv1alpha1.PackageInfo{
		ObjectMeta: metav1.ObjectMeta{
			Name: pkg.Spec.PackageInfo.Name,
		},
		Spec: packagesv1alpha1.PackageInfoSpec{
			Name:          pkg.Spec.PackageInfo.Name,
			RepositoryUrl: pkg.Spec.PackageInfo.RepositoryUrl,
		},
	}
	err := controllerutil.SetOwnerReference(pkg, packageInfo, r.Scheme)
	return packageInfo, err
}

func (r *PackageReconciler) ensurePackageInfo(ctx context.Context, pkg packagesv1alpha1.Package) (*packagesv1alpha1.PackageInfo, error) {
	log := log.FromContext(ctx)

	desiredPackageInfo, err := r.desiredPackageInfo(&pkg)
	if err != nil {
		return nil, err
	}

	var actualPackageInfo packagesv1alpha1.PackageInfo
	if err := r.Get(ctx, types.NamespacedName{Name: pkg.Spec.PackageInfo.Name}, &actualPackageInfo); apierrors.IsNotFound(err) {
		log.V(1).Info("Creating PackageInfo", "packageinfo", desiredPackageInfo.Name)
		err = r.Create(ctx, desiredPackageInfo)
		if err != nil {
			log.Error(err, "Failed to create PackageInfo", "packageinfo", desiredPackageInfo.Name)
		}
		return desiredPackageInfo, err
	} else if err != nil {
		log.Error(err, "Failed to fetch PackageInfo", "packageinfo", desiredPackageInfo.Name)
		return nil, err
	}

	if err = r.updatePackageInfoIfNeeded(ctx, &pkg, desiredPackageInfo, &actualPackageInfo); err != nil {
		return nil, err
	}

	return &actualPackageInfo, nil
}

func (r *PackageReconciler) updatePackageInfoIfNeeded(ctx context.Context, pkg *packagesv1alpha1.Package, desiredPackageInfo, actualPackageInfo *packagesv1alpha1.PackageInfo) error {
	log := log.FromContext(ctx)

	updateNeeded := false

	if actualPackageInfo.Spec.Name != desiredPackageInfo.Spec.Name ||
		actualPackageInfo.Spec.RepositoryUrl != desiredPackageInfo.Spec.RepositoryUrl {
		log.V(1).Info("spec is out of date", "packageinfo", actualPackageInfo.Name)
		actualPackageInfo.Spec.Name = desiredPackageInfo.Spec.Name
		actualPackageInfo.Spec.RepositoryUrl = desiredPackageInfo.Spec.RepositoryUrl
		updateNeeded = true
	}

	if owned, err := objHasOwner(actualPackageInfo, pkg); err != nil {
		return err
	} else if !owned {
		log.V(1).Info("owner is missing", "packageinfo", actualPackageInfo.Name)
		// add the current pkg as owner
		if err := controllerutil.SetOwnerReference(pkg, actualPackageInfo, r.Scheme); err != nil {
			return err
		}
		updateNeeded = true
	}

	if updateNeeded {
		log.V(1).Info("Updating PackageInfo", "packageinfo", actualPackageInfo.Name)
		if err := r.Update(ctx, actualPackageInfo); err != nil {
			log.Error(err, "Failed update PackageInfo")
			return err
		}
	} else {
		log.V(1).Info("PackageInfo is already up-to-date", "packageinfo", actualPackageInfo.Name)
	}

	return nil
}

func objHasOwner(obj, owner client.Object) (bool, error) {
	refs := obj.GetOwnerReferences()
	for _, ref := range refs {
		if ref.Name != owner.GetName() {
			continue
		}
		if gv, err := schema.ParseGroupVersion(ref.APIVersion); err != nil {
			return false, err
		} else if owner.GetObjectKind().GroupVersionKind() == gv.WithKind(ref.Kind) {
			return true, nil
		}
	}
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PackageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	controllerBuilder := ctrl.NewControllerManagedBy(mgr).For(&packagesv1alpha1.Package{})
	if r.Helm != nil {
		r.Helm.ControllerInit(controllerBuilder)
	}
	if r.Kustomize != nil {
		r.Kustomize.ControllerInit(controllerBuilder)
	}
	return controllerBuilder.
		Owns(&packagesv1alpha1.PackageInfo{}, builder.MatchEveryOwner).
		Complete(r)
}

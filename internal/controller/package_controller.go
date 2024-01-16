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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/api/v1alpha1/condition"
)

// PackageReconciler reconciles a Package object
type PackageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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
		log.Error(err, "Failed to fetch Package")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Initialize the status conditions
	if pkg.Status.Conditions == nil || len(pkg.Status.Conditions) == 0 {
		meta.SetStatusCondition(
			&pkg.Status.Conditions,
			metav1.Condition{Type: condition.Ready, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"},
		)
		if err := r.Status().Update(ctx, &pkg); err != nil {
			log.Error(err, "Failed to update Package status")
			return ctrl.Result{}, err
		}
		if err := r.Get(ctx, req.NamespacedName, &pkg); err != nil {
			log.Error(err, "Failed to re-fetch Package")
			return ctrl.Result{}, err
		}
	}

	desiredPackageInfo, err := r.desiredPackageInfo(&pkg)
	if err != nil {
		return ctrl.Result{}, err
	}

	var actualPackageInfo packagesv1alpha1.PackageInfo
	if err := r.Get(ctx, types.NamespacedName{Name: pkg.Spec.PackageInfo.Name}, &actualPackageInfo); apierrors.IsNotFound(err) {
		log.V(1).Info("Creating PackageInfo", "packageinfo", desiredPackageInfo.Name)
		err = r.Create(ctx, desiredPackageInfo)
		if err != nil {
			log.Error(err, "Failed to create PackageInfo", "packageinfo", desiredPackageInfo.Name)
		}
		return ctrl.Result{}, err
	} else if err != nil {
		log.Error(err, "Failed to fetch PackageInfo", "packageinfo", desiredPackageInfo.Name)
		return ctrl.Result{Requeue: true}, err
	}

	if err = r.updatePackageInfoIfNeeded(ctx, &pkg, desiredPackageInfo, &actualPackageInfo); err != nil {
		return ctrl.Result{}, err
	}

	if meta.IsStatusConditionTrue(actualPackageInfo.Status.Conditions, condition.Ready) {
		log.V(1).Info("PackageInfo is ready", "packageinfo", desiredPackageInfo.Name)
		// TODO: Handle PackageInfo with condition Ready=True
	} else {
		log.V(1).Info("PackageInfo is not ready", "packageinfo", desiredPackageInfo.Name)
	}

	return ctrl.Result{}, nil
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
	err := controllerutil.SetControllerReference(pkg, packageInfo, r.Scheme)
	return packageInfo, err
}

func (r *PackageReconciler) updatePackageInfoIfNeeded(ctx context.Context, pkg *packagesv1alpha1.Package, desiredPackageInfo, actualPackageInfo *packagesv1alpha1.PackageInfo) error {
	log := log.FromContext(ctx)

	if actualPackageInfo.Spec.Name != desiredPackageInfo.Spec.Name ||
		actualPackageInfo.Spec.RepositoryUrl != desiredPackageInfo.Spec.RepositoryUrl {

		log.V(1).Info("Updating PackageInfo", "packageinfo", actualPackageInfo.Name)

		actualPackageInfo.Spec.Name = desiredPackageInfo.Spec.Name
		actualPackageInfo.Spec.RepositoryUrl = desiredPackageInfo.Spec.RepositoryUrl

		if err := r.Update(ctx, actualPackageInfo); err != nil {
			log.Error(err, "Failed update PackageInfo")
			return err
		}
	} else {
		log.V(1).Info("PackageInfo is already up-to-date", "packageinfo", actualPackageInfo.Name)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PackageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&packagesv1alpha1.Package{}).
		Owns(&packagesv1alpha1.PackageInfo{}).
		Complete(r)
}

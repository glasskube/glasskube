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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/conditions"
	"github.com/glasskube/glasskube/pkg/condition"
)

// ClusterPackageReconciler reconciles a ClusterPackage object
type ClusterPackageReconciler struct {
	PackageReconcilerCommon
}

//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=clusterpackages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=clusterpackages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=clusterpackages/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterPackage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ClusterPackageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pkg v1alpha1.ClusterPackage

	if err := r.Get(ctx, req.NamespacedName, &pkg); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO: consider moving this to common
	if !pkg.DeletionTimestamp.IsZero() {
		_, err := conditions.SetUnknownAndUpdate(ctx, r.Client, &pkg, &pkg.Status.Conditions,
			condition.Pending, "Package is being deleted")
		// TODO: make telementry work
		// if changed {
		// 	telemetry.ForOperator().ReportDelete(&pkg)
		// }
		return ctrl.Result{}, err
	}

	// telemetry.ForOperator().ReconcilePackage(&pkg)

	return r.reconcile(ctx, &pkg)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterPackageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if bld, err := r.baseSetup(mgr, &v1alpha1.Package{}); err != nil {
		return err
	} else {
		return bld.Complete(r)
	}
}

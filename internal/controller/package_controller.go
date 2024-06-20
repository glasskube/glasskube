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

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PackageReconciler reconciles a Package object
type PackageReconciler struct {
	PackageReconcilerCommon
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
	var pkg v1alpha1.Package
	if err := r.Get(ctx, req.NamespacedName, &pkg); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.reconcile(ctx, &pkg)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PackageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if bld, err := r.baseSetup(mgr, &v1alpha1.Package{}, r); err != nil {
		return err
	} else {
		return bld.Complete(r)
	}
}

func (r *PackageReconciler) ListPackages(ctx context.Context) ([]ctrlpkg.Package, error) {
	var l v1alpha1.PackageList
	if err := r.Client.List(ctx, &l); err != nil {
		return nil, err
	}
	res := make([]ctrlpkg.Package, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res, nil
}

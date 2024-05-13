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

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/requeue"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	repotypes "github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/pkg/condition"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PackageRepositoryReconciler reconciles a PackageRepository object
type PackageRepositoryReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	RepoClient repoclient.RepoClientset
}

//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packagerepositories,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packagerepositories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packagerepositories/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PackageRepository object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *PackageRepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var repo packagesv1alpha1.PackageRepository
	if err := r.Get(ctx, req.NamespacedName, &repo); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var index repotypes.PackageRepoIndex
	var cond metav1.Condition
	err := r.RepoClient.ForRepo(repo).FetchPackageRepoIndex(&index)
	if err != nil {
		cond = metav1.Condition{
			Type:    string(condition.Ready),
			Status:  metav1.ConditionFalse,
			Reason:  string(condition.SyncFailed),
			Message: err.Error(),
		}
	} else {
		cond = metav1.Condition{
			Type:    string(condition.Ready),
			Status:  metav1.ConditionTrue,
			Reason:  string(condition.SyncCompleted),
			Message: fmt.Sprintf("repo has %v packages", len(index.Packages)),
		}
	}

	if meta.SetStatusCondition(&repo.Status.Conditions, cond) {
		multierr.AppendInto(&err, r.Status().Update(ctx, &repo))
	}

	return requeue.Always(ctx, err)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PackageRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&packagesv1alpha1.PackageRepository{}).
		Complete(r)
}

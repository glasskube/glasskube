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
	"time"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/conditions"
	"github.com/glasskube/glasskube/internal/controller/requeue"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/condition"
	"go.uber.org/multierr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PackageInfoReconciler reconciles a PackageInfo object
type PackageInfoReconciler struct {
	client.Client
	record.EventRecorder
	Scheme *runtime.Scheme
}

var (
	// 5 minutes in nanoseconds.
	// TODO: let users configure this value per PackageInfo or per repository
	repositorySyncInterval = 5 * time.Minute
)

//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packageinfos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packageinfos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packageinfos/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PackageInfo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PackageInfoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var packageInfo packagesv1alpha1.PackageInfo

	if err := r.Get(ctx, req.NamespacedName, &packageInfo); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(1).Info("Failed to fetch PackageInfo: " + err.Error())
			return ctrl.Result{}, nil
		} else {
			return requeue.Always(ctx, err)
		}
	}

	if err := conditions.SetInitialAndUpdate(ctx, r.Client, &packageInfo, &packageInfo.Status.Conditions); err != nil {
		return requeue.Always(ctx, err)
	}

	if shouldSyncFromRepo(packageInfo) {
		if err := repo.FetchPackageManifest(ctx, &packageInfo); err != nil {
			err1 := conditions.SetFailedAndUpdate(ctx, r.Client, r.EventRecorder, &packageInfo, &packageInfo.Status.Conditions, condition.SyncFailed, err.Error())
			return requeue.Always(ctx, multierr.Append(err, err1))
		} else {
			now := metav1.Now()
			packageInfo.Status.LastUpdateTimestamp = &now
			conditions.SetReady(ctx, r.EventRecorder, &packageInfo, &packageInfo.Status.Conditions, condition.SyncCompleted, "PackageInfo is up-to-date")
			if err := r.Status().Update(ctx, &packageInfo); err != nil {
				r.Event(&packageInfo, "Warning", string(condition.SyncFailed), err.Error())
				return requeue.Always(ctx, err)
			}
		}
	}

	return requeue.Always(ctx, nil)
}

func shouldSyncFromRepo(pi packagesv1alpha1.PackageInfo) bool {
	return pi.Status.LastUpdateTimestamp == nil || time.Since(pi.Status.LastUpdateTimestamp.Time) > repositorySyncInterval
}

// SetupWithManager sets up the controller with the Manager.
func (r *PackageInfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&packagesv1alpha1.PackageInfo{}).
		Complete(r)
}

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
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/api/v1alpha1/condition"
	"github.com/glasskube/glasskube/internal/controller/requeue"
)

// PackageInfoReconciler reconciles a PackageInfo object
type PackageInfoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	// 5 minutes in nanoseconds.
	// TODO: let users configure this value per PackageInfo or per repository
	repositorySyncInterval = 5 * time.Minute
	defaultRepositoryUrl   = "https://packages.dl.glasskube.dev/packages/"
)

//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packageinfos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packageinfos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=packages.glasskube.dev,resources=packageinfos/finalizers,verbs=update

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
			return requeue.Never(ctx, nil)
		} else {
			return requeue.Always(ctx, err)
		}
	}

	if err := r.setStatusConditionInitialAndUpdate(ctx, &packageInfo); err != nil {
		return requeue.Always(ctx, err)
	}

	if shouldSyncFromRepo(packageInfo) {
		if err := fetchManifestFromRepo(ctx, &packageInfo); err != nil {
			log.Error(err, "could not fetch package manifest")
			if err := r.setStatusConditionFailedAndUpdate(ctx, &packageInfo); err != nil {
				return requeue.Always(ctx, err)
			}
		} else {
			now := metav1.Now()
			packageInfo.Status.LastUpdateTimestamp = &now
			r.setStatusConditionReady(ctx, &packageInfo)
			if err := r.Status().Update(ctx, &packageInfo); err != nil {
				return requeue.Always(ctx, err)
			}
		}
	}

	return requeue.Always(ctx, nil)
}

func (r *PackageInfoReconciler) setStatusConditionInitialAndUpdate(ctx context.Context, pi *packagesv1alpha1.PackageInfo) error {
	if pi.Status.Conditions == nil || len(pi.Status.Conditions) == 0 {
		log := log.FromContext(ctx)
		log.V(1).Info("set initial conditions")
		return r.setStatusConditionsAndUpdate(ctx, pi,
			metav1.Condition{Type: condition.Ready, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"},
			metav1.Condition{Type: condition.Failed, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"},
		)
	} else {
		return nil
	}
}

func (r *PackageInfoReconciler) setStatusConditionReady(ctx context.Context, pi *packagesv1alpha1.PackageInfo) bool {
	log := log.FromContext(ctx)
	log.V(1).Info("set condition to ready")
	return setStatusConditions(pi,
		metav1.Condition{Type: condition.Ready, Status: metav1.ConditionTrue, Reason: "SyncCompleted"},
		metav1.Condition{Type: condition.Failed, Status: metav1.ConditionFalse, Reason: "SyncCompleted"},
	)
}

func (r *PackageInfoReconciler) setStatusConditionFailedAndUpdate(ctx context.Context, pi *packagesv1alpha1.PackageInfo) error {
	log := log.FromContext(ctx)
	log.V(1).Info("set condition to failed")
	return r.setStatusConditionsAndUpdate(ctx, pi,
		metav1.Condition{Type: condition.Ready, Status: metav1.ConditionFalse, Reason: "SyncCompleted"},
		metav1.Condition{Type: condition.Failed, Status: metav1.ConditionTrue, Reason: "SyncFailed"},
	)
}

func setStatusConditions(pi *packagesv1alpha1.PackageInfo, conditions ...metav1.Condition) bool {
	needsUpdate := false
	for _, condition := range conditions {
		changed := meta.SetStatusCondition(&pi.Status.Conditions, condition)
		needsUpdate = changed || needsUpdate
	}
	return needsUpdate
}

func (r *PackageInfoReconciler) setStatusConditionsAndUpdate(ctx context.Context, pi *packagesv1alpha1.PackageInfo, conditions ...metav1.Condition) error {
	log := log.FromContext(ctx)
	if setStatusConditions(pi, conditions...) {
		log.V(1).Info("Updating status after conditions changed")
		if err := r.Status().Update(ctx, pi); err != nil {
			log.Error(err, "Failed to update PackageInfo status")
			return err
		}
		if err := r.Get(ctx, client.ObjectKeyFromObject(pi), pi); err != nil {
			log.Error(err, "Failed to re-fetch PackageInfo")
			return err
		}
	}
	return nil
}

func getManifestUrl(pi packagesv1alpha1.PackageInfo) (string, error) {
	var baseUrl string
	if len(pi.Spec.RepositoryUrl) > 0 {
		baseUrl = pi.Spec.RepositoryUrl
	} else {
		baseUrl = defaultRepositoryUrl
	}
	return url.JoinPath(baseUrl, pi.Spec.Name, "package.yaml")
}

// TODO: Migrate to client package once it is available
func fetchManifestFromRepo(ctx context.Context, pi *packagesv1alpha1.PackageInfo) error {
	log := log.FromContext(ctx)
	url, err := getManifestUrl(*pi)
	if err != nil {
		log.Error(err, "can not get manifest url")
		return err
	}
	log.Info("starting to fetch manifset from " + url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("could not fetch package manifest: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var manifest packagesv1alpha1.PackageManifest
	if err := yaml.Unmarshal(body, &manifest); err != nil {
		return err
	}
	pi.Status.Manifest = &manifest
	return nil
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

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

package webhook

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PackageRepositoryValidatingWebhook struct {
	client.Client
}

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (p *PackageRepositoryValidatingWebhook) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.PackageRepository{}).
		WithValidator(p).
		Complete()
}

// +kubebuilder:webhook:path=/validate-packages-glasskube-dev-v1alpha1-packagerepository,mutating=false,failurePolicy=fail,sideEffects=None,groups=packages.glasskube.dev,resources=packagerepositories,verbs=update;delete,versions=v1alpha1,name=vpackagerepository.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &PackageRepositoryValidatingWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type
func (p *PackageRepositoryValidatingWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := ctrl.LoggerFrom(ctx)
	if repo, ok := obj.(*v1alpha1.PackageRepository); ok {
		log.Info("validate create", "name", repo.Name)
		return nil, nil
	}
	return nil, ErrInvalidObject
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type
func (p *PackageRepositoryValidatingWebhook) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (admission.Warnings, error) {
	log := ctrl.LoggerFrom(ctx)
	if oldRepo, ok := oldObj.(*v1alpha1.PackageRepository); ok {
		if newRepo, ok := newObj.(*v1alpha1.PackageRepository); ok {
			log.Info("validate update", "name", newRepo.Name)
			if oldRepo.Spec.Url != newRepo.Spec.Url {
				return nil, p.validateUpdateOrDelete(ctx, oldRepo)
			} else {
				return nil, nil
			}
		}
	}
	return nil, ErrInvalidObject
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (p *PackageRepositoryValidatingWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := ctrl.LoggerFrom(ctx)
	if repo, ok := obj.(*v1alpha1.PackageRepository); ok {
		log.Info("validate delete", "name", repo.Name)
		return nil, p.validateUpdateOrDelete(ctx, repo)
	}
	return nil, ErrInvalidObject
}

func (p *PackageRepositoryValidatingWebhook) validateUpdateOrDelete(ctx context.Context, repo *v1alpha1.PackageRepository) error {
	var pkgLs v1alpha1.PackageInfoList
	if err := p.Client.List(ctx, &pkgLs); err != nil {
		return err
	}
	var repoPkgLs []v1alpha1.PackageInfoSpec
	for _, item := range pkgLs.Items {
		if item.Spec.RepositoryName == repo.Name {
			repoPkgLs = append(repoPkgLs, item.Spec)
		}
	}
	if len(repoPkgLs) > 0 {
		return newErrPackagesInstalled(repoPkgLs)
	}
	return nil
}

package webhook

import (
	"context"
	"errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PackageValidatingWebhook struct {
	client.Client
}

//+kubebuilder:webhook:path=/validate-packages-glasskube-dev-v1alpha1-package,mutating=false,failurePolicy=fail,sideEffects=None,groups=packages.glasskube.dev,resources=packages,verbs=create;update;delete,versions=v1alpha1,name=vpackage.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &PackageValidatingWebhook{}

var ErrInvalidObject = errors.New("validator called with unexpected object type")

func (r *PackageValidatingWebhook) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.Package{}).
		WithValidator(r).
		Complete()
}

// ValidateCreate implements admission.CustomValidator.
func (p *PackageValidatingWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)

	if pkg, ok := obj.(*v1alpha1.Package); ok {
		log.Info("validate create", "name", pkg.Name)

		// TODO(user): fill in your validation logic upon object creation.

		return nil, nil
	}

	return nil, ErrInvalidObject
}

// ValidateUpdate implements admission.CustomValidator.
func (p *PackageValidatingWebhook) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)

	if _, ok := oldObj.(*v1alpha1.Package); ok {
		if newPkg, ok := newObj.(*v1alpha1.Package); ok {
			log.Info("validate update", "name", newPkg.Name)

			// TODO(user): fill in your validation logic upon object update.

			return nil, nil
		}
	}

	return nil, ErrInvalidObject
}

// ValidateDelete implements admission.CustomValidator.
func (p *PackageValidatingWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)

	if pkg, ok := obj.(*v1alpha1.Package); ok {
		log.Info("validate delete", "name", pkg.Name)

		// TODO(user): fill in your validation logic upon object deletion.

		return nil, nil
	}

	return nil, ErrInvalidObject
}

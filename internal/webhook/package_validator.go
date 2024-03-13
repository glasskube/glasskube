package webhook

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/dependency"
	ctrladapter "github.com/glasskube/glasskube/internal/dependency/adapter/controllerruntime"
	"github.com/glasskube/glasskube/internal/repo"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PackageValidatingWebhook struct {
	client.Client
	*owners.OwnerManager
	*dependency.DependendcyManager
	repo repoclient.RepoClient
}

//+kubebuilder:webhook:path=/validate-packages-glasskube-dev-v1alpha1-package,mutating=false,failurePolicy=fail,sideEffects=None,groups=packages.glasskube.dev,resources=packages,verbs=create;update;delete,versions=v1alpha1,name=vpackage.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &PackageValidatingWebhook{}

func (p *PackageValidatingWebhook) SetupWithManager(mgr ctrl.Manager) error {
	p.OwnerManager = owners.NewOwnerManager(p.Scheme())
	p.DependendcyManager = dependency.NewDependencyManager(ctrladapter.NewControllerRuntimeAdapter(p.Client))
	p.repo = repo.DefaultClient
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.Package{}).
		WithValidator(p).
		Complete()
}

// ValidateCreate implements admission.CustomValidator.
func (p *PackageValidatingWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)
	if pkg, ok := obj.(*v1alpha1.Package); ok {
		log.Info("validate create", "name", pkg.Name)
		return nil, p.validateCreateOrUpdate(ctx, pkg)
	}
	return nil, ErrInvalidObject
}

// ValidateUpdate implements admission.CustomValidator.
func (p *PackageValidatingWebhook) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)
	if oldPkg, ok := oldObj.(*v1alpha1.Package); ok {
		if newPkg, ok := newObj.(*v1alpha1.Package); ok {
			if oldPkg.Spec.PackageInfo == newPkg.Spec.PackageInfo {
				// If the package info did not change, we are already done
				return nil, nil
			} else {
				log.Info("validate update", "name", newPkg.Name)
				return nil, p.validateCreateOrUpdate(ctx, newPkg)
			}
		}
	}
	return nil, ErrInvalidObject
}

// ValidateDelete implements admission.CustomValidator.
//
// A package can not be deleted if another package depends on it. This is the case if it has owner references
// to at least one other package. However, we do have to allow packages to be deleted if all their dependant
// packages are already being deleted, because the webhook is also called if a package is garbage-collected and
// this would lead to a dead lock otherwise:
// The dependant package is not deleted because of the foreground deletion finalizer and the dependency is not
// deleted because the dependant is not deleted.
func (p *PackageValidatingWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)
	if pkg, ok := obj.(*v1alpha1.Package); ok {
		log.Info("validate delete", "name", pkg.Name)
		if dependants, err := p.OwnersOfType(&v1alpha1.Package{}, pkg); err != nil {
			return nil, err
		} else if len(dependants) > 0 {
			for _, dep := range dependants {
				var depPkg v1alpha1.Package
				if err := p.Get(ctx, types.NamespacedName{Name: dep.Name}, &depPkg); err != nil {
					return nil, err
				} else if depPkg.DeletionTimestamp.IsZero() {
					return nil, newConflictErrorDelete(dependants[0])
				}
			}
		}
		return nil, nil
	}
	return nil, ErrInvalidObject
}

func (p *PackageValidatingWebhook) validateCreateOrUpdate(ctx context.Context, pkg *v1alpha1.Package) error {
	// We must expect that this package is not installed in this version, so the PackageInfo does not exist.
	var manifest v1alpha1.PackageManifest
	err := p.repo.FetchPackageManifest("", pkg.Spec.PackageInfo.Name, pkg.Spec.PackageInfo.Version, &manifest)
	if err != nil {
		return err
	}

	if result, err := p.Validate(ctx, pkg, &manifest); err != nil {
		return err
	} else if len(result.Conflicts) > 0 {
		// Conflicts are not allowed.
		return newConflictError(result.Conflicts)
	} else if len(result.Requirements) > 0 {
		// Transitive dependencies are not supported yet, so we validate that a required package does not have any
		// dependencies itself.
		// TODO: Add support for validating transitive dependencies.
		var errs error
		for _, req := range result.Requirements {
			var depManifest v1alpha1.PackageManifest
			err := p.repo.FetchPackageManifest("", req.Name, req.Version, &depManifest)
			if err != nil {
				errs = multierr.Append(errs, err)
			} else if len(depManifest.Dependencies) > 0 {
				errs = multierr.Append(errs, newTransitiveError(req, depManifest.Dependencies[0]))
			}
		}
		return errs
	} else {
		return nil
	}
}

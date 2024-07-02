package webhook

import (
	"context"
	"errors"
	"reflect"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/dependency/graph"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PackageValidatingWebhook struct {
	client.Client
	*owners.OwnerManager
	*dependency.DependendcyManager
	RepoClient repoclient.RepoClientset
}

//+kubebuilder:webhook:path=/validate-packages-glasskube-dev-v1alpha1-package,mutating=false,failurePolicy=fail,sideEffects=None,groups=packages.glasskube.dev,resources=packages,verbs=create;update;delete,versions=v1alpha1,name=vpackage.kb.io,admissionReviewVersions=v1
//+kubebuilder:webhook:path=/validate-packages-glasskube-dev-v1alpha1-clusterpackage,mutating=false,failurePolicy=fail,sideEffects=None,groups=packages.glasskube.dev,resources=clusterpackages,verbs=create;update;delete,versions=v1alpha1,name=vclusterpackage.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &PackageValidatingWebhook{}

func (p *PackageValidatingWebhook) SetupWithManager(mgr ctrl.Manager) error {
	p.OwnerManager = owners.NewOwnerManager(p.Scheme())
	return multierr.Combine(
		ctrl.NewWebhookManagedBy(mgr).
			For(&v1alpha1.Package{}).
			WithValidator(p).
			Complete(),
		ctrl.NewWebhookManagedBy(mgr).
			For(&v1alpha1.ClusterPackage{}).
			WithValidator(p).
			Complete(),
	)

}

// ValidateCreate implements admission.CustomValidator.
func (p *PackageValidatingWebhook) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)
	if pkg, ok := obj.(ctrlpkg.Package); ok {
		log.Info("validate create", "name", pkg.GetName())
		return nil, p.validateCreateOrUpdate(ctx, pkg)
	}
	return nil, ErrInvalidObject
}

// ValidateUpdate implements admission.CustomValidator.
func (p *PackageValidatingWebhook) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (warnings admission.Warnings, err error) {
	log := ctrl.LoggerFrom(ctx)
	if oldPkg, ok := oldObj.(ctrlpkg.Package); ok {
		if newPkg, ok := newObj.(ctrlpkg.Package); ok {
			log.Info("validate update", "name", newPkg.GetName())
			if reflect.DeepEqual(oldPkg.GetSpec(), newPkg.GetSpec()) {
				// If the package info did not change, we are already done
				return nil, nil
			} else {
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
	if pkg, ok := obj.(ctrlpkg.Package); ok {
		log.Info("validate delete", "name", pkg.GetName())
		if g, err := p.NewGraph(ctx); err != nil {
			return nil, err
		} else if !pkg.IsNamespaceScoped() {
			// deletion is only validated for cluster-scoped packages
			if _, err := g.ValidateDelete(pkg.GetName()); err != nil {
				for _, err1 := range multierr.Errors(err) {
					if !errors.Is(err1, &graph.DependencyError{}) {
						return nil, err1
					}
				}
				return nil, newConflictErrorDelete(err)
			}
		}
		return nil, nil
	}
	return nil, ErrInvalidObject
}

func (p *PackageValidatingWebhook) validateCreateOrUpdate(ctx context.Context, pkg ctrlpkg.Package) error {
	// We must expect that this package is not installed in this version, so the PackageInfo does not exist.
	var manifest v1alpha1.PackageManifest
	err := p.RepoClient.ForPackage(pkg).FetchPackageManifest(
		pkg.GetSpec().PackageInfo.Name,
		pkg.GetSpec().PackageInfo.Version,
		&manifest,
	)
	if err != nil {
		return err
	}

	if err := manifestvalues.ValidatePackage(manifest, pkg); err != nil {
		return err
	}

	if result, err := p.Validate(ctx, &manifest, pkg.GetSpec().PackageInfo.Version); err != nil {
		return err
	} else if len(result.Conflicts) > 0 {
		// Conflicts are not allowed.
		return newConflictError(result.Conflicts)
	} else {
		return nil
	}
}

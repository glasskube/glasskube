package plain

import (
	"context"
	"fmt"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/controller/owners"
	ownerutils "github.com/glasskube/glasskube/internal/controller/owners/utils"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var fieldOwner = client.FieldOwner("packages.glasskube.dev/package-controller")

type Adapter struct {
	client.Client
	*owners.OwnerManager
}

func NewAdapter() manifest.ManifestAdapter {
	return &Adapter{}
}

// ControllerInit implements manifest.ManifestAdapter.
func (a *Adapter) ControllerInit(builder *builder.Builder, client client.Client, scheme *runtime.Scheme) error {
	if a.OwnerManager == nil {
		a.OwnerManager = owners.NewOwnerManager(scheme)
	}
	a.Client = client
	return nil
}

// Reconcile implements manifest.ManifestAdapter.
func (a *Adapter) Reconcile(
	ctx context.Context,
	pkg *packagesv1alpha1.Package,
	manifest *packagesv1alpha1.PackageManifest,
) (*result.ReconcileResult, error) {
	var allOwned []packagesv1alpha1.OwnedResourceRef
	for _, m := range manifest.Manifests {
		if owned, err := a.reconcilePlainManifest(ctx, *pkg, m); err != nil {
			return nil, err
		} else {
			allOwned = append(allOwned, owned...)
		}
	}
	return result.Ready(fmt.Sprintf("%v manifests reconciled", len(allOwned)), allOwned), nil
}

func (r *Adapter) reconcilePlainManifest(
	ctx context.Context,
	pkg packagesv1alpha1.Package,
	manifest packagesv1alpha1.PlainManifest,
) ([]packagesv1alpha1.OwnedResourceRef, error) {
	log := ctrl.LoggerFrom(ctx)
	objectsToApply, err := clientutils.FetchResources(manifest.Url)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("fetched "+manifest.Url, "objectCount", len(*objectsToApply))

	ownedResources := make([]packagesv1alpha1.OwnedResourceRef, 0, len(*objectsToApply))

	// TODO: check if namespace is terminating before applying

	for _, obj := range *objectsToApply {
		if err := r.SetOwner(&pkg, &obj, owners.BlockOwnerDeletion); err != nil {
			return nil, fmt.Errorf("could set owner reference: %w", err)
		}
		if err := r.Patch(ctx, &obj, client.Apply, fieldOwner, client.ForceOwnership); err != nil {
			return nil, fmt.Errorf("could not apply resource: %w", err)
		}
		log.V(1).Info("applied resource",
			"kind", obj.GroupVersionKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
		if _, err := ownerutils.AddOwnedResourceRef(r.Scheme(), &ownedResources, &obj); err != nil {
			return nil, err
		}
	}
	return ownedResources, nil
}

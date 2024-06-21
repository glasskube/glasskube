package plain

import (
	"context"
	"fmt"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/controller/owners"
	ownerutils "github.com/glasskube/glasskube/internal/controller/owners/utils"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var fieldOwner = client.FieldOwner("packages.glasskube.dev/package-controller")

type Adapter struct {
	client.Client
	*owners.OwnerManager
	namespaceGVK schema.GroupVersionKind
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
	if nsGVK, err := client.GroupVersionKindFor(&corev1.Namespace{}); err != nil {
		return err
	} else {
		a.namespaceGVK = nsGVK
	}
	return nil
}

// Reconcile implements manifest.ManifestAdapter.
func (a *Adapter) Reconcile(
	ctx context.Context,
	pkg ctrlpkg.Package,
	manifest *packagesv1alpha1.PackageManifest,
	patches manifestvalues.TargetPatches,
) (*result.ReconcileResult, error) {
	var allOwned []packagesv1alpha1.OwnedResourceRef
	for _, m := range manifest.Manifests {
		if owned, err := a.reconcilePlainManifest(ctx, pkg, *manifest, m, patches); err != nil {
			return nil, err
		} else {
			allOwned = append(allOwned, owned...)
		}
	}
	return result.Ready(fmt.Sprintf("%v manifests reconciled", len(allOwned)), allOwned), nil
}

func (r *Adapter) reconcilePlainManifest(
	ctx context.Context,
	pkg ctrlpkg.Package,
	pkgManifest packagesv1alpha1.PackageManifest,
	manifest packagesv1alpha1.PlainManifest,
	patches manifestvalues.TargetPatches,
) ([]packagesv1alpha1.OwnedResourceRef, error) {
	log := ctrl.LoggerFrom(ctx)
	var objectsToApply []client.Object
	if unstructured, err := clientutils.FetchResources(manifest.Url); err != nil {
		return nil, err
	} else {
		// Unstructured implements client.Object but we need it as a reference so the interface is fulfilled.
		objectsToApply = make([]client.Object, len(unstructured))
		for i := range unstructured {
			objectsToApply[i] = &unstructured[i]
		}
	}

	log.V(1).Info("fetched "+manifest.Url, "objectCount", len(objectsToApply))

	if pkg.IsNamespaceScoped() {
		for _, obj := range objectsToApply {
			if isNamespaced, err := r.IsObjectNamespaced(obj); err != nil {
				return nil, err
			} else if isNamespaced {
				obj.SetNamespace(pkg.GetNamespace())
			}
		}
	} else {
		// Determine the name of the default namespace. The more specific name takes precedence
		defaultNamespaceName := pkgManifest.DefaultNamespace
		if len(manifest.DefaultNamespace) > 0 {
			defaultNamespaceName = manifest.DefaultNamespace
		}

		if len(defaultNamespaceName) > 0 {
			defaultNamespaceRequired := false
			defaultNamespaceInList := false

			// Determine if the default namespace is needed (at least one resource would be created in the default namespace)
			// and set the namespace property to the default namespace for all objects that are known to be namespaced and do
			// not have an explicit namespace.
			// Also, determine if a namespace resource with the default namespace name already exists in the list.
			for _, obj := range objectsToApply {
				if obj.GetObjectKind().GroupVersionKind() == r.namespaceGVK && obj.GetName() == defaultNamespaceName {
					defaultNamespaceInList = true
				} else {
					if isNamespaced, err := r.IsObjectNamespaced(obj); err != nil {
						// It can not be determined whether this obj kind is namespaced.
						// This can happen if the obj kind is a CRD or some other type that the client does not know.
						// TODO: Should we assume it is namespaced or not? Or just throw an error?
						return nil, err
					} else if isNamespaced {
						if obj.GetNamespace() == "" {
							obj.SetNamespace(defaultNamespaceName)
						}

						if obj.GetNamespace() == defaultNamespaceName {
							defaultNamespaceRequired = true
						}
					}
				}
			}

			if defaultNamespaceRequired && !defaultNamespaceInList {
				defaultNamespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: defaultNamespaceName}}
				// It is necessary to set the GVK manually on this namespace.
				// This could be because we use SSA here.
				// TODO: Find out why!
				defaultNamespace.SetGroupVersionKind(r.namespaceGVK)
				objectsToApply = append([]client.Object{&defaultNamespace}, objectsToApply...)
			}
		}
	}

	// TODO: check if namespace is terminating before applying

	// Apply any modifications before changing anything on the cluster
	for _, obj := range objectsToApply {
		if err := r.SetOwnerIfManagedOrNotExists(r.Client, ctx, pkg, obj); err != nil {
			return nil, err
		}
		if err := patches.ApplyToResource(obj); err != nil {
			return nil, err
		}
	}

	ownedResources := make([]packagesv1alpha1.OwnedResourceRef, 0, len(objectsToApply))
	for _, obj := range objectsToApply {
		if err := r.Patch(ctx, obj, client.Apply, fieldOwner, client.ForceOwnership); err != nil {
			return nil, fmt.Errorf("could not apply resource: %w", err)
		}
		log.V(1).Info("applied resource",
			"kind", obj.GetObjectKind().GroupVersionKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
		if _, err := ownerutils.AddOwnedResourceRef(r.Scheme(), &ownedResources, obj); err != nil {
			return nil, err
		}
	}
	return ownedResources, nil
}

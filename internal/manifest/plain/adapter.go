package plain

import (
	"context"
	"fmt"
	"strings"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/constants"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/controller/owners"
	ownerutils "github.com/glasskube/glasskube/internal/controller/owners/utils"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/resourcepatch"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var fieldOwner = client.FieldOwner("packages.glasskube.dev/package-controller")

type Adapter struct {
	client.Client
	repo repoclient.RepoClientset
	*owners.OwnerManager
	namespaceGVK schema.GroupVersionKind
}

func NewAdapter() manifest.ManifestAdapter {
	return &Adapter{}
}

// ControllerInit implements manifest.ManifestAdapter.
func (a *Adapter) ControllerInit(
	builder *builder.Builder,
	client client.Client,
	repo repoclient.RepoClientset,
	scheme *runtime.Scheme,
) error {
	if a.OwnerManager == nil {
		a.OwnerManager = owners.NewOwnerManager(scheme)
	}
	a.Client = client
	a.repo = repo
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
	pi *packagesv1alpha1.PackageInfo,
	patches resourcepatch.TargetPatches,
) (*result.ReconcileResult, error) {
	var allOwned []packagesv1alpha1.OwnedResourceRef
	for _, manifest := range pi.Status.Manifest.Manifests {
		if owned, err := a.reconcilePlainManifest(ctx, pkg, pi, manifest, patches); err != nil {
			return nil, err
		} else {
			allOwned = append(allOwned, owned...)
		}
	}

	var notReady []packagesv1alpha1.OwnedResourceRef
	var notReadyNames []string
	// of all owned deployments and stateful sets, check their readiness
	for _, ownedResourceRef := range allOwned {
		namespacedName := types.NamespacedName{Namespace: ownedResourceRef.Namespace, Name: ownedResourceRef.Name}
		switch ownedResourceRef.Kind {
		case constants.Deployment:
			deployment := appsv1.Deployment{}
			if err := a.Get(ctx, namespacedName, &deployment); err != nil {
				return nil, fmt.Errorf("failed to get Deployment %v for status check: %w", namespacedName, err)
			}
			if !isReady(deployment.Status.ReadyReplicas, deployment.Spec.Replicas) {
				notReady = append(notReady, ownedResourceRef)
				notReadyNames = append(notReadyNames, namespacedName.String())
			}
		case constants.StatefulSet:
			statefulSet := appsv1.StatefulSet{}
			if err := a.Get(ctx, namespacedName, &statefulSet); err != nil {
				return nil, fmt.Errorf("failed to get StatefulSet for status check: %w", err)
			}
			if !isReady(statefulSet.Status.ReadyReplicas, statefulSet.Spec.Replicas) {
				notReady = append(notReady, ownedResourceRef)
				notReadyNames = append(notReadyNames, namespacedName.String())
			}
		}
	}

	if len(notReady) > 0 {
		return result.Waiting(fmt.Sprintf("%v resources not ready: %v", len(notReady),
			strings.Join(notReadyNames, ",")), allOwned), nil
	} else {
		return result.Ready(fmt.Sprintf("%v manifests reconciled", len(allOwned)), allOwned), nil
	}
}

func isReady(readyReplicas int32, specReplicas *int32) bool {
	return (specReplicas != nil && readyReplicas == *specReplicas) || (specReplicas == nil && readyReplicas > 0)
}

func (r *Adapter) reconcilePlainManifest(
	ctx context.Context,
	pkg ctrlpkg.Package,
	pi *packagesv1alpha1.PackageInfo,
	manifest packagesv1alpha1.PlainManifest,
	patches resourcepatch.TargetPatches,
) ([]packagesv1alpha1.OwnedResourceRef, error) {
	log := ctrl.LoggerFrom(ctx)
	var objectsToApply []client.Object
	if request, err := r.newManifestRequest(pi, manifest.Url); err != nil {
		return nil, err
	} else if unstructured, err := clientutils.FetchResources(request); err != nil {
		return nil, err
	} else {
		// Unstructured implements client.Object but we need it as a reference so the interface is fulfilled.
		objectsToApply = make([]client.Object, len(unstructured))
		for i := range unstructured {
			objectsToApply[i] = &unstructured[i]
		}
		log.V(1).Info("fetched manifest resources",
			"url", request.URL.Redacted(),
			"objectCount", len(objectsToApply))
	}

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
		defaultNamespaceName := pi.Status.Manifest.DefaultNamespace
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

	specHash, specHashErr := pkg.GetSpec().Hashed()
	if specHashErr != nil {
		log.Error(specHashErr, "failed to get spec hash for package – restarts might not happen", "package", pkg)
	}

	// TODO: check if namespace is terminating before applying
	// Apply any modifications before changing anything on the cluster
	for _, obj := range objectsToApply {
		if specHashErr == nil {
			if err := r.annotateWithSpecHash(obj, specHash); err != nil {
				log.Error(err, "could not annotate object with spec hash", "package", pkg, "object", obj)
			}
		}
		if err := r.SetOwnerIfManagedOrNotExists(r.Client, ctx, pkg, obj); err != nil {
			return nil, err
		}
		if err := patches.ApplyToResource(obj); err != nil {
			return nil, err
		}
	}

	if objs, err := prefixAndUpdateReferences(pkg, pi.Status.Manifest, objectsToApply); err != nil {
		return nil, err
	} else {
		objectsToApply = objs
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

// if the obj kind is Deployment or StatefulSet annotateWithSpecHash sets the AnnotationPackageSpecHashed annotation of the
// template to the given specHash. For any other kind it does nothing. Updating the template's annotation to a
// different value than the existing one, will trigger a rolling restart of the resource. When the value stays the same,
// the resource will not be restarted.
func (r *Adapter) annotateWithSpecHash(obj client.Object, specHash string) error {
	switch obj.GetObjectKind().GroupVersionKind().Kind {
	case constants.Deployment, constants.StatefulSet:
		if unstructuredObj, ok := obj.(runtime.Unstructured); ok {
			objContent := unstructuredObj.UnstructuredContent()
			err := unstructured.SetNestedField(objContent, specHash,
				"spec", "template", "metadata", "annotations", packagesv1alpha1.AnnotationPackageSpecHashed)
			if err != nil {
				return err
			}
			unstructuredObj.SetUnstructuredContent(objContent)
		}
	}
	return nil
}

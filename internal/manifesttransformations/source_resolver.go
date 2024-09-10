package manifesttransformations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceResolver interface {
	GetResource(ctx context.Context, ref corev1.TypedObjectReference) (*unstructured.Unstructured, error)
}

type ctrlResourceResolver struct {
	client client.Client
}

func (r *ctrlResourceResolver) GetResource(
	ctx context.Context,
	ref corev1.TypedObjectReference,
) (*unstructured.Unstructured, error) {
	var obj unstructured.Unstructured
	obj.SetName(ref.Name)
	obj.SetNamespace(*ref.Namespace)
	if gv, err := schema.ParseGroupVersion(*ref.APIGroup); err != nil {
		return nil, err
	} else {
		obj.SetGroupVersionKind(gv.WithKind(ref.Kind))
	}
	if err := r.client.Get(ctx, client.ObjectKeyFromObject(&obj), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

type SourceResolver struct {
	resourceResolver ResourceResolver
}

func NewResolver(client client.Client) *SourceResolver {
	return &SourceResolver{
		&ctrlResourceResolver{client: client},
	}
}

func refWithNamespace(ref corev1.TypedLocalObjectReference, namespace string) corev1.TypedObjectReference {
	return corev1.TypedObjectReference{
		APIGroup:  ref.APIGroup,
		Kind:      ref.Kind,
		Name:      ref.Name,
		Namespace: &namespace,
	}
}

func (resolver *SourceResolver) Resolve(
	ctx context.Context,
	pkg ctrlpkg.Package,
	source v1alpha1.TransformationSource,
) (any, error) {
	jp := jsonpath.New("")
	if err := jp.Parse(source.Path); err != nil {
		return nil, err
	}

	var resource *unstructured.Unstructured
	if source.Resource != nil {
		ref := refWithNamespace(*source.Resource, pkg.GetNamespace())
		if r, err := resolver.resourceResolver.GetResource(ctx, ref); err != nil {
			return nil, err
		} else {
			resource = r
		}
	} else {
		resource = &unstructured.Unstructured{}
		if data, err := json.Marshal(pkg); err != nil {
			return nil, err
		} else if err = json.Unmarshal(data, resource); err != nil {
			return nil, err
		}
	}

	if results, err := jp.FindResults(resource.UnstructuredContent()); err != nil {
		return nil, err
	} else if len(results) != 1 {
		return nil, fmt.Errorf("jsonpath produced unexpected number of results (%v)", len(results))
	} else if len(results[0]) != 1 {
		return nil, fmt.Errorf("jsonpath produced unexpected number of subresults (%v)", len(results[0]))
	} else {
		return results[0][0].Interface(), nil
	}
}

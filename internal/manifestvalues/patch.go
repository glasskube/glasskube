package manifestvalues

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"text/template"

	jsonpatch "github.com/evanphx/json-patch/v5"
	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/maputils"
	corev1 "k8s.io/api/core/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// object combines the runtime and metav1 object interfaces.
// This is the same as the controller-runtime client Object but we don't want any controller-runtime dependencies
// in packages that are imported by the client.
type object interface {
	runtime.Object
	metav1.Object
}

type patchWithValue struct {
	v1alpha1.PartialJsonPatch `json:",inline"`
	Value                     any `json:"value"`
}

type targetResource struct {
	schema.GroupVersionKind
	name      string
	namespace *string
}

func generateTargetResource(ref *corev1.TypedObjectReference) (*targetResource, error) {
	var groupVersion schema.GroupVersion
	if ref.APIGroup != nil {
		if gv, err := schema.ParseGroupVersion(*ref.APIGroup); err != nil {
			return nil, err
		} else {
			groupVersion = gv
		}
	}
	return &targetResource{
		GroupVersionKind: groupVersion.WithKind(ref.Kind),
		name:             ref.Name,
		namespace:        ref.Namespace,
	}, nil
}

func (r *targetResource) Match(obj object) bool {
	return r != nil &&
		r.GroupVersionKind == obj.GetObjectKind().GroupVersionKind() &&
		r.name == obj.GetName() &&
		(r.namespace == nil || *r.namespace == obj.GetNamespace())
}

type TargetPatch struct {
	resource  *targetResource
	helmChart *string
	patch     jsonpatch.Patch
}

// generateTargetPatch does three things:
//   - execute the targets valueTemplate if it exists
//   - create an applicable jsonpatch
//   - resolve the targets resource GVK
func generateTargetPatch(target v1alpha1.ValueDefinitionTarget, value string) (*TargetPatch, error) {
	if actualValue, err := getActualValue(target, value); err != nil {
		return nil, err
	} else if jsonPatch, err := generateJsonPatch(target.Patch, actualValue); err != nil {
		return nil, err
	} else {
		newResult := TargetPatch{patch: jsonPatch}
		if target.Resource != nil {
			if resource, err := generateTargetResource(target.Resource); err != nil {
				return nil, err
			} else {
				newResult.resource = resource
			}
		} else if target.ChartName != nil {
			newResult.helmChart = target.ChartName
		}
		return &newResult, nil
	}
}

func getActualValue(target v1alpha1.ValueDefinitionTarget, value string) (any, error) {
	if len(target.ValueTemplate) == 0 {
		return value, nil
	}

	tmplBase := template.New("").Funcs(template.FuncMap{
		"base64": func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) },
	})

	if tmpl, err := tmplBase.Parse(target.ValueTemplate); err != nil {
		return nil, err
	} else {
		var bw bytes.Buffer
		if err := tmpl.Execute(&bw, value); err != nil {
			return nil, err
		}
		var result any
		if err := json.Unmarshal(bw.Bytes(), &result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

func generateJsonPatch(p v1alpha1.PartialJsonPatch, value any) (jsonpatch.Patch, error) {
	// jsonpatch works with json.RawMessage, so the patch must be converted to JSON first.
	if data, err := json.Marshal([]patchWithValue{{p, value}}); err != nil {
		return nil, err
	} else {
		return jsonpatch.DecodePatch(data)
	}
}

func (p *TargetPatch) MatchResource(obj object) bool {
	return p.resource.Match(obj)
}

func (p *TargetPatch) MatchHelmRelease(obj *helmv2.HelmRelease) bool {
	return p.helmChart != nil && *p.helmChart == obj.Spec.Chart.Spec.Chart
}

func (p *TargetPatch) ApplyToResource(obj object) error {
	if p.MatchResource(obj) {
		return p.apply(obj)
	} else {
		return nil
	}
}

func (p *TargetPatch) ApplyToHelmRelease(obj *helmv2.HelmRelease) error {
	if p.MatchHelmRelease(obj) {
		if obj.Spec.Values == nil || len(obj.Spec.Values.Raw) == 0 {
			obj.Spec.Values = &extv1.JSON{Raw: []byte("{}")}
		}
		return p.apply(obj.Spec.Values)
	} else {
		return nil
	}
}

func (p *TargetPatch) apply(obj any) error {
	if data, err := json.Marshal(obj); err != nil {
		return err
	} else if patched, err := p.patch.Apply(data); err != nil {
		return err
	} else if err := json.Unmarshal(patched, &obj); err != nil {
		return err
	} else {
		return nil
	}
}

type TargetPatches []TargetPatch

func (patches *TargetPatches) ApplyToResource(obj object) error {
	for _, patch := range *patches {
		if err := patch.ApplyToResource(obj); err != nil {
			return err
		}
	}
	return nil
}

func (patches *TargetPatches) ApplyToHelmRelease(obj *helmv2.HelmRelease) error {
	for _, patch := range *patches {
		if err := patch.ApplyToHelmRelease(obj); err != nil {
			return err
		}
	}
	return nil
}

// GeneratePatches creates an applicable patch for each value definition in the supplied manifest that a
// value in values exists for. It performs no validation of the supplied values. Please run validate
// before using this!
func GeneratePatches(manifest v1alpha1.PackageManifest, values map[string]string) (TargetPatches, error) {
	var result []TargetPatch
	for _, name := range maputils.KeysSorted(manifest.ValueDefinitions) {
		def := manifest.ValueDefinitions[name]
		if value, ok := values[name]; ok {
			for _, target := range def.Targets {
				if patch, err := generateTargetPatch(target, value); err != nil {
					return nil, err
				} else {
					result = append(result, *patch)
				}
			}
		}
	}
	return result, nil
}

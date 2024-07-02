package manifestvalues

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch/v5"
	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func jsonPatch(s string) jsonpatch.Patch {
	p, err := jsonpatch.DecodePatch([]byte("[" + s + "]"))
	Expect(err).NotTo(HaveOccurred())
	return p
}

func newUnstructured(s string) unstructured.Unstructured {
	var u unstructured.Unstructured
	err := json.Unmarshal([]byte(s), &u)
	Expect(err).NotTo(HaveOccurred())
	return u
}

var (
	foo         = "foo"
	appsv1group = "apps/v1"
)

var _ = Describe("GeneratePatches", func() {
	It("should skip patch for missing value", func() {
		patches, err := GeneratePatches(
			v1alpha1.PackageManifest{
				ValueDefinitions: map[string]v1alpha1.ValueDefinition{
					"foo": {Targets: []v1alpha1.ValueDefinitionTarget{{
						ChartName: &foo,
						Patch:     v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"},
					}}},
				}},
			map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		Expect(patches).To(BeEmpty())
	})
	It("should pick correct value", func() {
		patches, err := GeneratePatches(
			v1alpha1.PackageManifest{
				ValueDefinitions: map[string]v1alpha1.ValueDefinition{
					"foo": {Targets: []v1alpha1.ValueDefinitionTarget{{
						ChartName: &foo,
						Patch:     v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"},
					}}},
				}},
			map[string]string{"foo": "test"})
		Expect(err).NotTo(HaveOccurred())
		Expect(patches).To(ConsistOf(TargetPatch{
			helmChart: &foo,
			patch:     jsonPatch(`{"op":"add","path":"/spec/replicas","value":"test"}`),
		}))
	})
})

var _ = Describe("generateTargetPatch", func() {
	DescribeTable("generateTargetPatch",
		func(target v1alpha1.ValueDefinitionTarget, value string, expectError bool, resource *targetResource,
			chartName *string, patch string) {
			result, err := generateTargetPatch(target, value)
			if expectError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
			if resource != nil {
				Expect(result.resource).NotTo(BeNil())
				Expect(*result.resource).To(Equal(*resource))
			} else {
				Expect(result.resource).To(BeNil())
			}
			if chartName != nil {

				Expect(result.helmChart).NotTo(BeNil())
				Expect(*result.helmChart).To(Equal(*chartName))
			}
			Expect(result.patch).To(Equal(jsonPatch(patch)))
		},
		Entry("when resource present",
			v1alpha1.ValueDefinitionTarget{
				Resource: &corev1.TypedObjectReference{APIGroup: &appsv1group, Kind: "Deployment", Name: "foo", Namespace: &foo},
				Patch:    v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"},
			},
			"2",
			false,
			&targetResource{schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, "foo", &foo},
			nil,
			`{"op":"add","path":"/spec/replicas","value":"2"}`,
		),
		Entry("when chart name present",
			v1alpha1.ValueDefinitionTarget{ChartName: &foo, Patch: v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"}},
			"2",
			false,
			nil,
			&foo,
			`{"op":"add","path":"/spec/replicas","value":"2"}`,
		),
		Entry("when valueTemplate is a number",
			v1alpha1.ValueDefinitionTarget{
				ValueTemplate: "{{.}}",
				Patch:         v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"},
			},
			"2",
			false,
			nil,
			nil,
			`{"op":"add","path":"/spec/replicas","value":2}`,
		),
		Entry("when valueTemplate is an object",
			v1alpha1.ValueDefinitionTarget{
				ValueTemplate: `{{if eq . "2"}} {"foo":{{.}}} {{else}} {"bar":{{.}}} {{end}}`,
				Patch:         v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"},
			},
			"2",
			false,
			nil,
			nil,
			`{"op":"add","path":"/spec/replicas","value":{"foo":2}}`,
		),
	)
})

var _ = Describe("ApplyToResource", func() {
	It("should patch resource", func() {
		patch, err := generateTargetPatch(v1alpha1.ValueDefinitionTarget{
			Resource: &corev1.TypedObjectReference{APIGroup: &appsv1group, Kind: "Deployment", Name: "foo", Namespace: &foo},
			Patch:    v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"},
		}, "2")
		Expect(err).NotTo(HaveOccurred())

		obj := newUnstructured(
			`{"apiVersion": "apps/v1",
			"kind": "Deployment",
			"metadata": {"name": "foo", "namespace": "foo"},
			"spec":{}}`)
		expected := newUnstructured(
			`{"apiVersion": "apps/v1",
			"kind": "Deployment",
			"metadata": {"name": "foo", "namespace": "foo"},
			"spec": {"replicas": "2"}}`)

		Expect(patch.MatchResource(&obj)).To(BeTrue())
		Expect(patch.ApplyToResource(&obj)).NotTo(HaveOccurred())
		Expect(obj).To(Equal(expected))
	})
	It("should patch helm values", func() {
		patch, err := generateTargetPatch(v1alpha1.ValueDefinitionTarget{
			ChartName: &foo,
			Patch:     v1alpha1.PartialJsonPatch{Op: "add", Path: "/spec/replicas"},
		}, "2")
		Expect(err).NotTo(HaveOccurred())

		obj := helmv2.HelmRelease{
			TypeMeta: metav1.TypeMeta{
				APIVersion: helmv2.GroupVersion.Version + "/" + helmv2.GroupVersion.Version,
				Kind:       "HelmRelease"},
			ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foo"},
			Spec: helmv2.HelmReleaseSpec{
				Chart:  &helmv2.HelmChartTemplate{Spec: helmv2.HelmChartTemplateSpec{Chart: "foo"}},
				Values: &extv1.JSON{Raw: []byte(`{"spec":{}}`)}}}
		expected := *obj.DeepCopy()
		expected.Spec.Values.Raw = []byte(`{"spec":{"replicas":"2"}}`)

		Expect(patch.MatchHelmRelease(&obj)).To(BeTrue())
		Expect(patch.ApplyToHelmRelease(&obj)).NotTo(HaveOccurred())
		Expect(obj).To(Equal(expected))
	})
})

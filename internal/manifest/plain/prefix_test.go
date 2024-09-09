package plain

import (
	"bytes"
	"io"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var pkg = &v1alpha1.Package{
	ObjectMeta: metav1.ObjectMeta{Name: "foo"},
	Spec: v1alpha1.PackageSpec{
		PackageInfo: v1alpha1.PackageInfoTemplate{
			Name: "test",
		},
	},
}

var secretAndDeployment = `
apiVersion: v1
kind: Secret
metadata:
  name: test
---
apiVersion: v1
kind: Deployment
metadata:
  name: test
spec:
  template:
    spec:
      containers:
        - env:
          - valueFrom:
              secretKeyRef:
                name: test
          - valueFrom:
              secretKeyRef:
                name: test1
`
var secretAndDeploymentExpected = `
apiVersion: v1
kind: Secret
metadata:
  name: foo-test
  labels:
    packages.glasskube.dev/package: test
    packages.glasskube.dev/instance: foo
---
apiVersion: v1
kind: Deployment
metadata:
  name: foo-test
  labels:
    packages.glasskube.dev/package: test
    packages.glasskube.dev/instance: foo
spec:
  selector:
    matchLabels:
      packages.glasskube.dev/package: test
      packages.glasskube.dev/instance: foo
  template:
    metadata:
      labels:
        packages.glasskube.dev/package: test
        packages.glasskube.dev/instance: foo
    spec:
      containers:
        - env:
          - valueFrom:
              secretKeyRef:
                name: foo-test
          - valueFrom:
              secretKeyRef:
                name: test1
`
var secretAndDeploymentExpectedWithTransitive = `
apiVersion: v1
kind: Secret
metadata:
  name: foo-test
  labels:
    packages.glasskube.dev/package: test
    packages.glasskube.dev/instance: foo
---
apiVersion: v1
kind: Deployment
metadata:
  name: foo-test
  labels:
    packages.glasskube.dev/package: test
    packages.glasskube.dev/instance: foo
spec:
  selector:
    matchLabels:
      packages.glasskube.dev/package: test
      packages.glasskube.dev/instance: foo
  template:
    metadata:
      labels:
        packages.glasskube.dev/package: test
        packages.glasskube.dev/instance: foo
    spec:
      containers:
        - env:
          - valueFrom:
              secretKeyRef:
                name: foo-test
          - valueFrom:
              secretKeyRef:
                name: foo-test1
`
var serviceAndIngress = `
apiVersion: v1
kind: Service
metadata:
  name: test
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test
spec:
  rules:
    - http:
        paths:
          - backend:
              service:
                name: test
`
var serviceAndIngressExpected = `
apiVersion: v1
kind: Service
metadata:
  name: foo-test
  labels:
    packages.glasskube.dev/package: test
    packages.glasskube.dev/instance: foo
spec:
  selector:
    packages.glasskube.dev/package: test
    packages.glasskube.dev/instance: foo
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: foo-test
  labels:
    packages.glasskube.dev/package: test
    packages.glasskube.dev/instance: foo
spec:
  rules:
    - http:
        paths:
          - backend:
              service:
                name: foo-test
`

func parseObjs(data string) []client.Object {
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer([]byte(data)), 4096)
	var objs []client.Object
	for {
		var obj unstructured.Unstructured
		if err := dec.Decode(&obj); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		objs = append(objs, &obj)
	}
	return objs
}

var _ = Describe("updateReferences", func() {
	It("should handle empty list", func() {
		objs := []client.Object{}
		newObjs, err := prefixAndUpdateReferences(pkg, &v1alpha1.PackageManifest{}, objs)
		Expect(err).NotTo(HaveOccurred())
		Expect(newObjs).To(BeEmpty())
	})

	It("should handle deployment with env", func() {
		objs := parseObjs(secretAndDeployment)
		expectedObj := parseObjs(secretAndDeploymentExpected)
		newObjs, err := prefixAndUpdateReferences(pkg, &v1alpha1.PackageManifest{}, objs)
		Expect(err).NotTo(HaveOccurred())
		Expect(newObjs).To(ConsistOf(expectedObj))
	})

	It("should handle deployment with env and transitive resource", func() {
		objs := parseObjs(secretAndDeployment)
		expectedObj := parseObjs(secretAndDeploymentExpectedWithTransitive)
		newObjs, err := prefixAndUpdateReferences(pkg, &v1alpha1.PackageManifest{
			TransitiveResources: []corev1.TypedLocalObjectReference{
				{
					APIGroup: util.Pointer("v1"),
					Kind:     "Secret",
					Name:     "test1",
				},
			},
		}, objs)
		Expect(err).NotTo(HaveOccurred())
		Expect(newObjs).To(ConsistOf(expectedObj))
	})

	It("should handle ingress", func() {
		objs := parseObjs(serviceAndIngress)
		expectedObjs := parseObjs(serviceAndIngressExpected)
		newObjs, err := prefixAndUpdateReferences(pkg, &v1alpha1.PackageManifest{}, objs)
		Expect(err).NotTo(HaveOccurred())
		Expect(newObjs).To(ConsistOf(expectedObjs))
	})
})

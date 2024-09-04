package plain

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type fakeScopeChecker struct{}

// IsObjectNamespaced implements ScopeChecker.
func (*fakeScopeChecker) IsObjectNamespaced(runtime.Object) (bool, error) {
	return true, nil
}

var p = newPrefixer(&fakeScopeChecker{})
var pkg = &v1alpha1.Package{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}

var _ = Describe("updateReferences", func() {
	It("should handle empty list", func() {
		objs := []client.Object{}
		err := p.prefixAndUpdateReferences(pkg, &v1alpha1.PackageManifest{}, objs)
		Expect(err).NotTo(HaveOccurred())
		Expect(objs).To(BeEmpty())
	})

	It("should handle deployment with env", func() {
		objs := []client.Object{
			&corev1.Secret{
				TypeMeta:   metav1.TypeMeta{Kind: "Secret"},
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
			},
			&appsv1.Deployment{
				TypeMeta:   metav1.TypeMeta{Kind: "Deployment"},
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Env: []corev1.EnvVar{
										{
											ValueFrom: &corev1.EnvVarSource{
												SecretKeyRef: &corev1.SecretKeySelector{
													LocalObjectReference: corev1.LocalObjectReference{Name: "test"},
												},
											},
										},
										{
											ValueFrom: &corev1.EnvVarSource{
												SecretKeyRef: &corev1.SecretKeySelector{
													LocalObjectReference: corev1.LocalObjectReference{Name: "test1"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		err := p.prefixAndUpdateReferences(pkg, &v1alpha1.PackageManifest{}, objs)
		Expect(err).NotTo(HaveOccurred())
		Expect(objs).To(HaveLen(2))
		Expect(objs[0].GetName()).To(Equal("foo-test"))
		Expect(objs[1]).To(Equal(&appsv1.Deployment{
			TypeMeta:   metav1.TypeMeta{Kind: "Deployment"},
			ObjectMeta: metav1.ObjectMeta{Name: "foo-test"},
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Env: []corev1.EnvVar{
									{
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{Name: "foo-test"},
											},
										},
									},
									{
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{Name: "test1"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}))
	})

	It("should handle ingress", func() {
		objs := []client.Object{
			&corev1.Service{
				TypeMeta:   metav1.TypeMeta{Kind: "Service"},
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
			},
			&netv1.Ingress{
				TypeMeta:   metav1.TypeMeta{Kind: "Ingress"},
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "test",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		err := p.prefixAndUpdateReferences(pkg, &v1alpha1.PackageManifest{}, objs)
		Expect(err).NotTo(HaveOccurred())
		Expect(objs).To(HaveLen(2))
		Expect(objs[0].GetName()).To(Equal("foo-test"))
		Expect(objs[1]).To(Equal(&netv1.Ingress{
			TypeMeta:   metav1.TypeMeta{Kind: "Ingress"},
			ObjectMeta: metav1.ObjectMeta{Name: "foo-test"},
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-test",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}))
	})
})

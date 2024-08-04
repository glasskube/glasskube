package manifestvalues

import (
	"context"
	"encoding/base64"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/adapter/controllerruntime"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newTestResolver(initObjs ...runtime.Object) *Resolver {
	scheme := runtime.NewScheme()
	Expect(clientscheme.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(v1alpha1.AddToScheme(scheme)).NotTo(HaveOccurred())
	client := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(initObjs...).
		Build()
	return NewResolver(
		controllerruntime.NewPackageClientAdapter(client),
		controllerruntime.NewKubernetesClientAdapter(client),
	)
}

var _ = Describe("resolver", func() {
	var testConst = "test"

	It("should resolve literal value", func(ctx context.Context) {
		resolver := newTestResolver()
		result, err := resolver.Resolve(ctx, map[string]v1alpha1.ValueConfiguration{
			"test": {Value: &testConst},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(map[string]string{"test": "test"}))
	})

	It("should resolve ConfigMap reference value", func(ctx context.Context) {
		resolver := newTestResolver(&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
			Data:       map[string]string{"test": "test"},
		})
		result, err := resolver.Resolve(ctx, map[string]v1alpha1.ValueConfiguration{
			"test": {
				ValueFrom: &v1alpha1.ValueReference{
					ConfigMapRef: &v1alpha1.ObjectKeyValueSource{Name: "test", Namespace: "test", Key: "test"},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(map[string]string{"test": "test"}))
	})

	It("should resolve Secret reference value", func(ctx context.Context) {
		resolver := newTestResolver(&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
			Data:       map[string][]byte{"test": []byte(base64.StdEncoding.EncodeToString([]byte("test")))},
		})
		result, err := resolver.Resolve(ctx, map[string]v1alpha1.ValueConfiguration{
			"test": {
				ValueFrom: &v1alpha1.ValueReference{
					SecretRef: &v1alpha1.ObjectKeyValueSource{Name: "test", Namespace: "test", Key: "test"},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(map[string]string{"test": "test"}))
	})

	It("should resolve Package reference value", func(ctx context.Context) {
		resolver := newTestResolver(
			&v1alpha1.ClusterPackage{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: v1alpha1.PackageSpec{
					Values: map[string]v1alpha1.ValueConfiguration{
						"test": {
							ValueFrom: &v1alpha1.ValueReference{
								SecretRef: &v1alpha1.ObjectKeyValueSource{Name: "test", Namespace: "test", Key: "test"},
							},
						},
					},
				},
			},
			&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Data:       map[string][]byte{"test": []byte(base64.StdEncoding.EncodeToString([]byte("test")))},
			},
		)
		result, err := resolver.Resolve(ctx, map[string]v1alpha1.ValueConfiguration{
			"test": {
				ValueFrom: &v1alpha1.ValueReference{
					PackageRef: &v1alpha1.PackageValueSource{Name: "test", Value: "test"},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(map[string]string{"test": "test"}))
	})
})

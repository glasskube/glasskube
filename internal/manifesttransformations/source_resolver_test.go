package manifesttransformations

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var ctrlClient client.Client

var resolver SourceResolver

var _ = Describe("Resolve", func() {
	BeforeEach(func() {
		ctrlClient = fake.NewClientBuilder().Build()
		NewResolver(ctrlClient)
	})

	It("should return package name", func(ctx context.Context) {
		pkg := v1alpha1.Package{
			TypeMeta:   metav1.TypeMeta{Kind: "Package", APIVersion: "packages.glasskube.dev/v1alpha1"},
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		}
		src := v1alpha1.TransformationSource{
			Path: "{ $.metadata.name }",
		}
		result, err := resolver.Resolve(ctx, &pkg, src)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("test"))
	})
})

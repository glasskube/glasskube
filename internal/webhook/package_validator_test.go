package webhook

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func newPackageValidatingWebhook(objects ...client.Object) webhook.CustomValidator {
	return &PackageValidatingWebhook{
		Client: fake.NewClientBuilder().
			WithObjects(objects...).
			Build(),
	}
}

// These tests demonstrate how unit tests for a CustomValidator COULD be implemented using the controller-runtime fake client.
// They will likely become useless once the PackageValidatingWebhook does some actual validation.
// TODO: Add some tests for package validation

var _ = Describe("PackageValidatingWebhook", Ordered, func() {
	BeforeAll(func() {
		err := v1alpha1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should ValidateCreate without error", func(ctx context.Context) {
		webhook := newPackageValidatingWebhook()
		_, err := webhook.ValidateCreate(ctx, &v1alpha1.Package{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("should ValidateUpdate without error", func(ctx context.Context) {
		webhook := newPackageValidatingWebhook()
		_, err := webhook.ValidateUpdate(ctx, &v1alpha1.Package{}, &v1alpha1.Package{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("should ValidateDelete without error", func(ctx context.Context) {
		pkg := v1alpha1.Package{}
		webhook := newPackageValidatingWebhook(&pkg)
		_, err := webhook.ValidateDelete(ctx, &pkg)
		Expect(err).NotTo(HaveOccurred())
	})

	When("called with invalid object", func() {
		It("should ValidateCreate with error", func(ctx context.Context) {
			webhook := newPackageValidatingWebhook()
			_, err := webhook.ValidateCreate(ctx, &unstructured.Unstructured{})
			Expect(err).To(HaveOccurred())
		})
	})
})

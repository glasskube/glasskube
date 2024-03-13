package webhook

import (
	"context"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/dependency"
	ctrladapter "github.com/glasskube/glasskube/internal/dependency/adapter/controllerruntime"
	fakerepo "github.com/glasskube/glasskube/internal/repo/client/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var fakeRepo = fakerepo.FakeClient{}

func newPackageValidatingWebhook(objects ...client.Object) *PackageValidatingWebhook {
	fakeClient := fake.NewClientBuilder().
		WithObjects(objects...).
		Build()
	return &PackageValidatingWebhook{
		Client:       fakeClient,
		OwnerManager: owners.NewOwnerManager(scheme.Scheme),
		DependendcyManager: dependency.NewDependencyManager(ctrladapter.NewControllerRuntimeAdapter(fakeClient)).
			WithRepo(&fakeRepo),
		repo: &fakeRepo,
	}
}

var (
	foov1pkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "foo", Version: "v1"}}}
	foov2pkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "foo", Version: "v2"}}}
	barv1pkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "bar"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "bar", Version: "v1"}}}
	barv2pkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "bar"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "bar", Version: "v2"}}}
	bazv1pkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "baz"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "baz", Version: "v1"}}}
	bazv2pkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "baz"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "baz", Version: "v2"}}}
	notExistsPkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "doesnotexist"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "doesnotexist", Version: "v1"}}}
)

// These tests demonstrate how unit tests for a CustomValidator COULD be implemented using the controller-runtime fake client.
// They will likely become useless once the PackageValidatingWebhook does some actual validation.
// TODO: Add some tests for package validation

var _ = Describe("PackageValidatingWebhook", Ordered, func() {
	BeforeAll(func() {
		err := v1alpha1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())
		fakeRepo.Packages = map[string]map[string]v1alpha1.PackageManifest{
			"foo": {
				"v1": {Dependencies: []v1alpha1.Dependency{{Name: "bar", Version: "v1"}}},
				"v2": {Dependencies: []v1alpha1.Dependency{{Name: "bar", Version: "v2"}}},
			},
			"bar": {
				"v1": {},
				"v2": {},
			},
			"baz": {
				"v1": {Dependencies: []v1alpha1.Dependency{{Name: "foo", Version: "v1"}}},
				"v2": {Dependencies: []v1alpha1.Dependency{{Name: "foo", Version: "v2"}}},
			},
		}
	})

	Context("ValidateCreate", func() {
		When("cluster is empty", func() {
			It("should not return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateCreate(ctx, &foov1pkg)
				Expect(err).NotTo(HaveOccurred())
			})
			When("a package with that name does not exist", func() {
				It("should return an error", func(ctx context.Context) {
					webhook := newPackageValidatingWebhook()
					_, err := webhook.ValidateCreate(ctx, &notExistsPkg)
					Expect(err).To(HaveOccurred())
				})
			})
		})
		When("the dependency exists", func() {
			When("dependency version ok", func() {
				It("should not return error", func(ctx context.Context) {
					webhook := newPackageValidatingWebhook(&barv1pkg)
					_, err := webhook.ValidateCreate(ctx,
						&foov1pkg)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			When("there is a version mismatch", func() {
				It("should return conflict error", func(ctx context.Context) {
					webhook := newPackageValidatingWebhook(&barv1pkg)
					_, err := webhook.ValidateCreate(ctx, &foov2pkg)
					Expect(err).To(And(HaveOccurred(), Satisfy(isErrDependencyConflict)))
				})
			})
		})
		When("package has transitive dependency", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateCreate(ctx, &bazv1pkg)
				Expect(err).To(And(HaveOccurred(), Satisfy(isErrTransitiveDependency)))
			})
		})
		When("called with object other than Package", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateCreate(ctx, &unstructured.Unstructured{})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("ValidateUpdate", func() {
		When("cluster is empty", func() {
			It("should not return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateUpdate(ctx, &foov1pkg, &foov2pkg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("the dependency exists", func() {
			When("dependency version ok", func() {
				It("should not return error", func(ctx context.Context) {
					webhook := newPackageValidatingWebhook(&barv2pkg)
					_, err := webhook.ValidateUpdate(ctx,
						&foov1pkg, &foov2pkg)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			When("there is a version mismatch", func() {
				It("should return conflict error", func(ctx context.Context) {
					webhook := newPackageValidatingWebhook(&barv1pkg)
					_, err := webhook.ValidateUpdate(ctx, &foov1pkg, &foov2pkg)
					Expect(err).To(And(HaveOccurred(), Satisfy(isErrDependencyConflict)))
				})
			})
		})
		When("package has transitive dependency", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateUpdate(ctx, &bazv1pkg, &bazv2pkg)
				Expect(err).To(And(HaveOccurred(), Satisfy(isErrTransitiveDependency)))
			})
		})
		When("called with object other than Package", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateCreate(ctx, &unstructured.Unstructured{})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("ValidateDelete", func() {
		When("package has no owner references", func() {
			It("should not return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook(&foov1pkg)
				_, err := webhook.ValidateDelete(ctx, &foov1pkg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("package has owner reference", func() {
			It("should return an error", func(ctx context.Context) {
				barv1Copy := barv1pkg.DeepCopy()
				err := controllerutil.SetOwnerReference(&foov1pkg, barv1Copy, scheme.Scheme)
				Expect(err).NotTo(HaveOccurred())

				webhook := newPackageValidatingWebhook(&foov1pkg, barv1Copy)
				_, err = webhook.ValidateDelete(ctx, barv1Copy)
				Expect(err).To(And(HaveOccurred(), Satisfy(isErrDependencyConflict)))
			})
			When("owning package is also being deleted", func() {
				It("should not return an error", func(ctx context.Context) {
					foov1Copy := foov1pkg.DeepCopy()
					foov1Copy.DeletionTimestamp = &v1.Time{Time: time.Now()}
					foov1Copy.Finalizers = []string{"foregroundDeletion"}
					barv1Copy := barv1pkg.DeepCopy()
					err := controllerutil.SetOwnerReference(foov1Copy, barv1Copy, scheme.Scheme)
					Expect(err).NotTo(HaveOccurred())

					webhook := newPackageValidatingWebhook(foov1Copy, barv1Copy)
					_, err = webhook.ValidateDelete(ctx, barv1Copy)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
})

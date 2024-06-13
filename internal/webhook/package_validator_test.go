package webhook

import (
	"context"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	ctrladapter "github.com/glasskube/glasskube/internal/adapter/controllerruntime"
	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/names"
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

var fakeRepoClient = fakerepo.FakeClient{
	PackageRepositories: []v1alpha1.PackageRepository{{}},
}
var fakeRepoClientset = fakerepo.FakeClientset{Client: &fakeRepoClient}

func newPackageValidatingWebhook(objects ...client.Object) *PackageValidatingWebhook {
	fakeClient := fake.NewClientBuilder().
		WithObjects(objects...).
		WithObjects(&v1alpha1.PackageRepository{}).
		Build()
	ownerManager := owners.NewOwnerManager(scheme.Scheme)
	return &PackageValidatingWebhook{
		Client:       fakeClient,
		OwnerManager: ownerManager,
		DependendcyManager: dependency.NewDependencyManager(
			ctrladapter.NewPackageClientAdapter(fakeClient),
			&fakeRepoClientset,
		),
		RepoClient: &fakeRepoClientset,
	}
}

var (
	foov1pkg = v1alpha1.ClusterPackage{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "foo", Version: "v1"}},
		Status:     v1alpha1.PackageStatus{OwnedPackageInfos: []v1alpha1.OwnedResourceRef{{Name: "foo--v1"}}}}
	foov1pi = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: names.PackageInfoName(&foov1pkg)},
		Spec:       v1alpha1.PackageInfoSpec{Name: "foo", Version: "v1"},
		Status: v1alpha1.PackageInfoStatus{
			Manifest: &v1alpha1.PackageManifest{Name: "foo", Dependencies: []v1alpha1.Dependency{{Name: "bar", Version: "v1"}}}}}
	foov2pkg = v1alpha1.ClusterPackage{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "foo", Version: "v2"}},
		Status:     v1alpha1.PackageStatus{OwnedPackageInfos: []v1alpha1.OwnedResourceRef{{Name: "foo--v2"}}}}
	foov2pi = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: names.PackageInfoName(&foov2pkg)},
		Spec:       v1alpha1.PackageInfoSpec{Name: "foo", Version: "v2"},
		Status: v1alpha1.PackageInfoStatus{
			Manifest: &v1alpha1.PackageManifest{Name: "foo", Dependencies: []v1alpha1.Dependency{{Name: "bar", Version: "v2"}}}}}
	barv1pkg = v1alpha1.ClusterPackage{
		ObjectMeta: v1.ObjectMeta{Name: "bar"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "bar", Version: "v1"}},
		Status:     v1alpha1.PackageStatus{OwnedPackageInfos: []v1alpha1.OwnedResourceRef{{Name: "bar--v1"}}}}
	barv1pi = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: names.PackageInfoName(&barv1pkg)},
		Spec:       v1alpha1.PackageInfoSpec{Name: "bar", Version: "v1"},
		Status:     v1alpha1.PackageInfoStatus{Manifest: &v1alpha1.PackageManifest{Name: "bar"}}}
	barv2pkg = v1alpha1.ClusterPackage{
		ObjectMeta: v1.ObjectMeta{Name: "bar"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "bar", Version: "v2"}},
		Status:     v1alpha1.PackageStatus{OwnedPackageInfos: []v1alpha1.OwnedResourceRef{{Name: "bar--v2"}}}}
	barv2pi = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: names.PackageInfoName(&barv2pkg)},
		Spec:       v1alpha1.PackageInfoSpec{Name: "bar", Version: "v2"},
		Status:     v1alpha1.PackageInfoStatus{Manifest: &v1alpha1.PackageManifest{Name: "bar"}}}
	bazv1pkg = v1alpha1.ClusterPackage{
		ObjectMeta: v1.ObjectMeta{Name: "baz"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "baz", Version: "v1"}},
		Status:     v1alpha1.PackageStatus{OwnedPackageInfos: []v1alpha1.OwnedResourceRef{{Name: "baz--v1"}}}}
	bazv1pi = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: names.PackageInfoName(&bazv1pkg)},
		Spec:       v1alpha1.PackageInfoSpec{Name: "baz", Version: "v1"},
		Status: v1alpha1.PackageInfoStatus{
			Manifest: &v1alpha1.PackageManifest{Name: "baz", Dependencies: []v1alpha1.Dependency{{Name: "foo", Version: "v1"}}}}}
	bazv2pkg = v1alpha1.ClusterPackage{
		ObjectMeta: v1.ObjectMeta{Name: "baz"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "baz", Version: "v2"}},
		Status:     v1alpha1.PackageStatus{OwnedPackageInfos: []v1alpha1.OwnedResourceRef{{Name: "baz--v2"}}}}
	bazv2pi = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: names.PackageInfoName(&bazv2pkg)},
		Spec:       v1alpha1.PackageInfoSpec{Name: "baz", Version: "v1"},
		Status: v1alpha1.PackageInfoStatus{
			Manifest: &v1alpha1.PackageManifest{Name: "baz", Dependencies: []v1alpha1.Dependency{{Name: "foo", Version: "v2"}}}}}
	nspv1pkg = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "nsp", Namespace: "default"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "nsp", Version: "v1"}}}
	nspv1pi = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: names.PackageInfoName(&nspv1pkg)},
		Spec:       v1alpha1.PackageInfoSpec{Name: "nsp", Version: "v1"},
		Status: v1alpha1.PackageInfoStatus{
			Manifest: &v1alpha1.PackageManifest{Name: "nsp", Dependencies: []v1alpha1.Dependency{{Name: "foo", Version: "v1"}}}}}
	notExistsPkg = v1alpha1.ClusterPackage{
		ObjectMeta: v1.ObjectMeta{Name: "doesnotexist"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "doesnotexist", Version: "v1"}}}
)

// These tests demonstrate how unit tests for a CustomValidator COULD be implemented using the controller-runtime
// fake client.
// They will likely become useless once the PackageValidatingWebhook does some actual validation.
// TODO: Add some tests for package validation

var _ = Describe("PackageValidatingWebhook", Ordered, func() {
	BeforeAll(func() {
		err := v1alpha1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())
		fakeRepoClient.Packages = map[string]map[string]*v1alpha1.PackageManifest{
			"foo": {
				"v1": foov1pi.Status.Manifest,
				"v2": foov2pi.Status.Manifest,
			},
			"bar": {
				"v1": barv1pi.Status.Manifest,
				"v2": barv2pi.Status.Manifest,
			},
			"baz": {
				"v1": bazv1pi.Status.Manifest,
				"v2": bazv2pi.Status.Manifest,
			},
			"nsp": {
				"v1": nspv1pi.Status.Manifest,
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
					webhook := newPackageValidatingWebhook(&barv1pkg, &barv1pi)
					_, err := webhook.ValidateCreate(ctx,
						&foov1pkg)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			When("there is a version mismatch", func() {
				It("should return conflict error", func(ctx context.Context) {
					webhook := newPackageValidatingWebhook(&barv1pkg, &barv1pi)
					_, err := webhook.ValidateCreate(ctx, &foov2pkg)
					Expect(err).To(And(HaveOccurred(), Satisfy(isErrDependencyConflict)))
				})
			})
		})
		When("package has transitive dependency", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateCreate(ctx, &bazv1pkg)
				Expect(err).NotTo(HaveOccurred())
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
					webhook := newPackageValidatingWebhook(&barv2pkg, &barv2pi)
					_, err := webhook.ValidateUpdate(ctx, &foov1pkg, &foov2pkg)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			When("there is a version mismatch", func() {
				It("should return conflict error", func(ctx context.Context) {
					webhook := newPackageValidatingWebhook(&barv1pkg, &barv1pi)
					_, err := webhook.ValidateUpdate(ctx, &foov1pkg, &foov2pkg)
					Expect(err).To(And(HaveOccurred(), Satisfy(isErrDependencyConflict)))
				})
			})
		})
		When("package has transitive dependency", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateUpdate(ctx, &bazv1pkg, &bazv2pkg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("called with object other than Package", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook()
				_, err := webhook.ValidateCreate(ctx, &unstructured.Unstructured{})
				Expect(err).To(HaveOccurred())
			})
		})
		When("a namespaced dependent exists with constraint", func() {
			It("should return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook(&nspv1pkg, &nspv1pi, &foov1pkg, &foov1pi)
				_, err := webhook.ValidateUpdate(ctx, &foov1pkg, &foov2pkg)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("ValidateDelete", func() {
		When("package has no owner references", func() {
			It("should not return error", func(ctx context.Context) {
				webhook := newPackageValidatingWebhook(&foov1pkg, &foov1pi)
				_, err := webhook.ValidateDelete(ctx, &foov1pkg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("package has owner reference", func() {
			It("should return an error", func(ctx context.Context) {
				barv1Copy := barv1pkg.DeepCopy()
				err := controllerutil.SetOwnerReference(&foov1pkg, barv1Copy, scheme.Scheme)
				Expect(err).NotTo(HaveOccurred())

				webhook := newPackageValidatingWebhook(&foov1pkg, &foov1pi, barv1Copy, &barv1pi)
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

					webhook := newPackageValidatingWebhook(foov1Copy, &foov1pi, barv1Copy, &barv1pi)
					_, err = webhook.ValidateDelete(ctx, barv1Copy)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		// TODO tests without owner reference but with owning ref from the other side
	})
})

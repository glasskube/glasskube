package dependency

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/names"
	"github.com/glasskube/glasskube/internal/repo/client/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var fakeRepo = &fake.FakeClient{
	PackageRepositories: []v1alpha1.PackageRepository{
		{},
	},
}

type testClientAdapter struct {
	packages     []v1alpha1.ClusterPackage
	packageInfos []v1alpha1.PackageInfo
}

func (a *testClientAdapter) GetPackageInfo(ctx context.Context, pkgInfoName string) (*v1alpha1.PackageInfo, error) {
	for _, pi := range a.packageInfos {
		if pkgInfoName == pi.Name {
			return &pi, nil
		}
	}
	return nil, errors.NewNotFound(schema.GroupResource{}, pkgInfoName)
}

func (a *testClientAdapter) ListPackages(ctx context.Context) (*v1alpha1.ClusterPackageList, error) {
	return &v1alpha1.ClusterPackageList{Items: a.packages}, nil
}

// GetPackage implements adapter.PackageClientAdapter.
func (a *testClientAdapter) GetPackage(ctx context.Context, name string) (*v1alpha1.ClusterPackage, error) {
	panic("unimplemented")
}

// GetPackageRepository implements adapter.PackageClientAdapter.
func (a *testClientAdapter) GetPackageRepository(ctx context.Context, name string) (
	*v1alpha1.PackageRepository, error) {
	panic("unimplemented")
}

// ListPackageRepositories implements adapter.PackageClientAdapter.
func (a *testClientAdapter) ListPackageRepositories(ctx context.Context) (*v1alpha1.PackageRepositoryList, error) {
	panic("unimplemented")
}

func createDependencyManager() *DependendcyManager {
	testClient = &testClientAdapter{}
	return NewDependencyManager(testClient, &fake.FakeClientset{Client: fakeRepo})
}

var testClient *testClientAdapter
var dm *DependendcyManager

// For the following test suite, we always use the Package p (name "P") as the package who's dependencies should be
// checked.
// Package d (name "D") is the dependency (such that P depends on D) or does not exist, and Packages x (name "X") and
// y (name "Y") are additional optional packages having a dependency on D. For tests where P has multiple dependencies,
// we additionally use Package e (name "E").
var d, e, p, x, y *v1alpha1.ClusterPackage

// di, ei, pi, xi, yi are the corresponding PackageInfo's to the package d, p, x, y
var di, ei, pi, xi, yi *v1alpha1.PackageInfo

func createPackageAndInfo(name string, version string, installed bool) (*v1alpha1.ClusterPackage, *v1alpha1.PackageInfo) {
	manifest := v1alpha1.PackageManifest{
		Name: name,
	}
	fakeRepo.AddPackage(name, version, &manifest)
	pkg := v1alpha1.ClusterPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.PackageSpec{
			PackageInfo: v1alpha1.PackageInfoTemplate{
				Name:    name,
				Version: version,
			},
		},
		Status: v1alpha1.PackageStatus{OwnedPackageInfos: []v1alpha1.OwnedResourceRef{{Name: name}}},
	}

	pkgi := v1alpha1.PackageInfo{
		ObjectMeta: metav1.ObjectMeta{
			Name: names.PackageInfoName(&pkg),
		},
		Spec: v1alpha1.PackageInfoSpec{
			Name:    name,
			Version: version,
		},
		Status: v1alpha1.PackageInfoStatus{
			Version:  version,
			Manifest: &manifest,
		},
	}
	if installed {
		testClient.packages = append(testClient.packages, pkg)
		testClient.packageInfos = append(testClient.packageInfos, pkgi)
	}
	return &pkg, &pkgi
}

var _ = Describe("Dependency Manager", func() {

	BeforeEach(func() {
		dm = createDependencyManager()
		p, pi = createPackageAndInfo("P", "12.2.0", false)
	})

	AfterEach(func() {
		fakeRepo.Clear()
		p = nil
		pi = nil
		d = nil
		di = nil
		e = nil
		ei = nil
		x = nil
		xi = nil
		y = nil
		yi = nil
		dm = nil
	})

	Describe("Validation", func() {

		When("P has no dependencies", func() {
			It("should return OK", func(ctx context.Context) {
				res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).ShouldNot(BeNil())
				Expect(res.Status).Should(Equal(ValidationResultStatusOk))
				Expect(res.Requirements).Should(BeEmpty())
				Expect(res.Conflicts).Should(BeEmpty())
			})
		})

		When("P has a dependency on D", func() {

			When("P requires no version range of D", func() {

				BeforeEach(func() {
					pi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
						Name: "D",
					}}
				})

				When("D exists", func() {

					BeforeEach(func() {
						d, di = createPackageAndInfo("D", "1.1.1", true)
					})

					When("no other package dependent on D", func() {
						It("should return OK", func(ctx context.Context) {
							res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(res).ShouldNot(BeNil())
							Expect(res.Status).Should(Equal(ValidationResultStatusOk))
							Expect(res.Requirements).Should(BeEmpty())
							Expect(res.Conflicts).Should(BeEmpty())
						})
					})

					When("other existing packages X, Y dependent on D", func() {

						BeforeEach(func() {
							x, xi = createPackageAndInfo("X", "0.17.2", true)
							xi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
								Name: "D",
							}}
							y, yi = createPackageAndInfo("Y", "3.2.0-beta.7", true)
							yi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
								Name: "D",
							}}
						})

						When("X and Y require no version range of D", func() {
							It("should return OK", func(ctx context.Context) {
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusOk))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(BeEmpty())
							})
						})

						When("X requires D in version range", func() {
							It("should return OK", func(ctx context.Context) {
								xi.Status.Manifest.Dependencies[0].Version = ">= 1, < 2"
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusOk))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(BeEmpty())
							})
						})

						When("Y requires D in version range", func() {
							It("should return OK", func(ctx context.Context) {
								yi.Status.Manifest.Dependencies[0].Version = "1.x.x"
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusOk))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(BeEmpty())
							})
						})

						When("X and Y require D in version ranges", func() {
							It("should return OK", func(ctx context.Context) {
								xi.Status.Manifest.Dependencies[0].Version = ">= 1, < 2"
								yi.Status.Manifest.Dependencies[0].Version = "1.x.x"
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusOk))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(BeEmpty())
							})
						})
					})
				})

				When("D does not exist", func() {
					BeforeEach(func() {
						fakeRepo.AddPackage("D", "1.1.7", &v1alpha1.PackageManifest{Name: "D"})
					})
					It("should return a RESOLVABLE result with D in latest", func(ctx context.Context) {
						res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res).ShouldNot(BeNil())
						Expect(res.Status).Should(Equal(ValidationResultStatusResolvable))
						Expect(res.Requirements).Should(HaveLen(1))
						Expect(res.Requirements[0].Name).Should(Equal("D"))
						Expect(res.Requirements[0].Version).Should(Equal("1.1.7"))
						Expect(res.Conflicts).Should(BeEmpty())
					})
				})
			})

			When("P requires D in a version range", func() {
				BeforeEach(func() {
					pi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
						Name:    "D",
						Version: "^1.2.3",
					}}
				})

				When("D exists", func() {

					When("no other existing package dependent on D requires a version range of D", func() {

						When("there is no other existing package dependent on D", func() {

							BeforeEach(func() {
								x, xi = createPackageAndInfo("X", "0.17.0", true) // X here has no dependency on D
							})

							It("should return OK if D's version is in required range", func(ctx context.Context) {
								d, di = createPackageAndInfo("D", "1.3.0", true)
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusOk))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(BeEmpty())
							})

							It("should return CONFLICT if D's version is too old", func(ctx context.Context) {
								d, di = createPackageAndInfo("D", "1.2.1", true)
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusConflict))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(HaveLen(1))
								Expect(res.Conflicts[0].Actual.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Actual.Version).Should(Equal("1.2.1"))
								Expect(res.Conflicts[0].Required).ShouldNot(BeNil())
								Expect(res.Conflicts[0].Required.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Required.Version).Should(Equal("^1.2.3"))
							})

							It("should return CONFLICT if D's version is too new", func(ctx context.Context) {
								d, di = createPackageAndInfo("D", "2.0.0-alpha.2", true)
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusConflict))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(HaveLen(1))
								Expect(res.Conflicts[0].Actual.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Actual.Version).Should(Equal("2.0.0-alpha.2"))
								Expect(res.Conflicts[0].Required).ShouldNot(BeNil())
								Expect(res.Conflicts[0].Required.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Required.Version).Should(Equal("^1.2.3"))
							})
						})

						When("existing packages X and Y are dependent on D but require no version range of D", func() {

							BeforeEach(func() {
								x, xi = createPackageAndInfo("X", "0.17.3", true)
								xi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
									Name: "D",
								}}
								y, yi = createPackageAndInfo("Y", "3.2.0-beta.8", true)
								yi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
									Name: "D",
								}}
							})

							// these are the same tests as in the previous When("there is no other existing package dependent on D")

							It("should return OK if D's version is in required range", func(ctx context.Context) {
								d, di = createPackageAndInfo("D", "1.3.1", true)
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusOk))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(BeEmpty())
							})

							It("should return CONFLICT if D's version is too old", func(ctx context.Context) {
								d, di = createPackageAndInfo("D", "1.1.7", true)
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusConflict))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(HaveLen(1))
								Expect(res.Conflicts[0].Actual.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Actual.Version).Should(Equal("1.1.7"))
								Expect(res.Conflicts[0].Required).ShouldNot(BeNil())
								Expect(res.Conflicts[0].Required.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Required.Version).Should(Equal("^1.2.3"))
							})

							It("should return CONFLICT if D's version is too new", func(ctx context.Context) {
								d, di = createPackageAndInfo("D", "2.0.0-alpha.2", true)
								res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(res).ShouldNot(BeNil())
								Expect(res.Status).Should(Equal(ValidationResultStatusConflict))
								Expect(res.Requirements).Should(BeEmpty())
								Expect(res.Conflicts).Should(HaveLen(1))
								Expect(res.Conflicts[0].Actual.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Actual.Version).Should(Equal("2.0.0-alpha.2"))
								Expect(res.Conflicts[0].Required).ShouldNot(BeNil())
								Expect(res.Conflicts[0].Required.Name).Should(Equal("D"))
								Expect(res.Conflicts[0].Required.Version).Should(Equal("^1.2.3"))
							})
						})
					})

					When("other existing packages X, Y are dependent on D and require D in version ranges", func() {

						BeforeEach(func() {
							x, xi = createPackageAndInfo("X", "0.18.3", true)
							xi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
								Name:    "D",
								Version: "^1.0.0 || 2.0.0",
							}}
							y, yi = createPackageAndInfo("X", "3.3.3", true)
							yi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
								Name:    "D",
								Version: ">= 1.1.0, < 3",
							}}
						})

						It("should return OK if D's version is in required range", func(ctx context.Context) {
							d, di = createPackageAndInfo("D", "1.4.0", true)
							res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(res).ShouldNot(BeNil())
							Expect(res.Status).Should(Equal(ValidationResultStatusOk))
							Expect(res.Requirements).Should(BeEmpty())
							Expect(res.Conflicts).Should(BeEmpty())
						})

						It("should return CONFLICT if D's version is too old", func(ctx context.Context) {
							d, di = createPackageAndInfo("D", "1.2.1", true)
							res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(res).ShouldNot(BeNil())
							Expect(res.Status).Should(Equal(ValidationResultStatusConflict))
							Expect(res.Requirements).Should(BeEmpty())
							Expect(res.Conflicts).Should(HaveLen(1))
							Expect(res.Conflicts[0].Actual.Name).Should(Equal("D"))
							Expect(res.Conflicts[0].Actual.Version).Should(Equal("1.2.1"))
							Expect(res.Conflicts[0].Required).ShouldNot(BeNil())
							Expect(res.Conflicts[0].Required.Name).Should(Equal("D"))
							Expect(res.Conflicts[0].Required.Version).Should(Equal("^1.2.3"))
						})

						It("should return CONFLICT if D's version is too new", func(ctx context.Context) {
							d, di = createPackageAndInfo("D", "2.0.0", true)
							res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(res).ShouldNot(BeNil())
							Expect(res.Status).Should(Equal(ValidationResultStatusConflict))
							Expect(res.Requirements).Should(BeEmpty())
							Expect(res.Conflicts).Should(HaveLen(1))
							Expect(res.Conflicts[0].Actual.Name).Should(Equal("D"))
							Expect(res.Conflicts[0].Actual.Version).Should(Equal("2.0.0"))
							Expect(res.Conflicts[0].Required).ShouldNot(BeNil())
							Expect(res.Conflicts[0].Required.Name).Should(Equal("D"))
							Expect(res.Conflicts[0].Required.Version).Should(Equal("^1.2.3"))
						})
					})
				})

				When("D does not exist", func() {
					BeforeEach(func() {
						fakeRepo.AddPackage("D", "1.7.4", &v1alpha1.PackageManifest{Name: "D"})
					})

					It("should return a RESOLVABLE result with latest D", func(ctx context.Context) {
						res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res).ShouldNot(BeNil())
						Expect(res.Status).Should(Equal(ValidationResultStatusResolvable))
						Expect(res.Requirements).Should(HaveLen(1))
						Expect(res.Requirements[0].Name).Should(Equal("D"))
						Expect(res.Requirements[0].Version).Should(Equal("1.7.4"))
						Expect(res.Conflicts).Should(BeEmpty())
					})
				})
			})
		})

		When("P has dependencies on D and E", func() {

			// TODO some more tests here

			BeforeEach(func() {
				pi.Status.Manifest.Dependencies = []v1alpha1.Dependency{{
					Name: "D",
				}, {
					Name: "E",
				}}
			})

			When("P requires no version ranges of D and E", func() {
				When("D, E exist", func() {
					It("Should return OK", func(ctx context.Context) {
						d, di = createPackageAndInfo("D", "118.0.0", true)
						e, ei = createPackageAndInfo("E", "11.80.0", true)
						res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res).ShouldNot(BeNil())
						Expect(res.Status).Should(Equal(ValidationResultStatusOk))
						Expect(res.Requirements).Should(BeEmpty())
						Expect(res.Conflicts).Should(BeEmpty())
					})
				})

				When("D, E do not exist", func() {
					BeforeEach(func() {
						fakeRepo.AddPackage("D", "1.1.7", &v1alpha1.PackageManifest{Name: "D"})
						fakeRepo.AddPackage("E", "1.1.7", &v1alpha1.PackageManifest{Name: "E"})
					})

					It("Should return RESOLVABLE result with D, E as requirements", func(ctx context.Context) {
						res, err := dm.Validate(ctx, pi.Status.Manifest, p.Spec.PackageInfo.Version)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res).ShouldNot(BeNil())
						Expect(res.Status).Should(Equal(ValidationResultStatusResolvable))
						Expect(res.Requirements).Should(HaveLen(2))
						Expect(res.Requirements[0].Name).Should(Equal("D"))
						Expect(res.Requirements[0].Version).Should(Equal("1.1.7"))
						Expect(res.Requirements[1].Name).Should(Equal("E"))
						Expect(res.Requirements[1].Version).Should(Equal("1.1.7"))
						Expect(res.Conflicts).Should(BeEmpty())
					})
				})

			})

			When("P requires D and E to be in version ranges", func() {

				// TODO

			})

		})

	})

	// TODO test IsUpdateAllowed()

	Describe("CheckConflict", func() {

	})

})

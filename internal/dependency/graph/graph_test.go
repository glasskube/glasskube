package graph

import (
	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DependencyGraph", func() {
	var graph *DependencyGraph
	var foo = "foo"
	var bar = "bar"
	var baz = "baz"
	var defaultNs = "default"

	BeforeEach(func() {
		graph = NewGraph()
	})

	Describe("Validate", func() {
		When("graph is empty", func() {
			It("should not return an error", func() {
				Expect(graph.Validate()).NotTo(HaveOccurred())
			})
		})

		When("there is a package without dependencies", func() {
			It("should not return an error", func() {
				Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: foo}, "v1.0.0", true)).NotTo(HaveOccurred())
				Expect(graph.Validate()).NotTo(HaveOccurred())
			})
		})

		When("there is a package with dependencies", func() {
			When("the dependency exists", func() {
				When("there is a constraint", func() {
					It("should not return an error", func() {
						fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
						Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is not violated", func() {
					It("should not return an error", func() {
						fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}}
						Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is violated", func() {
					It("should return an error", func() {
						fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.1.x"}}}
						Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).To(MatchError((error)(&DependencyError{})))
					})
				})
			})

			When("the dependency does not exist", func() {
				It("should return an error", func() {
					fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
					Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
					Expect(graph.Validate()).To(MatchError((error)(&DependencyError{})))
				})
			})
		})

		When("there is a namespaced package with dependencies", func() {
			When("the dependency exists", func() {
				When("there is a constraint", func() {
					It("should not return an error", func() {
						fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
						Expect(graph.AddNamespaced(foo, defaultNs, fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is not violated", func() {
					It("should not return an error", func() {
						fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}}
						Expect(graph.AddNamespaced(foo, defaultNs, fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is violated", func() {
					It("should return an error", func() {
						fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.1.x"}}}
						Expect(graph.AddNamespaced(foo, defaultNs, fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).To(MatchError((error)(&DependencyError{})))
					})
				})
			})

			When("the dependency does not exist", func() {
				It("should return an error", func() {
					fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
					Expect(graph.AddNamespaced(foo, defaultNs, fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
					Expect(graph.Validate()).To(MatchError((error)(&DependencyError{})))
				})
			})
		})
	})

	Describe("Delete", func() {
		It("should remove all properties", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.Version(foo, "")).NotTo(BeNil())
			Expect(graph.Manual(foo, "")).To(BeTrue())
			Expect(graph.Dependencies(foo, "")).NotTo(BeEmpty())
			Expect(graph.Delete(foo, "")).To(BeTrue())
			Expect(graph.Version(foo, "")).To(BeNil())
			Expect(graph.Manual(foo, "")).To(BeFalse())
			Expect(graph.Dependencies(foo, "")).To(BeEmpty())
		})
	})

	Describe("Prune", func() {
		It("should remove orphaned vertex", func() {
			Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
			Expect(graph.Version(bar, "")).To(Equal(semver.MustParse("v1.0.0")))
			Expect(graph.Prune()).To(ConsistOf(PackageRef{bar, "", bar}))
			Expect(graph.Version(bar, "")).To(BeNil())
		})

		It("should remove orphaned vertex transitively", func() {
			Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar}}},
				"v1.0.0", false)).NotTo(HaveOccurred())
			Expect(graph.Prune()).To(ConsistOf(PackageRef{bar, "", bar},
				PackageRef{foo, "", foo}))
			Expect(graph.Version(bar, "")).To(BeNil())
		})
	})

	Describe("DeleteAndPrune", func() {
		It("should remove dependency", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
			Expect(graph.Version(bar, "")).To(Equal(semver.MustParse("v1.0.0")))
			Expect(graph.DeleteAndPrune(foo, "")).To(ConsistOf(PackageRef{foo, "", foo}, PackageRef{bar, "", bar}))
			Expect(graph.Version(bar, "")).To(BeNil())
		})
	})

	Describe("Max", func() {
		It("should return error for empty slice", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.2"}}}
			bazManifest := v1alpha1.PackageManifest{Name: baz, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.1.0"}}}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bazManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			v, err := graph.Max(bar, "", []*semver.Version{})
			Expect(err).To(HaveOccurred())
			Expect(v).To(BeNil())
		})

		It("should return error for no matching version", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.1"}}}
			bazManifest := v1alpha1.PackageManifest{Name: baz, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.1.0"}}}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bazManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			versions := []*semver.Version{semver.MustParse("1.0.0"), semver.MustParse("1.1.1"),
				semver.MustParse("1.2.0"), semver.MustParse("2.0.0")}
			v, err := graph.Max(bar, "", versions)
			Expect(err).To(HaveOccurred())
			Expect(v).To(BeNil())
		})

		It("should return correct version", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.1"}}}
			bazManifest := v1alpha1.PackageManifest{Name: baz, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.1.0"}}}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bazManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			versions := []*semver.Version{semver.MustParse("1.0.0"), semver.MustParse("1.1.0"),
				semver.MustParse("1.1.1"), semver.MustParse("1.2.0"), semver.MustParse("2.0.0")}
			v, err := graph.Max(bar, "", versions)
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(semver.MustParse("1.1.0")))
		})

		It("should consider version metadata for comparison", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.1"}}}
			bazManifest := v1alpha1.PackageManifest{Name: baz, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.1.0"}}}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bazManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			versions := []*semver.Version{semver.MustParse("1.0.0"), semver.MustParse("1.1.0+1"), semver.MustParse("1.1.0+2"),
				semver.MustParse("1.1.1"), semver.MustParse("1.2.0"), semver.MustParse("2.0.0")}
			v, err := graph.Max(bar, "", versions)
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(semver.MustParse("1.1.0+2")))
		})

		It("should not consider version metadata for constraints", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.1"}}}
			bazManifest := v1alpha1.PackageManifest{Name: baz, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "<= 1.1.0+1"}}}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bazManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			versions := []*semver.Version{semver.MustParse("1.0.0"), semver.MustParse("1.1.0+1"), semver.MustParse("1.1.0+2"),
				semver.MustParse("1.1.1"), semver.MustParse("1.2.0"), semver.MustParse("2.0.0")}
			v, err := graph.Max(bar, "", versions)
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(semver.MustParse("1.1.0+2"))) // of baz's <= 1.1.0+1 constraint, the metadata +1 is ignored!
		})
	})

	Describe("Dependencies", func() {
		It("should return all dependencies", func() {
			manifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}, {Name: baz}}}
			Expect(graph.AddCluster(manifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.Dependencies(foo, "")).To(ConsistOf(PackageRef{bar, "", bar}, PackageRef{baz, "", baz}))
		})
	})

	Describe("Dependants", func() {
		It("should return all dependants", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}}
			bazManifest := v1alpha1.PackageManifest{Name: baz, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
			Expect(graph.AddNamespaced(foo, defaultNs, fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bazManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
			Expect(graph.Dependants(bar, "")).To(ContainElements(PackageRef{foo, defaultNs, foo}, PackageRef{baz, "", baz}))
			Expect(graph.Dependants(bar, "")).To(HaveLen(3))
		})
	})

	Describe("Constraints", func() {
		It("should return constraints of dependants", func() {
			fooManifest1 := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.2.x"}}}
			fooManifest2 := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}}
			bazManifest := v1alpha1.PackageManifest{Name: baz, Dependencies: []v1alpha1.Dependency{{Name: bar}}}
			Expect(graph.AddNamespaced(foo, defaultNs, fooManifest1, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(fooManifest2, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bazManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(v1alpha1.PackageManifest{Name: bar}, "v1.0.0", false)).NotTo(HaveOccurred())
			Expect(graph.Constraints(bar, "")).To(ConsistOf(constraint("1.x.x"), constraint("1.2.x")))
		})
	})

	Describe("DeepCopy", func() {
		It("should produce equal graph", func() {
			fooManifest := v1alpha1.PackageManifest{Name: foo, Dependencies: []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}}
			barManifest := v1alpha1.PackageManifest{Name: bar}
			Expect(graph.AddCluster(fooManifest, "v1.0.0", true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(barManifest, "v1.0.0", false)).NotTo(HaveOccurred())
			newGraph := graph.DeepCopy()
			Expect(newGraph).To(Equal(graph))
			Expect(graph.AddCluster(barManifest, "v1.1.0", false)).NotTo(HaveOccurred())
			Expect(newGraph).NotTo(Equal(graph))
		})
	})

	Describe("Version", func() {
		It("should return nil for missing package", func() {
			Expect(graph.Version("foo", "")).To(BeNil())
		})
	})

	Describe("Manual", func() {
		It("should return false for missing package", func() {
			Expect(graph.Manual("foo", "")).To(BeFalse())
		})
	})
})

func constraint(pattern string) *semver.Constraints {
	c, err := semver.NewConstraint(pattern)
	Expect(err).NotTo(HaveOccurred())
	return c
}

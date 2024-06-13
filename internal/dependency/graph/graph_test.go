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
				Expect(graph.AddCluster(foo, "v1.0.0", nil, true)).NotTo(HaveOccurred())
				Expect(graph.Validate()).NotTo(HaveOccurred())
			})
		})

		When("there is a package with dependencies", func() {
			When("the dependency exists", func() {
				When("there is a constraint", func() {
					It("should not return an error", func() {
						Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar}}, true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is not violated", func() {
					It("should not return an error", func() {
						Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}, true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is violated", func() {
					It("should return an error", func() {
						Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.1.x"}}, true)).NotTo(HaveOccurred())
						Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).To(MatchError(&DependencyError{}))
					})
				})
			})

			When("the dependency does not exist", func() {
				It("should return an error", func() {
					Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar}}, true)).NotTo(HaveOccurred())
					Expect(graph.Validate()).To(MatchError(&DependencyError{}))
				})
			})
		})

		When("there is a namespaced package with dependencies", func() {
			When("the dependency exists", func() {
				When("there is a constraint", func() {
					It("should not return an error", func() {
						Expect(graph.AddNamespaced(foo, foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar}})).NotTo(HaveOccurred())
						Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is not violated", func() {
					It("should not return an error", func() {
						Expect(graph.AddNamespaced(foo, foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}})).NotTo(HaveOccurred())
						Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).NotTo(HaveOccurred())
					})
				})

				When("there is a constraint and it is violated", func() {
					It("should return an error", func() {
						Expect(graph.AddNamespaced(foo, foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.1.x"}})).NotTo(HaveOccurred())
						Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
						Expect(graph.Validate()).To(MatchError(&DependencyError{}))
					})
				})
			})

			When("the dependency does not exist", func() {
				It("should return an error", func() {
					Expect(graph.AddNamespaced(foo, foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar}})).NotTo(HaveOccurred())
					Expect(graph.Validate()).To(MatchError(&DependencyError{}))
				})
			})
		})
	})

	Describe("Delete", func() {
		It("should remove all properties", func() {
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar}}, true)).NotTo(HaveOccurred())
			Expect(graph.Version(foo)).NotTo(BeNil())
			Expect(graph.Manual(foo)).To(BeTrue())
			Expect(graph.Dependencies(foo)).NotTo(BeEmpty())
			Expect(graph.Delete(foo)).To(BeTrue())
			Expect(graph.Version(foo)).To(BeNil())
			Expect(graph.Manual(foo)).To(BeFalse())
			Expect(graph.Dependencies(foo)).To(BeEmpty())
		})
	})

	Describe("Prune", func() {
		It("should remove orphaned vertex", func() {
			Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
			Expect(graph.Version(bar)).To(Equal(semver.MustParse("v1.0.0")))
			Expect(graph.Prune()).To(ConsistOf(bar))
			Expect(graph.Version(bar)).To(BeNil())
		})
	})

	Describe("DeleteAndPrune", func() {
		It("should remove dependency", func() {
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar}}, true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
			Expect(graph.Version(bar)).To(Equal(semver.MustParse("v1.0.0")))
			Expect(graph.DeleteAndPrune(foo)).To(ConsistOf(foo, bar))
			Expect(graph.Version(bar)).To(BeNil())
		})
	})

	Describe("Max", func() {
		It("should return error for empty slice", func() {
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.2"}}, true)).
				NotTo(HaveOccurred())
			Expect(graph.AddCluster(baz, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: ">= 1.1.0"}}, true)).NotTo(HaveOccurred())
			v, err := graph.Max(bar, []*semver.Version{})
			Expect(err).To(HaveOccurred())
			Expect(v).To(BeNil())
		})

		It("should return error for no matching version", func() {
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.1"}}, true)).
				NotTo(HaveOccurred())
			Expect(graph.AddCluster(baz, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: ">= 1.1.0"}}, true)).NotTo(HaveOccurred())
			versions := []*semver.Version{semver.MustParse("1.0.0"), semver.MustParse("1.1.1"), semver.MustParse("1.2.0"),
				semver.MustParse("2.0.0")}
			v, err := graph.Max(bar, versions)
			Expect(err).To(HaveOccurred())
			Expect(v).To(BeNil())
		})

		It("should return correct version", func() {
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: ">= 1.0.0, < 1.1.1"}}, true)).
				NotTo(HaveOccurred())
			Expect(graph.AddCluster(baz, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: ">= 1.1.0"}}, true)).
				NotTo(HaveOccurred())
			versions := []*semver.Version{semver.MustParse("1.0.0"), semver.MustParse("1.1.0"), semver.MustParse("1.1.1"),
				semver.MustParse("1.2.0"), semver.MustParse("2.0.0")}
			v, err := graph.Max(bar, versions)
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(semver.MustParse("1.1.0")))
		})
	})

	Describe("Dependencies", func() {
		It("should return all dependencies", func() {
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}, {Name: baz}}, true)).
				NotTo(HaveOccurred())
			Expect(graph.Dependencies(foo)).To(ConsistOf(bar, baz))
		})
	})

	Describe("Dependants", func() {
		It("should return all dependants", func() {
			Expect(graph.AddNamespaced(foo, foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}})).NotTo(HaveOccurred())
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}, true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(baz, "v1.0.0", []v1alpha1.Dependency{{Name: bar}}, true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
			Expect(graph.Dependants(bar)).To(ContainElements(foo, baz))
			Expect(graph.Dependants(bar)).To(HaveLen(3))
		})
	})

	Describe("Constraints", func() {
		It("should return constraints of dependants", func() {
			Expect(graph.AddNamespaced(foo, foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.2.x"}})).NotTo(HaveOccurred())
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}, true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(baz, "v1.0.0", []v1alpha1.Dependency{{Name: bar}}, true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
			Expect(graph.Constraints(bar)).To(ConsistOf(constraint("1.x.x"), constraint("1.2.x")))
		})
	})

	Describe("DeepCopy", func() {
		It("should produce equal graph", func() {
			Expect(graph.AddCluster(foo, "v1.0.0", []v1alpha1.Dependency{{Name: bar, Version: "1.x.x"}}, true)).NotTo(HaveOccurred())
			Expect(graph.AddCluster(bar, "v1.0.0", nil, false)).NotTo(HaveOccurred())
			newGraph := graph.DeepCopy()
			Expect(graph.AddCluster(bar, "v1.1.0", nil, false)).NotTo(HaveOccurred())
			Expect(newGraph).NotTo(Equal(graph))
		})
	})

	Describe("Version", func() {
		It("should return nil for missing package", func() {
			Expect(graph.Version("foo")).To(BeNil())
		})
	})

	Describe("Manual", func() {
		It("should return false for missing package", func() {
			Expect(graph.Manual("foo")).To(BeFalse())
		})
	})
})

func constraint(pattern string) *semver.Constraints {
	c, err := semver.NewConstraint(pattern)
	Expect(err).NotTo(HaveOccurred())
	return c
}

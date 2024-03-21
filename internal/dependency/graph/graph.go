package graph

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	isemver "github.com/glasskube/glasskube/internal/semver"
	"go.uber.org/multierr"
)

type DependencyGraph struct {
	vertices map[string]*vertex
}

type vertex struct {
	version *semver.Version
	manual  bool
	edges   map[string]*edge
}

type edge struct {
	constraint *semver.Constraints
	vertex     *vertex
}

func NewGraph() *DependencyGraph {
	return &DependencyGraph{vertices: map[string]*vertex{}}
}

// Add simulates installing or updating a package by
// 1. Creating a vertex if necessary
// 2. Setting its version and
// 3. Updating the outgoing edges of the vertex to match the manifests dependencies declaration
func (g *DependencyGraph) Add(name, version string, dependencies []v1alpha1.Dependency, manual bool) error {
	_, err := g.add(name, version, dependencies, manual)
	return err
}

// Manual returns whether a package has been manually added by a user
func (g *DependencyGraph) Manual(name string) bool {
	if vertex, ok := g.vertices[name]; ok {
		return vertex.manual
	} else {
		return false
	}
}

// Version returns the installed version of a package or nil if that package is not installed
func (g *DependencyGraph) Version(of string) *semver.Version {
	if vertex, ok := g.vertices[of]; ok {
		return vertex.version
	} else {
		return nil
	}
}

// Dependencies returns the names of packages that this package depends on
func (g *DependencyGraph) Dependencies(of string) []string {
	if vertex, ok := g.vertices[of]; ok {
		dependencies := make([]string, len(vertex.edges))
		i := 0
		for dep := range vertex.edges {
			dependencies[i] = dep
			i++
		}
		return dependencies
	} else {
		return nil
	}
}

// Dependants returns the names of packages that depend on this package
func (g *DependencyGraph) Dependants(of string) []string {
	var dependants []string
	for name, vertex := range g.vertices {
		if _, ok := vertex.edges[of]; vertex.version != nil && ok {
			dependants = append(dependants, name)
		}
	}
	return dependants
}

// Constraints returns all constraints of dependants of this package
func (g *DependencyGraph) Constraints(of string) []*semver.Constraints {
	var constraints []*semver.Constraints
	for _, vertex := range g.vertices {
		if edge, ok := vertex.edges[of]; ok && vertex.version != nil && edge.constraint != nil {
			constraints = append(constraints, edge.constraint)
		}
	}
	return constraints
}

// Max returns the maximum element of versions that does not violate any constraint of this package
func (g *DependencyGraph) Max(of string, versions []*semver.Version) (*semver.Version, error) {
	var maxVersion *semver.Version
outer:
	for _, version := range versions {
		if maxVersion == nil || maxVersion.LessThan(version) {
			for _, constraint := range g.Constraints(of) {
				if isemver.ValidateVersionConstraint(version, constraint) != nil {
					continue outer
				}
			}
			maxVersion = version
		}
	}
	if maxVersion != nil {
		return maxVersion, nil
	} else {
		return nil, fmt.Errorf("no matching version for %v found", of)
	}
}

// DeepCopy returns an exact copy of this graph
func (oldGraph *DependencyGraph) DeepCopy() *DependencyGraph {
	newGraph := NewGraph()
	for vertexName, vertex := range oldGraph.vertices {
		newVertex := newGraph.vertex(vertexName)
		newVertex.version = vertex.version
		newVertex.manual = vertex.manual
		for edgeName, edge := range vertex.edges {
			newGraph.edge(vertexName, edgeName, edge.constraint)
		}
	}
	return newGraph
}

// Delete simulates uninstalling a package.
//
// The vertex is not actually removed from the graph, as it may still be referenced by other
// packages and needs to be kept for validation! Instead, its version is unset and its dependencies
// are cleared.
func (g *DependencyGraph) Delete(name string) bool {
	return g.delete(name)
}

// Prune deletes all vertices for which all of the following applies:
// 1. It has not been installed manually
// 2. It does not have any dependants
func (g *DependencyGraph) Prune() []string {
	stable := false
	var removed []string
	for !stable {
		stable = true
		for name, vertex := range g.vertices {
			if !vertex.manual && len(g.Dependants(name)) == 0 && g.delete(name) {
				stable = false
				removed = append(removed, name)
			}
		}
	}
	return removed
}

func (g *DependencyGraph) DeleteAndPrune(name string) []string {
	if g.delete(name) {
		return append([]string{name}, g.Prune()...)
	}
	return nil
}

func (g *DependencyGraph) ValidateDelete(name string) ([]string, error) {
	gc := g.DeepCopy()
	return gc.DeleteAndPrune(name), gc.Validate()
}

// Validate checks the consistency of the entire graph by checking that
// 1. All vertices with at least one dependency have a version that is not nil
// 2. There are no violated version constraints
func (g *DependencyGraph) Validate() error {
	var err error
	for name, vertex := range g.vertices {
		for dep, edge := range vertex.edges {
			if edge.vertex.version == nil {
				multierr.AppendInto(&err, ErrDependency(name, dep, ErrNotInstalled(dep)))
			} else if edge.constraint != nil {
				if err1 := isemver.ValidateVersionConstraint(edge.vertex.version, edge.constraint); err1 != nil {
					multierr.AppendInto(&err, ErrDependency(name, dep, ErrConstraint(dep, edge.vertex.version, edge.constraint, err1)))
				}
			}
		}
	}
	return err
}

func (g *DependencyGraph) vertex(name string) *vertex {
	if n, ok := g.vertices[name]; ok {
		return n
	} else {
		return g.createVertex(name, nil)
	}
}

func (g *DependencyGraph) createVertex(name string, version *semver.Version) *vertex {
	n := &vertex{
		version: version,
		edges:   map[string]*edge{},
	}
	g.vertices[name] = n
	return n
}

func (g *DependencyGraph) edge(from string, to string, constraint *semver.Constraints) {
	g.vertex(from).edges[to] = &edge{
		vertex:     g.vertex(to),
		constraint: constraint,
	}
}

func (g *DependencyGraph) add(name, version string, dependencies []v1alpha1.Dependency, manual bool) (*vertex, error) {
	if version == "" {
		g.delete(name)
		return g.vertex(name), nil
	}

	parsedVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	vertex := g.vertex(name)
	vertex.version = parsedVersion
	vertex.manual = manual
	vertex.edges = map[string]*edge{}

	for _, dep := range dependencies {
		var constraint *semver.Constraints
		if len(dep.Version) > 0 {
			if c, err := semver.NewConstraint(dep.Version); err != nil {
				return nil, err
			} else {
				constraint = c
			}
		}
		g.edge(name, dep.Name, constraint)
	}

	return vertex, nil
}

func (g *DependencyGraph) delete(name string) bool {
	vertex := g.vertex(name)
	deleted := vertex.version != nil
	vertex.version = nil
	vertex.manual = false
	vertex.edges = map[string]*edge{}
	return deleted
}

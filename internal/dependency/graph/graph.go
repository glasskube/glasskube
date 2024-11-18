package graph

import (
	"fmt"

	"k8s.io/client-go/tools/cache"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/dependency/util"
	isemver "github.com/glasskube/glasskube/internal/semver"
	"go.uber.org/multierr"
)

type PackageRef struct {
	Name, Namespace, PackageName string
}

func (ref PackageRef) String() string {
	return cache.ObjectName{Namespace: ref.Namespace, Name: ref.Name}.String()
}

type vertexRef struct {
	name, namespace string
}

type vertexMap map[vertexRef]*vertex

func (m vertexMap) vertex(key vertexRef, packageName string) *vertex {
	if v, ok := m[key]; ok {
		return v
	} else {
		return m.createVertex(key, packageName)
	}
}

func (m vertexMap) createVertex(key vertexRef, packageName string) *vertex {
	n := &vertex{packageName: packageName, edges: make(map[vertexRef]*edge)}
	m[key] = n
	return n
}

func (m *vertexMap) edge(from *vertex, to vertexRef, constraint *semver.Constraints) {
	from.edges[to] = &edge{
		vertex:     m.vertex(to, ""),
		constraint: constraint,
	}
}

func (oldMap vertexMap) deepCopy() vertexMap {
	newMap := vertexMap{}
	for ref, vertex := range oldMap {
		newVertex := newMap.vertex(ref, vertex.packageName)
		newVertex.version = vertex.version
		newVertex.manual = vertex.manual
		for edgeRef, edge := range vertex.edges {
			newMap.vertex(edgeRef, edge.vertex.packageName)
			newMap.edge(newVertex, edgeRef, edge.constraint)
		}
	}
	return newMap
}

type DependencyGraph struct {
	vertices vertexMap
}

type vertex struct {
	packageName string
	version     *semver.Version
	manual      bool
	edges       map[vertexRef]*edge
}

type edge struct {
	constraint *semver.Constraints
	vertex     *vertex
}

func NewGraph() *DependencyGraph {
	return &DependencyGraph{vertices: make(vertexMap)}
}

// AddCluster simulates installing or updating a ClusterPackage by
// 1. Creating a vertex if necessary
// 2. Setting its version and
// 3. Updating the outgoing edges of the vertex to match the manifests dependencies declaration
func (g *DependencyGraph) AddCluster(manifest v1alpha1.PackageManifest, version string, manual bool) error {
	return g.add(vertexRef{name: manifest.Name}, manifest, version, manual)
}

func (g *DependencyGraph) AddNamespaced(
	name, namespace string,
	manifest v1alpha1.PackageManifest,
	version string,
	manual bool,
) error {
	return g.add(vertexRef{name: name, namespace: namespace}, manifest, version, manual)
}

// Manual returns whether a package has been manually added by a user
func (g *DependencyGraph) Manual(name string, namespace string) bool {
	if vertex, ok := g.vertices[vertexRef{name, namespace}]; ok {
		return vertex.manual
	} else {
		return false
	}
}

// Version returns the installed version of a package or nil if that package is not installed
func (g *DependencyGraph) Version(of string, namespace string) *semver.Version {
	if vertex, ok := g.vertices[vertexRef{name: of, namespace: namespace}]; ok {
		return vertex.version
	} else {
		return nil
	}
}

// Dependencies returns the names of packages that this package depends on
func (g *DependencyGraph) Dependencies(of, namespace string) []PackageRef {
	if vertex, ok := g.vertices[vertexRef{name: of, namespace: namespace}]; ok {
		dependencies := make([]PackageRef, len(vertex.edges))
		i := 0
		for ref, dep := range vertex.edges {
			dependencies[i] = PackageRef{Name: ref.name, Namespace: ref.namespace, PackageName: dep.vertex.packageName}
			i++
		}
		return dependencies
	} else {
		return nil
	}
}

// Dependants returns the names of packages that depend on this package
func (g *DependencyGraph) Dependants(of, namespace string) []PackageRef {
	var dependants []PackageRef
	for ref, vertex := range g.vertices {
		if _, ok := vertex.edges[vertexRef{name: of, namespace: namespace}]; vertex.version != nil && ok {
			dependants = append(dependants,
				PackageRef{Name: ref.name, Namespace: ref.namespace, PackageName: vertex.packageName},
			)
		}
	}
	return dependants
}

// Constraints returns all constraints of dependants of this package
func (g *DependencyGraph) Constraints(of, namespace string) []*semver.Constraints {
	var constraints []*semver.Constraints
	for _, vertex := range g.vertices {
		if edge, ok := vertex.edges[vertexRef{name: of, namespace: namespace}]; ok && vertex.version != nil && edge.constraint != nil {
			constraints = append(constraints, edge.constraint)
		}
	}
	return constraints
}

// Max returns the maximum element of versions that does not violate any constraint of this package. Note that it
// also interprets the metadata of the versions, just as in IsVersionUpgradable.
func (g *DependencyGraph) Max(of, namespace string, versions []*semver.Version) (*semver.Version, error) {
	var maxVersion *semver.Version
outer:
	for _, version := range versions {
		if maxVersion == nil || isemver.IsVersionUpgradable(maxVersion, version) {
			for _, constraint := range g.Constraints(of, namespace) {
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
	return &DependencyGraph{
		vertices: oldGraph.vertices.deepCopy(),
	}
}

// Delete simulates uninstalling a package.
//
// The vertex is not actually removed from the graph, as it may still be referenced by other
// packages and needs to be kept for validation! Instead, its version is unset and its dependencies
// are cleared.
func (g *DependencyGraph) Delete(name, namespace string) bool {
	return g.delete(g.vertices.vertex(vertexRef{name: name, namespace: namespace}, ""))
}

// Prune deletes all vertices for which all of the following applies:
// 1. It has not been installed manually
// 2. It does not have any dependants
func (g *DependencyGraph) Prune() []PackageRef {
	stable := false
	var removed []PackageRef
	for !stable {
		stable = true
		for ref, vertex := range g.vertices {
			if !vertex.manual && len(g.Dependants(ref.name, ref.namespace)) == 0 && g.delete(vertex) {
				stable = false
				removed = append(removed,
					PackageRef{Name: ref.name, Namespace: ref.namespace, PackageName: vertex.packageName},
				)
			}
		}
	}
	return removed
}

func (g *DependencyGraph) DeleteAndPrune(name, namespace string) []PackageRef {
	vertex := g.vertices.vertex(vertexRef{name: name, namespace: namespace}, "")
	if g.delete(vertex) {
		return append([]PackageRef{{Name: name, Namespace: namespace, PackageName: vertex.packageName}}, g.Prune()...)
	}
	return nil
}

func (g *DependencyGraph) ValidateDelete(name, namespace string) ([]PackageRef, error) {
	gc := g.DeepCopy()
	return gc.DeleteAndPrune(name, namespace), gc.Validate()
}

// Validate checks the consistency of the entire graph by checking that
// 1. All vertices with at least one dependency have a version that is not nil
// 2. There are no violated version constraints
func (g *DependencyGraph) Validate() error {
	var err error
	for ref, vertex := range g.vertices {
		pkgRef := PackageRef{Name: ref.name, Namespace: ref.namespace, PackageName: vertex.packageName}
		for dep, edge := range vertex.edges {
			depRef := PackageRef{Name: dep.name, Namespace: dep.namespace, PackageName: edge.vertex.packageName}
			if edge.vertex.version == nil {
				multierr.AppendInto(&err, ErrDependency(pkgRef, depRef, ErrNotInstalled(depRef)))
			} else if edge.constraint != nil {
				if err1 := isemver.ValidateVersionConstraint(edge.vertex.version, edge.constraint); err1 != nil {
					multierr.AppendInto(&err, ErrDependency(pkgRef, depRef,
						ErrConstraint(depRef, edge.vertex.version, edge.constraint, err1)))
				}
			}
		}
	}
	return err
}

func (g *DependencyGraph) add(
	ref vertexRef,
	manifest v1alpha1.PackageManifest,
	version string,
	manual bool,
) error {
	vertex := g.vertices.vertex(ref, manifest.Name)

	if version == "" {
		g.delete(vertex)
		return nil
	}

	parsedVersion, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	vertex.version = parsedVersion
	vertex.manual = manual
	vertex.edges = map[vertexRef]*edge{}

	for _, dep := range manifest.Dependencies {
		depRef := vertexRef{name: dep.Name}
		g.vertices.vertex(depRef, dep.Name)

		var constraint *semver.Constraints
		if len(dep.Version) > 0 {
			if c, err := semver.NewConstraint(dep.Version); err != nil {
				return err
			} else {
				constraint = c
			}
		}
		g.vertices.edge(vertex, depRef, constraint)
	}

	for _, cmp := range manifest.Components {
		cmpRef := vertexRef{
			name:      util.ComponentName(ref.name, cmp),
			namespace: ref.namespace,
		}
		if cmpRef.namespace == "" {
			cmpRef.namespace = manifest.DefaultNamespace
		}
		g.vertices.vertex(cmpRef, cmp.Name)

		var constraint *semver.Constraints
		if len(cmp.Version) > 0 {
			if c, err := semver.NewConstraint(cmp.Version); err != nil {
				return err
			} else {
				constraint = c
			}
		}
		g.vertices.edge(vertex, cmpRef, constraint)
	}

	return nil
}

func (g *DependencyGraph) delete(vertex *vertex) bool {
	deleted := vertex.version != nil
	vertex.version = nil
	vertex.manual = false
	vertex.edges = make(map[vertexRef]*edge)
	return deleted
}

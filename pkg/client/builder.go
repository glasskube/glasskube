package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type packageBuilder struct {
	manifestName, version, repositoryName string
	namespace, name                       string
	autoUpdate                            bool
	values                                map[string]v1alpha1.ValueConfiguration
}

func PackageBuilder(name string) *packageBuilder {
	return &packageBuilder{
		manifestName: name,
		values:       make(map[string]v1alpha1.ValueConfiguration),
	}
}

func (b *packageBuilder) WithVersion(version string) *packageBuilder {
	b.version = version
	return b
}

func (b *packageBuilder) WithAutoUpdates(enabled bool) *packageBuilder {
	b.autoUpdate = enabled
	return b
}

func (b *packageBuilder) WithRepositoryName(repositoryName string) *packageBuilder {
	b.repositoryName = repositoryName
	return b
}

func (b *packageBuilder) WithNamespace(namespace string) *packageBuilder {
	b.namespace = namespace
	return b
}

func (b *packageBuilder) WithName(name string) *packageBuilder {
	b.name = name
	return b
}

func (b *packageBuilder) WithValues(values map[string]v1alpha1.ValueConfiguration) *packageBuilder {
	for name, value := range values {
		b.values[name] = value
	}
	return b
}

func (b *packageBuilder) BuildClusterPackage() *v1alpha1.ClusterPackage {
	pkg := v1alpha1.ClusterPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name: b.manifestName,
		},
		Spec: v1alpha1.PackageSpec{
			PackageInfo: v1alpha1.PackageInfoTemplate{
				Name:           b.manifestName,
				Version:        b.version,
				RepositoryName: b.repositoryName,
			},
			Values: b.values,
		},
	}
	pkg.SetAutoUpdatesEnabled(b.autoUpdate)
	return &pkg
}

func (b *packageBuilder) BuildPackage() *v1alpha1.Package {
	pkg := v1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.name,
			Namespace: b.namespace,
		},
		Spec: v1alpha1.PackageSpec{
			PackageInfo: v1alpha1.PackageInfoTemplate{
				Name:           b.manifestName,
				Version:        b.version,
				RepositoryName: b.repositoryName,
			},
			Values: b.values,
		},
	}
	pkg.SetAutoUpdatesEnabled(b.autoUpdate)
	return &pkg
}

func (b *packageBuilder) Build(scope *v1alpha1.PackageScope) ctrlpkg.Package {
	if scope.IsCluster() {
		return b.BuildClusterPackage()
	} else {
		return b.BuildPackage()
	}
}

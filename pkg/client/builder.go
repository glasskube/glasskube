package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type packageBuilder struct {
	name, version, repositoryName string
	autoUpdate                    bool
	values                        map[string]v1alpha1.ValueConfiguration
}

func ClusterPackageBuilder(name string) *packageBuilder {
	return &packageBuilder{
		name:   name,
		values: make(map[string]v1alpha1.ValueConfiguration),
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

func (b *packageBuilder) WithValues(values map[string]v1alpha1.ValueConfiguration) *packageBuilder {
	for name, value := range values {
		b.values[name] = value
	}
	return b
}

func (b *packageBuilder) Build() *v1alpha1.ClusterPackage {
	pkg := v1alpha1.ClusterPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name: b.name,
		},
		Spec: v1alpha1.PackageSpec{
			PackageInfo: v1alpha1.PackageInfoTemplate{
				Name:           b.name,
				Version:        b.version,
				RepositoryName: b.repositoryName,
			},
			Values: b.values,
		},
	}
	pkg.SetAutoUpdatesEnabled(b.autoUpdate)
	return &pkg
}

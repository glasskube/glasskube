package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type packageBuilder struct {
	name, version string
	autoUpdate    bool
	values        map[string]v1alpha1.ValueConfiguration
}

func PackageBuilder(name string) *packageBuilder {
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

func (b *packageBuilder) WithValues(values map[string]v1alpha1.ValueConfiguration) *packageBuilder {
	for name, value := range values {
		b.values[name] = value
	}
	return b
}

func (b *packageBuilder) Build() *v1alpha1.Package {
	pkg := v1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: b.name,
		},
		Spec: v1alpha1.PackageSpec{
			PackageInfo: v1alpha1.PackageInfoTemplate{
				Name:    b.name,
				Version: b.version,
			},
			Values: b.values,
		},
	}
	if b.autoUpdate {
		pkg.SetLabels(map[string]string{
			"packages.glasskube.dev/auto-update": "true",
		})
	}
	return &pkg
}

package client

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type PackageV1Alpha1Client interface {
	Packages() PackageInterface
	PackageInfos() PackageInfoInterface
	PackageRepositories() PackageRepositoryInterface
	WithStores(
		packageStore cache.Store,
		packageInfoStore cache.Store,
		packageRepositoryStore cache.Store,
	) PackageV1Alpha1Client
}

type PackageInterface interface {
	readWriteClientInterface[v1alpha1.Package, v1alpha1.PackageList]
}

type PackageInfoInterface interface {
	readOnlyClientInterface[v1alpha1.PackageInfo, v1alpha1.PackageInfoList]
}

type PackageRepositoryInterface interface {
	readWriteClientInterface[v1alpha1.PackageRepository, v1alpha1.PackageRepositoryList]
}

type readOnlyClientInterface[T any, L any] interface {
	Get(ctx context.Context, name string, target *T) error
	GetAll(ctx context.Context, target *L) error
	Watch(ctx context.Context) (watch.Interface, error)
}

type readWriteClientInterface[T any, L any] interface {
	readOnlyClientInterface[T, L]
	Create(ctx context.Context, target *T, opts metav1.CreateOptions) error
	Update(ctx context.Context, target *T) error
	Delete(ctx context.Context, target *T, options metav1.DeleteOptions) error
}

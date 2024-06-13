package client

import (
	"context"
	"errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type cacheClientset struct {
	PackageV1Alpha1Client
	clusterPackageStore    cache.Store
	packageStore           cache.Store
	packageInfoStore       cache.Store
	packageRepositoryStore cache.Store
}

func (c *cacheClientset) ClusterPackages() ClusterPackageInterface {
	p := c.PackageV1Alpha1Client.ClusterPackages()
	if c.clusterPackageStore == nil {
		return p
	}
	return &readWriteCacheClient[v1alpha1.ClusterPackage, v1alpha1.ClusterPackageList]{
		p,
		readOnlyCacheClient[v1alpha1.ClusterPackage, v1alpha1.ClusterPackageList]{
			p,
			c.clusterPackageStore,
			func(items []v1alpha1.ClusterPackage) v1alpha1.ClusterPackageList {
				return v1alpha1.ClusterPackageList{Items: items}
			},
		},
	}
}

func (c *cacheClientset) Packages(ns string) PackageInterface {
	p := c.PackageV1Alpha1Client.Packages(ns)
	if c.packageStore == nil {
		return p
	}
	return &readWriteCacheClient[v1alpha1.Package, v1alpha1.PackageList]{
		p,
		readOnlyCacheClient[v1alpha1.Package, v1alpha1.PackageList]{
			p,
			c.packageStore,
			func(items []v1alpha1.Package) v1alpha1.PackageList {
				return v1alpha1.PackageList{Items: items}
			},
		},
	}
}

func (c *cacheClientset) PackageInfos() PackageInfoInterface {
	pi := c.PackageV1Alpha1Client.PackageInfos()
	if c.packageInfoStore == nil {
		return pi
	}
	return &readOnlyCacheClient[v1alpha1.PackageInfo, v1alpha1.PackageInfoList]{
		pi,
		c.packageInfoStore,
		func(items []v1alpha1.PackageInfo) v1alpha1.PackageInfoList {
			return v1alpha1.PackageInfoList{Items: items}
		},
	}
}

func (c *cacheClientset) PackageRepositories() PackageRepositoryInterface {
	pr := c.PackageV1Alpha1Client.PackageRepositories()
	if c.packageRepositoryStore == nil {
		return pr
	}
	return &readWriteCacheClient[v1alpha1.PackageRepository, v1alpha1.PackageRepositoryList]{
		pr,
		readOnlyCacheClient[v1alpha1.PackageRepository, v1alpha1.PackageRepositoryList]{
			pr,
			c.packageRepositoryStore,
			func(items []v1alpha1.PackageRepository) v1alpha1.PackageRepositoryList {
				return v1alpha1.PackageRepositoryList{Items: items}
			},
		},
	}
}

type readOnlyCacheClient[T any, L any] struct {
	fallback    readOnlyClientInterface[T, L]
	store       cache.Store
	listFactory func(items []T) L
}

func (c *readOnlyCacheClient[T, L]) Watch(ctx context.Context) (watch.Interface, error) {
	return c.fallback.Watch(ctx)
}

func (c *readOnlyCacheClient[T, L]) Get(ctx context.Context, name string, target *T) error {
	if obj, ok, err := c.store.GetByKey(name); err != nil {
		return apierrors.NewInternalError(err)
	} else if !ok {
		return c.fallback.Get(ctx, name, target)
	} else if obj, ok := obj.(*T); !ok {
		return apierrors.NewInternalError(errors.New("resource exists but has wrong type"))
	} else {
		*target = *obj
		return nil
	}
}

func (c *readOnlyCacheClient[T, L]) GetAll(ctx context.Context, target *L) error {
	objs := c.store.List()
	items := make([]T, len(objs))
	for i, obj := range objs {
		if obj, ok := obj.(*T); !ok {
			return apierrors.NewInternalError(errors.New("resource has has wrong type"))
		} else {
			items[i] = *obj
		}
	}
	*target = c.listFactory(items)
	return nil
}

type readWriteCacheClient[T any, L any] struct {
	fallback readWriteClientInterface[T, L]
	readOnlyCacheClient[T, L]
}

func (c *readWriteCacheClient[T, L]) Create(ctx context.Context, target *T, opts metav1.CreateOptions) error {
	return c.fallback.Create(ctx, target, opts)
}

func (c *readWriteCacheClient[T, L]) Update(ctx context.Context, target *T) error {
	return c.fallback.Update(ctx, target)
}

func (c *readWriteCacheClient[T, L]) Delete(ctx context.Context, target *T, options metav1.DeleteOptions) error {
	return c.fallback.Delete(ctx, target, options)
}

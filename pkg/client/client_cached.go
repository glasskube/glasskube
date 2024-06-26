package client

import (
	"context"
	"errors"
	"fmt"
	"slices"

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
		fallback: p,
		readOnlyCacheClient: readOnlyCacheClient[v1alpha1.ClusterPackage, v1alpha1.ClusterPackageList]{
			fallback: p,
			store:    c.clusterPackageStore,
			listFactory: func(items []v1alpha1.ClusterPackage) v1alpha1.ClusterPackageList {
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
		fallback: p,
		readOnlyCacheClient: readOnlyCacheClient[v1alpha1.Package, v1alpha1.PackageList]{
			fallback: p,
			store:    c.packageStore,
			listFactory: func(items []v1alpha1.Package) v1alpha1.PackageList {
				return v1alpha1.PackageList{Items: items}
			},
			namespace: ns,
		},
	}
}

func (c *cacheClientset) PackageInfos() PackageInfoInterface {
	pi := c.PackageV1Alpha1Client.PackageInfos()
	if c.packageInfoStore == nil {
		return pi
	}
	return &readOnlyCacheClient[v1alpha1.PackageInfo, v1alpha1.PackageInfoList]{
		fallback: pi,
		store:    c.packageInfoStore,
		listFactory: func(items []v1alpha1.PackageInfo) v1alpha1.PackageInfoList {
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
		fallback: pr,
		readOnlyCacheClient: readOnlyCacheClient[v1alpha1.PackageRepository, v1alpha1.PackageRepositoryList]{
			fallback: pr,
			store:    c.packageRepositoryStore,
			listFactory: func(items []v1alpha1.PackageRepository) v1alpha1.PackageRepositoryList {
				return v1alpha1.PackageRepositoryList{Items: items}
			},
		},
	}
}

type readOnlyCacheClient[T any, L any] struct {
	fallback    readOnlyClientInterface[T, L]
	store       cache.Store
	listFactory func(items []T) L
	namespace   string
}

func getKey(namespace, name string) string {
	return cache.NewObjectName(namespace, name).String()
}

func (c *readOnlyCacheClient[T, L]) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.fallback.Watch(ctx, opts)
}

func (c *readOnlyCacheClient[T, L]) Get(ctx context.Context, name string, target *T) error {
	if obj, ok, err := c.store.GetByKey(getKey(c.namespace, name)); err != nil {
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
	keys := c.store.ListKeys()
	slices.Sort(keys)
	items := make([]T, len(keys))
	for i, key := range keys {
		if c.namespace != "" {
			if objName, err := cache.ParseObjectName(key); err != nil {
				return apierrors.NewInternalError(fmt.Errorf("bad cache key: %w", err))
			} else if objName.Namespace != c.namespace {
				continue
			}
		}

		if obj, exists, err := c.store.GetByKey(key); err != nil {
			return apierrors.NewInternalError(fmt.Errorf("resource not found: %w", err))
		} else if !exists {
			continue
		} else if obj, ok := obj.(*T); !ok {
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

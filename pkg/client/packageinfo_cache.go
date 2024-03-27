package client

import (
	"context"
	"errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

type cachedPackageInfoClient struct {
	PackageInfoInterface
	store cache.Store
}

func (c *cachedPackageInfoClient) Get(ctx context.Context, pkgInfoName string, result *v1alpha1.PackageInfo) error {
	if obj, ok, err := c.store.GetByKey(pkgInfoName); err != nil {
		return apierrors.NewInternalError(err)
	} else if !ok {
		return apierrors.NewNotFound(schema.GroupResource{}, pkgInfoName)
	} else if pkgInfo, ok := obj.(*v1alpha1.PackageInfo); !ok {
		return apierrors.NewInternalError(errors.New("not a packageinfo"))
	} else {
		*result = *pkgInfo
		return nil
	}
}

func (c *cachedPackageInfoClient) GetAll(ctx context.Context, result *v1alpha1.PackageInfoList) error {
	objs := c.store.List()
	items := make([]v1alpha1.PackageInfo, len(objs))
	for i, obj := range objs {
		if pkgInfo, ok := obj.(*v1alpha1.PackageInfo); !ok {
			return apierrors.NewInternalError(errors.New("not a packageinfo"))
		} else {
			items[i] = *pkgInfo
		}
	}
	*result = v1alpha1.PackageInfoList{Items: items}
	return nil
}

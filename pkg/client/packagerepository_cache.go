package client

import (
	"context"
	"errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"
)

type cachedPackageRepositoryClient struct {
	PackageRepositoryInterface
	store cache.Store
}

func (c *cachedPackageRepositoryClient) Get(
	ctx context.Context,
	repoName string,
	result *v1alpha1.PackageRepository,
) error {
	if obj, ok, err := c.store.GetByKey(repoName); err != nil {
		return apierrors.NewInternalError(err)
	} else if !ok {
		return c.PackageRepositoryInterface.Get(ctx, repoName, result)
	} else if repo, ok := obj.(*v1alpha1.PackageRepository); !ok {
		return apierrors.NewInternalError(errors.New("not a packagerepository"))
	} else {
		*result = *repo
		return nil
	}
}

func (c *cachedPackageRepositoryClient) GetAll(ctx context.Context, result *v1alpha1.PackageRepositoryList) error {
	objs := c.store.List()
	items := make([]v1alpha1.PackageRepository, len(objs))
	for i, obj := range objs {
		if repo, ok := obj.(*v1alpha1.PackageRepository); !ok {
			return apierrors.NewInternalError(errors.New("not a packagerepository"))
		} else {
			items[i] = *repo
		}
	}
	*result = v1alpha1.PackageRepositoryList{Items: items}
	return nil
}

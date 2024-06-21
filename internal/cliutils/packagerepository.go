package cliutils

import (
	"context"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

var NoDefaultRepo = fmt.Errorf("no default repo was found")

func GetDefaultRepo(ctx context.Context) (*v1alpha1.PackageRepository, error) {
	var repos v1alpha1.PackageRepositoryList

	client := PackageClient(ctx)
	if err := client.PackageRepositories().GetAll(ctx, &repos); err != nil {
		return nil, err
	}

	for _, repo := range repos.Items {
		if repo.IsDefaultRepository() {
			return &repo, nil
		}
	}

	return nil, NoDefaultRepo
}

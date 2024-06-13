package cliutils

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

func GetDefaultRepo(ctx context.Context) (*v1alpha1.PackageRepository, error) {
	var repos v1alpha1.PackageRepositoryList
	var defaultRepo v1alpha1.PackageRepository

	client := PackageClient(ctx)
	if err := client.PackageRepositories().GetAll(ctx, &repos); err != nil {
		return nil, err
	}

	for _, r := range repos.Items {
		if r.IsDefaultRepository() {
			defaultRepo = r
			break
		}
	}

	return &defaultRepo, nil
}

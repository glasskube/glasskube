package pkg_overview_btn

import (
	"fmt"

	"github.com/glasskube/glasskube/internal/web/util"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
)

const templateId = "clpkg-overview-btn"

type clpkgOverviewBtnInput struct {
	ButtonId        string
	PackageName     string
	Status          *client.PackageStatus
	Manifest        *v1alpha1.PackageManifest
	UpdateAvailable bool
	InDeletion      bool
	PackageHref     string
}

func getButtonId(pkgName string) string {
	return fmt.Sprintf("%v-%v", templateId, pkgName)
}

func ForClPkgOverviewBtn(packageWithStatus *list.PackageWithStatus, updateAvailable bool) *clpkgOverviewBtnInput {
	buttonId := getButtonId(packageWithStatus.Name)
	inDeletion := false
	if packageWithStatus.ClusterPackage != nil {
		inDeletion = !packageWithStatus.ClusterPackage.DeletionTimestamp.IsZero()
	}
	return &clpkgOverviewBtnInput{
		ButtonId:        buttonId,
		PackageName:     packageWithStatus.Name,
		Status:          packageWithStatus.Status,
		Manifest:        packageWithStatus.InstalledManifest,
		UpdateAvailable: updateAvailable,
		InDeletion:      inDeletion,
		PackageHref:     util.GetClusterPkgHref(packageWithStatus.Name),
	}
}

package pkg_overview_btn

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
)

const TemplateId = "pkg-overview-btn"

type pkgOverviewBtnInput struct {
	ButtonId        string
	PackageName     string
	Status          *client.PackageStatus
	Manifest        *v1alpha1.PackageManifest
	UpdateAvailable bool
}

func getButtonId(pkgName string) string {
	return fmt.Sprintf("%v-%v", TemplateId, pkgName)
}

func ForPkgOverviewBtn(packageWithStatus *list.PackageWithStatus, updateAvailable bool) *pkgOverviewBtnInput {
	buttonId := getButtonId(packageWithStatus.Name)
	return &pkgOverviewBtnInput{
		ButtonId:        buttonId,
		PackageName:     packageWithStatus.Name,
		Status:          packageWithStatus.Status,
		Manifest:        packageWithStatus.InstalledManifest,
		UpdateAvailable: updateAvailable,
	}
}

package pkg_overview_btn

import (
	"fmt"
	"html/template"
	"io"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
)

const TemplateId = "pkg-overview-btn"

type pkgOverviewBtnInput struct {
	ButtonId        string
	Swap            string
	PackageName     string
	Status          *client.PackageStatus
	Manifest        *v1alpha1.PackageManifest
	UpdateAvailable bool
}

func getButtonId(pkgName string) string {
	return fmt.Sprintf("%v-%v", TemplateId, pkgName)
}

func Render(w io.Writer, tmpl *template.Template, pkg *v1alpha1.Package, status *client.PackageStatus, manifest *v1alpha1.PackageManifest, updateAvailable bool) error {
	buttonId := getButtonId(pkg.Name)
	return tmpl.ExecuteTemplate(w, TemplateId, &pkgOverviewBtnInput{
		ButtonId:        buttonId,
		Swap:            fmt.Sprintf("outerHTML:#%s", buttonId),
		PackageName:     pkg.Name,
		Status:          status,
		Manifest:        manifest,
		UpdateAvailable: updateAvailable,
	})
}

func ForPkgOverviewBtn(packageWithStatus *list.PackageWithStatus, updateAvailable bool) *pkgOverviewBtnInput {
	buttonId := getButtonId(packageWithStatus.Name)
	return &pkgOverviewBtnInput{
		ButtonId:        buttonId,
		Swap:            "",
		PackageName:     packageWithStatus.Name,
		Status:          packageWithStatus.Status,
		Manifest:        packageWithStatus.InstalledManifest,
		UpdateAvailable: updateAvailable,
	}
}

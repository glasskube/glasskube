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
	ButtonId    string
	Swap        string
	PackageName string
	Status      *client.PackageStatus
	Manifest    *v1alpha1.PackageManifest
}

func getButtonId(pkgName string) string {
	return fmt.Sprintf("%v-%v", TemplateId, pkgName)
}

func Render(w io.Writer, tmpl *template.Template, pkgName string, status *client.PackageStatus, manifest *v1alpha1.PackageManifest) error {
	buttonId := getButtonId(pkgName)
	return tmpl.ExecuteTemplate(w, TemplateId, &pkgOverviewBtnInput{
		ButtonId:    buttonId,
		Swap:        fmt.Sprintf("outerHTML:#%s", buttonId),
		PackageName: pkgName,
		Status:      status,
		Manifest:    manifest,
	})
}

func ForPkgOverviewBtn(pkgTeaser *list.PackageWithStatus) *pkgOverviewBtnInput {
	buttonId := getButtonId(pkgTeaser.Name)
	return &pkgOverviewBtnInput{
		ButtonId:    buttonId,
		Swap:        "",
		PackageName: pkgTeaser.Name,
		Status:      pkgTeaser.Status,
		Manifest:    pkgTeaser.InstalledManifest,
	}
}

package pkg_detail_btns

import (
	"fmt"
	"html/template"
	"io"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
)

const TemplateId = "pkg-detail-btns"

type pkgDetailBtnsInput struct {
	ContainerId     string
	Swap            string
	PackageName     string
	Status          *client.PackageStatus
	Manifest        *v1alpha1.PackageManifest
	UpdateAvailable bool
}

func getId(pkgName string) string {
	return fmt.Sprintf("%v-%v", TemplateId, pkgName)
}

func getSwap(id string) string {
	return fmt.Sprintf("outerHTML:#%s", id)
}

func Render(w io.Writer, tmpl *template.Template, pkg *v1alpha1.Package, status *client.PackageStatus, manifest *v1alpha1.PackageManifest, latestVersion string) error {
	id := getId(pkg.Name)
	return tmpl.ExecuteTemplate(w, TemplateId, &pkgDetailBtnsInput{
		ContainerId:     id,
		Swap:            getSwap(id),
		PackageName:     pkg.Name,
		Status:          status,
		Manifest:        manifest,
		UpdateAvailable: latestVersion != "" && pkg.Spec.PackageInfo.Version != latestVersion,
	})
}

func ForPkgDetailBtns(
	pkgName string,
	status *client.PackageStatus,
	manifest *v1alpha1.PackageManifest,
	pkg *v1alpha1.Package,
	latestVersion string,
) *pkgDetailBtnsInput {
	id := getId(pkgName)
	return &pkgDetailBtnsInput{
		ContainerId:     id,
		Swap:            "",
		PackageName:     pkgName,
		Status:          status,
		Manifest:        manifest,
		UpdateAvailable: pkg != nil && pkg.Spec.PackageInfo.Version != latestVersion,
	}
}

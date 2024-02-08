package pkg_detail_btns

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
)

const TemplateId = "pkg-detail-btns"

type pkgDetailBtnsInput struct {
	ContainerId string
	Swap        string
	PackageName string
	Status      *client.PackageStatus
	Manifest    *v1alpha1.PackageManifest
}

func getId(pkgName string) string {
	return fmt.Sprintf("%v-%v", TemplateId, pkgName)
}

func getSwap(id string) string {
	return fmt.Sprintf("outerHTML:#%s", id)
}

func Render(w io.Writer, tmpl *template.Template, pkgName string, status *client.PackageStatus, manifest *v1alpha1.PackageManifest) {
	id := getId(pkgName)
	err := tmpl.ExecuteTemplate(w, TemplateId, &pkgDetailBtnsInput{
		ContainerId: id,
		Swap:        getSwap(id),
		PackageName: pkgName,
		Status:      status,
		Manifest:    manifest,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred rendering %v for %v: \n%v\n"+
			"This is most likely a BUG!", TemplateId, pkgName, err)
	}
}

func ForPkgDetailBtns(pkgName string, status *client.PackageStatus, manifest *v1alpha1.PackageManifest) *pkgDetailBtnsInput {
	id := getId(pkgName)
	return &pkgDetailBtnsInput{
		ContainerId: id,
		Swap:        "",
		PackageName: pkgName,
		Status:      status,
		Manifest:    manifest,
	}
}

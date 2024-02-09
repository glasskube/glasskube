package components

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
)

func GetButtonId(pkgName string) string {
	return fmt.Sprintf("install-%v", pkgName)
}

func GetSwap(buttonId string) string {
	return fmt.Sprintf("outerHTML:#%s", buttonId)
}

func RenderInstallButton(w io.Writer, tmpl *template.Template, pkgName string, status *client.PackageStatus, manifest *v1alpha1.PackageManifest) {
	buttonId := GetButtonId(pkgName)
	err := tmpl.ExecuteTemplate(w, "installbutton", &map[string]any{
		"ButtonId":    buttonId,
		"Swap":        GetSwap(buttonId),
		"PackageName": pkgName,
		"Status":      status,
		"Manifest":    manifest,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred rendering the install button for %v: \n%v\n"+
			"This is most likely a BUG!", pkgName, err)
	}
}

func ToInstallButtonInput(pkgTeaser list.PackageTeaserWithStatus) map[string]any {
	buttonId := GetButtonId(pkgTeaser.PackageName)
	return map[string]any{
		"ButtonId":    buttonId,
		"Swap":        GetSwap(buttonId),
		"PackageName": pkgTeaser.PackageName,
		"Status":      pkgTeaser.Status,
		"Manifest":    pkgTeaser.InstalledManifest,
	}
}

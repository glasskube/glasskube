package web

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/glasskube/glasskube/internal/web/components/alert"

	"github.com/glasskube/glasskube/api/v1alpha1"

	"github.com/glasskube/glasskube/internal/repo"

	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
	"github.com/glasskube/glasskube/internal/web/components/pkg_update_alert"
)

var (
	baseTemplate        *template.Template
	pkgsPageTmpl        *template.Template
	pkgPageTmpl         *template.Template
	supportPageTmpl     *template.Template
	bootstrapPageTmpl   *template.Template
	kubeconfigPageTmpl  *template.Template
	pkgOverviewBtnTmpl  *template.Template
	pkgDetailBtnsTmpl   *template.Template
	pkgInstallModalTmpl *template.Template
	pkgUpdateModalTmpl  *template.Template
	pkgUpdateAlertTmpl  *template.Template
	alertTmpl           *template.Template
	templatesDir        = "templates"
	componentsDir       = path.Join(templatesDir, "components")
	pagesDir            = path.Join(templatesDir, "pages")
)

func init() {
	templateFuncs := template.FuncMap{
		"ForPkgOverviewBtn": pkg_overview_btn.ForPkgOverviewBtn,
		"ForPkgDetailBtns":  pkg_detail_btns.ForPkgDetailBtns,
		"ForPkgUpdateAlert": pkg_update_alert.ForPkgUpdateAlert,
		"PackageManifestUrl": func(pkgName string, pkg *v1alpha1.Package, latestVersion string) string {
			var version string
			if pkg != nil && pkg.Spec.PackageInfo.Version != "" {
				version = pkg.Spec.PackageInfo.Version
			} else {
				version = latestVersion
			}
			if url, err := repo.GetPackageManifestURL("", pkgName, version); err != nil {
				return ""
			} else {
				return url
			}
		},
		"ForAlert": alert.ForAlert,
	}
	baseTemplate = template.Must(template.New("base.html").
		Funcs(templateFuncs).
		ParseFS(embededFs, path.Join(templatesDir, "layout", "base.html")))
	pkgsPageTmpl = pageTmpl("packages.html")
	pkgPageTmpl = pageTmpl("package.html")
	supportPageTmpl = pageTmpl("support.html")
	bootstrapPageTmpl = pageTmpl("bootstrap.html")
	kubeconfigPageTmpl = pageTmpl("kubeconfig.html")
	pkgOverviewBtnTmpl = componentTmpl(pkg_overview_btn.TemplateId, "pkg-overview-btn.html")
	pkgDetailBtnsTmpl = componentTmpl(pkg_detail_btns.TemplateId, "pkg-detail-btns.html")
	pkgUpdateAlertTmpl = componentTmpl(pkg_update_alert.TemplateId, "pkg-update-alert.html")
	pkgInstallModalTmpl = componentTmpl("pkg-install-modal", "pkg-install-modal.html")
	pkgUpdateModalTmpl = componentTmpl("pkg-update-modal", "pkg-update-modal.html")
	alertTmpl = componentTmpl("alert", "alert.html")
}

func pageTmpl(fileName string) *template.Template {
	return template.Must(
		template.Must(baseTemplate.Clone()).ParseFS(
			embededFs,
			path.Join(pagesDir, fileName),
			path.Join(componentsDir, "*.html")))
}

func componentTmpl(id string, fileName string) *template.Template {
	return template.Must(
		template.New(id).ParseFS(
			embededFs,
			path.Join(componentsDir, fileName)))
}

func checkTmplError(e error, tmplName string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "\nUnexpected error rendering %v: %v\n – This is most likely a BUG – "+
			"Please report it here: https://github.com/glasskube/glasskube\n\n", tmplName, e)
	}
}

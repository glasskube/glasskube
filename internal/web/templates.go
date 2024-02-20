package web

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
)

var (
	baseTemplate                *template.Template
	pkgsPageTmpl                *template.Template
	pkgPageTmpl                 *template.Template
	supportPageTmpl             *template.Template
	bootstrapPageTmpl           *template.Template
	pkgOverviewBtnTmpl          *template.Template
	pkgDetailBtnsTmpl           *template.Template
	pkgInstallModalTmpl         *template.Template
	pkgInstallModalVersionsTmpl *template.Template
	templatesDir                = "templates"
	componentsDir               = path.Join(templatesDir, "components")
	pagesDir                    = path.Join(templatesDir, "pages")
)

func loadTemplates() {
	templateFuncs := template.FuncMap{
		"ForPkgOverviewBtn": pkg_overview_btn.ForPkgOverviewBtn,
		"ForPkgDetailBtns":  pkg_detail_btns.ForPkgDetailBtns,
		"PackageManifestUrl": func(pkgName string) string {
			// TODO get configured repository URL instead
			return fmt.Sprintf("https://github.com/glasskube/packages/blob/main/packages/%s/package.yaml", pkgName)
		},
	}
	baseTemplate = template.Must(template.New("base.html").
		Funcs(templateFuncs).
		ParseFS(embededFs, path.Join(templatesDir, "layout", "base.html")))
	pkgsPageTmpl = template.Must(
		template.Must(baseTemplate.Clone()).
			ParseFS(embededFs, path.Join(pagesDir, "packages.html"), path.Join(componentsDir, "*.html")))
	pkgPageTmpl = template.Must(template.Must(baseTemplate.Clone()).
		ParseFS(embededFs, path.Join(pagesDir, "package.html"), path.Join(componentsDir, "*.html")))
	supportPageTmpl = template.Must(template.Must(baseTemplate.Clone()).
		ParseFS(embededFs, path.Join(pagesDir, "support.html"), path.Join(componentsDir, "*.html")))
	bootstrapPageTmpl = template.Must(template.Must(baseTemplate.Clone()).
		ParseFS(embededFs, path.Join(pagesDir, "bootstrap.html"), path.Join(componentsDir, "*.html")))
	pkgOverviewBtnTmpl = template.Must(template.New(pkg_overview_btn.TemplateId).
		ParseFS(embededFs, path.Join(componentsDir, "pkg-overview-btn.html")))
	pkgDetailBtnsTmpl = template.Must(template.New(pkg_detail_btns.TemplateId).
		ParseFS(embededFs, path.Join(componentsDir, "pkg-detail-btns.html")))
	pkgInstallModalTmpl = template.Must(template.New("pkg-install-modal").
		ParseFS(embededFs, path.Join(componentsDir, "pkg-install-modal.html")))
	pkgInstallModalVersionsTmpl = template.Must(template.New("pkg-install-modal-versions").
		ParseFS(embededFs, path.Join(componentsDir, "pkg-install-modal-versions.html")))
}

func checkTmplError(e error, tmplName string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "\nUnexpected error rendering %v: %v\n – This is most likely a BUG – "+
			"Please report it here: https://github.com/glasskube/glasskube\n\n", tmplName, e)
	}
}

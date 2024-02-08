package web

import (
	"fmt"
	"html/template"
	"path"

	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
)

var (
	baseTemplate       *template.Template
	pkgsPageTmpl       *template.Template
	pkgPageTmpl        *template.Template
	supportPageTmpl    *template.Template
	pkgOverviewBtnTmpl *template.Template
	pkgDetailBtnsTmpl  *template.Template
	templatesDir       = "templates"
	componentsDir      = path.Join(templatesDir, "components")
	pagesDir           = path.Join(templatesDir, "pages")
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
	pkgOverviewBtnTmpl = template.Must(template.New(pkg_overview_btn.TemplateId).
		ParseFS(embededFs, path.Join(componentsDir, "pkg-overview-btn.html")))
	pkgDetailBtnsTmpl = template.Must(template.New(pkg_detail_btns.TemplateId).
		ParseFS(embededFs, path.Join(componentsDir, "pkg-detail-btns.html")))
}

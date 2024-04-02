package web

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/glasskube/glasskube/internal/web/components/pkg_config_input"

	"github.com/fsnotify/fsnotify"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/internal/web/components/alert"
	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
	"github.com/glasskube/glasskube/internal/web/components/pkg_update_alert"
	"go.uber.org/multierr"
)

var (
	templateFuncs         template.FuncMap
	baseTemplate          *template.Template
	pkgsPageTmpl          *template.Template
	pkgPageTmpl           *template.Template
	supportPageTmpl       *template.Template
	bootstrapPageTmpl     *template.Template
	kubeconfigPageTmpl    *template.Template
	pkgOverviewBtnTmpl    *template.Template
	pkgDetailBtnsTmpl     *template.Template
	pkgUpdateModalTmpl    *template.Template
	pkgUpdateAlertTmpl    *template.Template
	pkgConfigInput        *template.Template
	pkgUninstallModalTmpl *template.Template
	alertTmpl             *template.Template
	templatesBaseDir      = "internal/web"
	templatesDir          = "templates"
	componentsDir         = path.Join(templatesDir, "components")
	pagesDir              = path.Join(templatesDir, "pages")
)

func watchTemplates() error {
	watcher, err := fsnotify.NewWatcher()
	err = multierr.Combine(
		err,
		watcher.Add(path.Join(templatesBaseDir, componentsDir)),
		watcher.Add(path.Join(templatesBaseDir, templatesDir, "layout")),
		watcher.Add(path.Join(templatesBaseDir, pagesDir)),
	)
	if err == nil {
		go func() {
			for range watcher.Events {
				parseTemplates()
			}
		}()
	}
	return err
}

func parseTemplates() {
	templateFuncs = template.FuncMap{
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
		"ForAlert":          alert.ForAlert,
		"ForPkgConfigInput": pkg_config_input.ForPkgConfigInput,
		"IsUpgradable":      semver.IsUpgradable,
	}

	baseTemplate = template.Must(template.New("base.html").
		Funcs(templateFuncs).
		ParseFS(webFs, path.Join(templatesDir, "layout", "base.html")))
	pkgsPageTmpl = pageTmpl("packages.html")
	pkgPageTmpl = pageTmpl("package.html")
	supportPageTmpl = pageTmpl("support.html")
	bootstrapPageTmpl = pageTmpl("bootstrap.html")
	kubeconfigPageTmpl = pageTmpl("kubeconfig.html")
	pkgOverviewBtnTmpl = componentTmpl(pkg_overview_btn.TemplateId, "pkg-overview-btn.html")
	pkgDetailBtnsTmpl = componentTmpl(pkg_detail_btns.TemplateId, "pkg-detail-btns.html")
	pkgUpdateAlertTmpl = componentTmpl(pkg_update_alert.TemplateId, "pkg-update-alert.html")
	pkgUpdateModalTmpl = componentTmpl("pkg-update-modal", "pkg-update-modal.html")
	pkgConfigInput = componentTmpl("pkg-config-input", "pkg-config-input.html")
	pkgUninstallModalTmpl = componentTmpl("pkg-uninstall-modal", "pkg-uninstall-modal.html")
	alertTmpl = componentTmpl("alert", "alert.html")
	componentTmpl("version-mismatch-warning", "version-mismatch-warning.html")
}

func pageTmpl(fileName string) *template.Template {
	return template.Must(
		template.Must(baseTemplate.Clone()).ParseFS(
			webFs,
			path.Join(pagesDir, fileName),
			path.Join(componentsDir, "*.html")))
}

func componentTmpl(id string, fileName string) *template.Template {
	return template.Must(
		template.New(id).Funcs(templateFuncs).ParseFS(
			webFs,
			path.Join(componentsDir, fileName)))
}

func checkTmplError(e error, tmplName string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "\nUnexpected error rendering %v: %v\n – This is most likely a BUG – "+
			"Please report it here: https://github.com/glasskube/glasskube\n\n", tmplName, e)
	}
}

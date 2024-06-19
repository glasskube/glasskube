package web

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/internal/web/components/alert"
	"github.com/glasskube/glasskube/internal/web/components/datalist"
	"github.com/glasskube/glasskube/internal/web/components/pkg_config_input"
	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
	"github.com/glasskube/glasskube/internal/web/components/pkg_update_alert"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type templates struct {
	templateFuncs           template.FuncMap
	baseTemplate            *template.Template
	clusterPkgsPageTemplate *template.Template
	pkgsPageTmpl            *template.Template
	pkgPageTmpl             *template.Template
	pkgDiscussionPageTmpl   *template.Template
	supportPageTmpl         *template.Template
	bootstrapPageTmpl       *template.Template
	kubeconfigPageTmpl      *template.Template
	settingsPageTmpl        *template.Template
	pkgUpdateModalTmpl      *template.Template
	pkgConfigInput          *template.Template
	pkgConfigAdvancedTmpl   *template.Template
	pkgUninstallModalTmpl   *template.Template
	alertTmpl               *template.Template
	datalistTmpl            *template.Template
	pkgDiscussionBadgeTmpl  *template.Template
	repoClientset           repoclient.RepoClientset
}

var (
	templatesBaseDir = "internal/web"
	templatesDir     = "templates"
	componentsDir    = path.Join(templatesDir, "components")
	pagesDir         = path.Join(templatesDir, "pages")
)

func (t *templates) watchTemplates() error {
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
				t.parseTemplates()
			}
		}()
	}
	return err
}

func (t *templates) parseTemplates() {
	t.templateFuncs = template.FuncMap{
		"ForClPkgOverviewBtn": pkg_overview_btn.ForClPkgOverviewBtn,
		"ForPkgDetailBtns":    pkg_detail_btns.ForPkgDetailBtns,
		"ForPkgUpdateAlert":   pkg_update_alert.ForPkgUpdateAlert,
		"PackageManifestUrl": func(pkg ctrlpkg.Package) string {
			if !pkg.IsNil() {
				url, err := t.repoClientset.ForPackage(pkg).
					GetPackageManifestURL(pkg.GetName(), pkg.GetSpec().PackageInfo.Version)
				if err == nil {
					return url
				}
			}
			return ""
		},
		"ForAlert":          alert.ForAlert,
		"ForPkgConfigInput": pkg_config_input.ForPkgConfigInput,
		"ForDatalist":       datalist.ForDatalist,
		"IsUpgradable":      semver.IsUpgradable,
		"Markdown": func(source string) template.HTML {
			var buf bytes.Buffer

			converter := goldmark.New(
				goldmark.WithParserOptions(
					parser.WithASTTransformers(
						util.Prioritized(&ASTTransformer{}, 1000),
					),
				),
			)

			if err := converter.Convert([]byte(source), &buf); err != nil {
				return template.HTML("<p>" + source + "</p>")
			}

			return template.HTML(buf.String())
		},
		"Reversed": func(param any) any {
			kind := reflect.TypeOf(param).Kind()
			switch kind {
			case reflect.Slice, reflect.Array:
				val := reflect.ValueOf(param)

				ln := val.Len()
				newVal := make([]interface{}, ln)
				for i := 0; i < ln; i++ {
					newVal[ln-i-1] = val.Index(i).Interface()
				}

				return newVal
			default:
				return param
			}
		},
		"UrlEscape": func(param string) string {
			return template.URLQueryEscaper(param)
		},
		"IsRepoStatusReady": func(repo v1alpha1.PackageRepository) bool {
			cond := meta.FindStatusCondition(repo.Status.Conditions, string(condition.Ready))
			return cond != nil && cond.Status == metav1.ConditionTrue
		},
		"PackageDetailRefreshId": func(manifest *v1alpha1.PackageManifest, pkg ctrlpkg.Package) string {
			var id string
			if manifest.Scope.IsCluster() {
				id = manifest.Name
			} else if !pkg.IsNil() {
				id = fmt.Sprintf("%s-%s", pkg.GetNamespace(), pkg.GetName())
			}
			return fmt.Sprintf("refresh-pkg-detail-%s", id)
		},
	}

	t.baseTemplate = template.Must(template.New("base.html").
		Funcs(t.templateFuncs).
		ParseFS(webFs, path.Join(templatesDir, "layout", "base.html")))
	t.clusterPkgsPageTemplate = t.pageTmpl("clusterpackages.html")
	t.pkgsPageTmpl = t.pageTmpl("packages.html")
	t.pkgPageTmpl = t.pageTmpl("package.html")
	t.pkgDiscussionPageTmpl = t.pageTmpl("discussion.html")
	t.supportPageTmpl = t.pageTmpl("support.html")
	t.bootstrapPageTmpl = t.pageTmpl("bootstrap.html")
	t.kubeconfigPageTmpl = t.pageTmpl("kubeconfig.html")
	t.settingsPageTmpl = t.pageTmpl("settings.html")
	t.pkgUpdateModalTmpl = t.componentTmpl("pkg-update-modal")
	t.pkgConfigInput = t.componentTmpl("pkg-config-input", "datalist")
	t.pkgConfigAdvancedTmpl = t.componentTmpl("pkg-config-advanced")
	t.pkgUninstallModalTmpl = t.componentTmpl("pkg-uninstall-modal")
	t.alertTmpl = t.componentTmpl("alert")
	t.datalistTmpl = t.componentTmpl("datalist")
	t.pkgDiscussionBadgeTmpl = t.componentTmpl("discussion-badge")
}

func (t *templates) pageTmpl(fileName string) *template.Template {
	return template.Must(
		template.Must(t.baseTemplate.Clone()).ParseFS(
			webFs,
			path.Join(pagesDir, fileName),
			path.Join(componentsDir, "*.html")))
}

func (t *templates) componentTmpl(id string, requiredTemplates ...string) *template.Template {
	tpls := make([]string, 0)
	for _, requiredTmpl := range requiredTemplates {
		tpls = append(tpls, path.Join(componentsDir, requiredTmpl+".html"))
	}
	tpls = append(tpls, path.Join(componentsDir, id+".html"))
	return template.Must(
		template.New(id).Funcs(t.templateFuncs).ParseFS(
			webFs,
			tpls...))
}

func checkTmplError(e error, tmplName string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "\nUnexpected error rendering %v: %v\n – This is most likely a BUG – "+
			"Please report it here: https://github.com/glasskube/glasskube\n\n", tmplName, e)
	}
}

type ASTTransformer struct{}

func (g *ASTTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch v := n.(type) {
		case *ast.Link:
			v.SetAttributeString("target", "_blank")
			v.SetAttributeString("rel", "noopener noreferrer")
		}

		return ast.WalkContinue, nil
	})
}

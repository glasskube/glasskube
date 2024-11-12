package responder

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"reflect"
	"strings"

	"github.com/glasskube/glasskube/internal/web/components"

	depUtil "github.com/glasskube/glasskube/internal/dependency/util"

	webutil "github.com/glasskube/glasskube/internal/web/sse/refresh"

	"github.com/fsnotify/fsnotify"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/internal/web/components/pkg_detail_btns"
	"github.com/glasskube/glasskube/internal/web/components/pkg_overview_btn"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type templates struct {
	templateFuncs template.FuncMap
	baseTemplate  *template.Template
	fs            fs.FS
}

var (
	BaseDir      = "internal/web"
	templatesDir = "templates"
)

func (t *templates) WatchTemplates() error {
	watcher, err := fsnotify.NewWatcher()
	err = multierr.Combine(
		err,
		watcher.Add(path.Join(BaseDir, templatesDir, "components")),
		watcher.Add(path.Join(BaseDir, templatesDir, "layout")),
		watcher.Add(path.Join(BaseDir, templatesDir, "pages")),
	)
	if err == nil {
		go func() {
			for range watcher.Events {
				t.ParseTemplates()
			}
		}()
	}
	return err
}

func (t *templates) ParseTemplates() {
	t.templateFuncs = template.FuncMap{
		"ForClPkgOverviewBtn": pkg_overview_btn.ForClPkgOverviewBtn,
		"ForPkgDetailBtns":    pkg_detail_btns.ForPkgDetailBtns,
		"ForToast":            toast.ForToast,
		"ForPkgConfigInput":   components.ForPkgConfigInput,
		"ForDatalist":         components.ForDatalist,
		"IsUpgradable":        semver.IsUpgradable,
		"Markdown": func(source string) template.HTML {
			var buf bytes.Buffer

			converter := goldmark.New(
				goldmark.WithExtensions(
					extension.Linkify,
				),
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
		"PackageDetailRefreshId":          webutil.PackageRefreshDetailId,
		"PackageDetailHeaderRefreshId":    webutil.PackageRefreshDetailHeaderId,
		"PackageOverviewRefreshId":        webutil.PackageOverviewRefreshId,
		"ClusterPackageOverviewRefreshId": webutil.ClusterPackageOverviewRefreshId,
		"ComponentName":                   depUtil.ComponentName,
		"AutoUpdateEnabled": func(pkg ctrlpkg.Package) bool {
			if pkg != nil && !pkg.IsNil() {
				return pkg.AutoUpdatesEnabled()
			}
			return false
		},
		"IsSuspended": func(pkg ctrlpkg.Package) bool {
			if pkg != nil && !pkg.IsNil() {
				return pkg.GetSpec().Suspend
			}
			return false
		},
	}

	var paths []string
	err := fs.WalkDir(t.fs, templatesDir, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(d.Name(), ".html") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to walk templates dir: %v", templatesDir))
	}

	t.baseTemplate = template.Must(
		template.New("root").
			Funcs(t.templateFuncs).
			ParseFS(t.fs, paths...))
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
		case *ast.Blockquote:
			v.SetAttributeString("class", "border-start border-primary border-3 ps-2")
		}

		return ast.WalkContinue, nil
	})
}

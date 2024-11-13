package handlers

import (
	"fmt"
	"net/http"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/adapter/goclient"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/dependency/graph"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/internal/web/types"
	"github.com/glasskube/glasskube/internal/web/util"
	"github.com/glasskube/glasskube/pkg/uninstall"
)

type uninstallModalData struct {
	types.TemplateContextHolder
	PackageName     string
	Namespace, Name string
	Pruned          []graph.PackageRef
	PackageHref     string
	ShownError      error
}

func GetUninstallClusterPackage(w http.ResponseWriter, r *http.Request) {
	handleUninstallModal(w, r)
}

func PostUninstallClusterPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pCtx := getPackageContext(r).request
	pkgClient := clicontext.PackageClientFromContext(ctx)
	uninstaller := uninstall.NewUninstaller(pkgClient)

	var pkg v1alpha1.ClusterPackage
	if err := pkgClient.ClusterPackages().Get(ctx, pCtx.manifestName, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch clusterpackage %v: %w", pCtx.manifestName, err)))
		return
	}
	if err := uninstaller.Uninstall(ctx, &pkg, false); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to uninstall clusterpackage %v: %w", pCtx.manifestName, err)))
		return
	}
}

func GetUninstallPackage(w http.ResponseWriter, r *http.Request) {
	handleUninstallModal(w, r)
}

func PostUninstallPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pCtx := getPackageContext(r).request
	pkgClient := clicontext.PackageClientFromContext(ctx)
	uninstaller := uninstall.NewUninstaller(pkgClient)

	var pkg v1alpha1.Package
	if err := pkgClient.Packages(pCtx.namespace).Get(ctx, pCtx.name, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch package %v/%v: %w", pCtx.namespace, pCtx.name, err)))
		return
	}
	if err := uninstaller.Uninstall(ctx, &pkg, false); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to uninstall package %v/%v: %w", pCtx.namespace, pCtx.name, err)))
		return
	}
}

func handleUninstallModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pCtx := getPackageContext(r).request
	pkgClient := clicontext.PackageClientFromContext(ctx)
	repoClientset := clicontext.RepoClientsetFromContext(ctx)
	dependencyMgr := dependency.NewDependencyManager(goclient.NewPackageClientAdapter(pkgClient), repoClientset)

	var pruned []graph.PackageRef
	var validationErr error
	var pkgHref string
	if g, err := dependencyMgr.NewGraph(r.Context()); err != nil {
		validationErr = fmt.Errorf("error validating uninstall: %w", err)
	} else {
		if pCtx.namespace == "" && pCtx.name == "" {
			pkgHref = util.GetClusterPkgHref(pCtx.manifestName)
			g.Delete(pCtx.manifestName, "")
		} else {
			pkgHref = util.GetNamespacedPkgHref(pCtx.manifestName, pCtx.namespace, pCtx.name)
			g.Delete(pCtx.name, pCtx.namespace)
		}
		pruned = g.Prune()
		if err := g.Validate(); err != nil {
			validationErr = fmt.Errorf("%v cannot be uninstalled: %w", pCtx.manifestName, err)
		}
	}
	responder.SendComponent(w, r, "components/pkg-uninstall-modal",
		responder.ContextualizedTemplate(&uninstallModalData{
			PackageName: pCtx.manifestName,
			Namespace:   pCtx.namespace,
			Name:        pCtx.name,
			Pruned:      pruned,
			PackageHref: pkgHref,
			ShownError:  validationErr,
		}))
}

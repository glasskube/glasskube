package handlers

import (
	"fmt"
	"net/http"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	webopen "github.com/glasskube/glasskube/internal/web/open"
	"github.com/glasskube/glasskube/internal/web/responder"
)

func PostOpenClusterPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)
	p := getPackageContext(r).request

	var pkg v1alpha1.ClusterPackage
	if err := pkgClient.ClusterPackages().Get(ctx, p.manifestName, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch clusterpackage %v: %w", p.manifestName, err)))
		return
	}
	if err := webopen.HandleOpen(ctx, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to open %v: %w", pkg.GetName(), err)))
		return
	}
}

func PostOpenPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)
	p := getPackageContext(r).request

	var pkg v1alpha1.Package
	if err := pkgClient.Packages(p.namespace).Get(ctx, p.name, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch package %v/%v: %w", p.namespace, p.name, err)))
		return
	}
	if err := webopen.HandleOpen(ctx, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to open %v: %w", pkg.GetName(), err)))
		return
	}
}

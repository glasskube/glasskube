package handlers

import (
	"fmt"
	"net/http"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/responder"
	webutil "github.com/glasskube/glasskube/internal/web/util"
	"github.com/glasskube/glasskube/pkg/suspend"
)

func PostSuspend(w http.ResponseWriter, r *http.Request) {
	var options suspend.Options
	gitopsModeEnabled := webutil.IsGitopsModeEnabled(r)
	if gitopsModeEnabled {
		options = append(options, suspend.DryRun())
	}

	if pkg, err := getPackageFromRequest(r); err != nil {
		responder.SendToast(w, toast.WithErr(err))
	} else if suspended, err := suspend.Suspend(r.Context(), pkg, options...); err != nil {
		responder.SendToast(w, toast.WithErr(err))
	} else if suspended {
		if gitopsModeEnabled {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				responder.SendYamlModal(w, yamlOutput, nil)
			}
		} else {
			responder.SendToast(w, toast.WithMessage(fmt.Sprintf("%v has been suspended", pkg.GetName())),
				toast.WithSeverity(toast.Info))
		}
	} else {
		responder.SendToast(w, toast.WithMessage(fmt.Sprintf("%v was already suspended", pkg.GetName())),
			toast.WithSeverity(toast.Info))
	}
}

func PostResume(w http.ResponseWriter, r *http.Request) {
	var options suspend.Options
	gitopsModeEnabled := webutil.IsGitopsModeEnabled(r)
	if gitopsModeEnabled {
		options = append(options, suspend.DryRun())
	}

	if pkg, err := getPackageFromRequest(r); err != nil {
		responder.SendToast(w, toast.WithErr(err))
	} else if resumed, err := suspend.Resume(r.Context(), pkg, options...); err != nil {
		responder.SendToast(w, toast.WithErr(err))
	} else if resumed {
		if gitopsModeEnabled {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				responder.SendYamlModal(w, yamlOutput, nil)
			}
		} else {
			responder.SendToast(w, toast.WithMessage(fmt.Sprintf("%v has been resumed", pkg.GetName())))
		}
	} else {
		responder.SendToast(w, toast.WithMessage(fmt.Sprintf("%v was not suspended", pkg.GetName())),
			toast.WithSeverity(toast.Info))
	}
}

func getPackageFromRequest(r *http.Request) (ctrlpkg.Package, error) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)
	pCtx := getPackageContext(r).request

	var pkg ctrlpkg.Package
	if pCtx.namespace == "" && pCtx.name == "" {
		var cp v1alpha1.ClusterPackage
		if err := pkgClient.ClusterPackages().Get(ctx, pCtx.manifestName, &cp); err != nil {
			return nil, err
		} else {
			pkg = &cp
		}
	} else {
		var p v1alpha1.Package
		if err := pkgClient.Packages(pCtx.namespace).Get(r.Context(), pCtx.name, &p); err != nil {
			return nil, err
		} else {
			pkg = &p
		}
	}
	return pkg, nil
}

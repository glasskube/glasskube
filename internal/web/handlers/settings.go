package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/cookie"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/internal/web/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type settingsPageData struct {
	types.TemplateContextHolder
	Repositories    []v1alpha1.PackageRepository
	AdvancedOptions bool
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)

	var repos v1alpha1.PackageRepositoryList
	if err := pkgClient.PackageRepositories().GetAll(ctx, &repos); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
		return
	}

	advancedOptions, err := cookie.GetAdvancedOptionsFromCookie(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get advanced options from cookie: %v\n", err)
	}

	responder.SendPage(w, r, "pages/settings",
		responder.ContextualizedTemplate(&settingsPageData{
			Repositories:    repos.Items,
			AdvancedOptions: advancedOptions,
		}))
}

func PostSettings(w http.ResponseWriter, r *http.Request) {
	formVal := r.FormValue(cookie.AdvancedOptionsKey)
	cookie.SetAdvancedOptionsCookie(w, formVal == "on")
}

type repositoryPageData struct {
	types.TemplateContextHolder
	Repository v1alpha1.PackageRepository
}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)

	repoName := r.PathValue("repoName")
	var repo v1alpha1.PackageRepository
	if err := pkgClient.PackageRepositories().Get(ctx, repoName, &repo); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
		return
	}
	responder.SendPage(w, r, "pages/repository",
		responder.ContextualizedTemplate(&repositoryPageData{
			Repository: repo,
		}))
}

func PostRepository(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)

	repoName := r.PathValue("repoName")
	repoUrl := r.FormValue("url")
	checkDefault := r.FormValue("default")
	opts := v1.UpdateOptions{}
	var repo v1alpha1.PackageRepository
	var defaultRepo *v1alpha1.PackageRepository
	var err error

	if err := pkgClient.PackageRepositories().Get(ctx, repoName, &repo); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
		return
	}

	if repoUrl != "" {
		if _, err := url.ParseRequestURI(repoUrl); err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("use a valid URL for the package repository (got %v)", err)))
			return
		}
		repo.Spec.Url = repoUrl
	}

	repo.Spec.Auth = nil

	if checkDefault == "on" {
		// TODO "cliutils" in the server??
		defaultRepo, err = cliutils.GetDefaultRepo(ctx)
		if errors.Is(err, cliutils.NoDefaultRepo) {
			repo.SetDefaultRepository()
		} else if err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
			return
		} else if defaultRepo.Name != repoName {
			defaultRepo.SetDefaultRepositoryBool(false)
			if err := pkgClient.PackageRepositories().Update(ctx, defaultRepo, opts); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf(" error updating current default package repository: %v", err)))
				return
			}
			repo.SetDefaultRepository()
		}
	}

	if err := pkgClient.PackageRepositories().Update(ctx, &repo, opts); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf(" error updating the package repository: %v", err)))
		if checkDefault == "on" && defaultRepo != nil && defaultRepo.Name != repoName {
			defaultRepo.SetDefaultRepositoryBool(true)
			if err := pkgClient.PackageRepositories().Update(ctx, defaultRepo, opts); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf(" error rolling back to default package repository: %v", err)))
			}
		}
		return
	}
	responder.Redirect(w, "/settings")
}

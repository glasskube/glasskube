package settings

import (
	"errors"
	"fmt"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/cookie"
	"github.com/glasskube/glasskube/internal/web/templates"
	"github.com/glasskube/glasskube/internal/web/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/url"
	"os"
)

type settingsPage struct {
	templates.Page
	Repositories    []v1alpha1.PackageRepository
	AdvancedOptions bool
}

func Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /settings", getSettings)
	m.HandleFunc("POST /settings", postSettings)
	m.HandleFunc("GET /settings/repository/{repoName}", getRepository)
	m.HandleFunc("GET /settings/repository/{repoName}", postRepository)
	return m
}

func getSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)

	var repos v1alpha1.PackageRepositoryList
	if err := pkgClient.PackageRepositories().GetAll(ctx, &repos); err != nil {
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
		return
	}

	advancedOptions, err := cookie.GetAdvancedOptionsFromCookie(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get advanced options from cookie: %v\n", err)
	}

	tmplErr := templates.Templates.SettingsPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Repositories":    repos.Items,
		"AdvancedOptions": advancedOptions,
	}, nil))
	util.CheckTmplError(tmplErr, "settings")
}

func postSettings(w http.ResponseWriter, r *http.Request) {
	formVal := r.FormValue(cookie.AdvancedOptionsKey)
	cookie.SetAdvancedOptionsCookie(w, formVal == "on")
}

type repositoryPage struct {
	templates.Page
	Repository v1alpha1.PackageRepository
}

func getRepository(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)

	repoName := r.PathValue("repoName")
	var repo v1alpha1.PackageRepository
	if err := pkgClient.PackageRepositories().Get(ctx, repoName, &repo); err != nil {
		// error handling
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
		return
	}
	tmplErr := templates.Templates.RepositoryPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Repository": repo,
	}, nil))
	util.CheckTmplError(tmplErr, "repository")
}

func postRepository(w http.ResponseWriter, r *http.Request) {
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
		s.sendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
		return
	}

	if repoUrl != "" {
		if _, err := url.ParseRequestURI(repoUrl); err != nil {
			s.sendToast(w, toast.WithErr(fmt.Errorf("use a valid URL for the package repository (got %v)", err)))
			return
		}
		repo.Spec.Url = repoUrl
	}

	repo.Spec.Auth = nil

	if checkDefault == "on" {
		defaultRepo, err = cliutils.GetDefaultRepo(ctx)
		if errors.Is(err, cliutils.NoDefaultRepo) {
			repo.SetDefaultRepository()
		} else if err != nil {
			s.sendToast(w, toast.WithErr(fmt.Errorf("failed to fetch repositories: %w", err)))
			return
		} else if defaultRepo.Name != repoName {
			defaultRepo.SetDefaultRepositoryBool(false)
			if err := pkgClient.PackageRepositories().Update(ctx, defaultRepo, opts); err != nil {
				s.sendToast(w, toast.WithErr(fmt.Errorf(" error updating current default package repository: %v", err)))
				return
			}
			repo.SetDefaultRepository()
		}
	}

	if err := pkgClient.PackageRepositories().Update(ctx, &repo, opts); err != nil {
		s.sendToast(w, toast.WithErr(fmt.Errorf(" error updating the package repository: %v", err)))
		if checkDefault == "on" && defaultRepo != nil && defaultRepo.Name != repoName {
			defaultRepo.SetDefaultRepositoryBool(true)
			if err := pkgClient.PackageRepositories().Update(ctx, defaultRepo, opts); err != nil {
				s.sendToast(w, toast.WithErr(fmt.Errorf(" error rolling back to default package repository: %v", err)))
			}
		}
		return
	}
	s.swappingRedirect(w, "/settings", "main", "main")
}

package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/glasskube/glasskube/internal/giscus"
	"github.com/glasskube/glasskube/internal/httperror"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/api/errors"
)

// packageDiscussion is a full page for showing various discussions, reactions, etc.
func (s *server) packageDiscussion(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		githubUrl := r.FormValue("githubUrl")
		telemetry.SetUserProperty("github_url", githubUrl)
		return
	}
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := mux.Vars(r)["repositoryName"]
	pkg, manifest, err := describe.DescribeInstalledPackage(r.Context(), pkgName)
	if err != nil && !errors.IsNotFound(err) {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching installed package %v", pkgName), "danger")
		return
	} else if err != nil {
		// implies that the package is not installed
		err = nil
	}

	var idx repo.PackageIndex
	if err := s.repoClientset.ForRepoWithName(repositoryName).FetchPackageIndex(pkgName, &idx); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred fetching versions of "+pkgName, "danger")
		return
	}

	if manifest == nil {
		manifest = &v1alpha1.PackageManifest{}
		if err := s.repoClientset.ForRepoWithName(repositoryName).
			FetchPackageManifest(pkgName, idx.LatestVersion, manifest); err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching manifest of %v in version %v in repository %v",
					pkgName, idx.LatestVersion, repositoryName), "danger")
			return
		}
	}

	err = s.templates.pkgDiscussionPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Giscus":             giscus.Client().Config,
		"Package":            pkg,
		"Status":             client.GetStatusOrPending(pkg),
		"Manifest":           manifest,
		"LatestVersion":      idx.LatestVersion,
		"UpdateAvailable":    pkg != nil && s.isUpdateAvailable(r.Context(), pkgName),
		"ShowDiscussionLink": true,
	}, err))
	checkTmplError(err, fmt.Sprintf("package-detail (%s)", pkgName))
}

func (s *server) discussionBadge(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkgName"]

	var totalCount int
	if counts, err := giscus.Client().GetCountsFor(pkgName); err != nil {
		if !httperror.IsNotFound(err) {
			fmt.Fprintf(os.Stderr, "failed to get discussion counts from giscus: %v\n", err)
		}
	} else {
		totalCount = counts.ReactionCount + counts.TotalCommentCount + counts.TotalReplyCount
	}

	var err error
	err = s.templates.pkgDiscussionBadgeTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"TotalCount": totalCount,
	}, err))
	checkTmplError(err, fmt.Sprintf("discussion-badge (%s)", pkgName))
}

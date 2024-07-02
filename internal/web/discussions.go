package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/web/util"

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
		s.handleGiscus(r)
		return
	}
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	repositoryName := mux.Vars(r)["repositoryName"]
	pkg, manifest, err := describe.DescribeInstalledPackage(r.Context(), namespace, name)
	if err != nil && !errors.IsNotFound(err) {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching installed package %v", name), "danger")
		return
	}

	s.handlePackageDiscussionPage(w, r, &packageDetailPageContext{
		repositoryName: repositoryName,
		manifestName:   manifestName,
		pkg:            pkg,
		manifest:       manifest,
	})

}

// clusterPackageDiscussion is a full page for showing various discussions, reactions, etc.
func (s *server) clusterPackageDiscussion(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.handleGiscus(r)
		return
	}
	pkgName := mux.Vars(r)["pkgName"]
	repositoryName := mux.Vars(r)["repositoryName"]
	pkg, manifest, err := describe.DescribeInstalledClusterPackage(r.Context(), pkgName)
	if err != nil && !errors.IsNotFound(err) {
		s.respondAlertAndLog(w, err,
			fmt.Sprintf("An error occurred fetching installed package %v", pkgName), "danger")
		return
	}

	s.handlePackageDiscussionPage(w, r, &packageDetailPageContext{
		repositoryName: repositoryName,
		manifestName:   pkgName,
		pkg:            pkg,
		manifest:       manifest,
	})
}

func (s *server) handleGiscus(r *http.Request) {
	githubUrl := r.FormValue("githubUrl")
	telemetry.SetUserProperty("github_url", githubUrl)
}

func (s *server) handlePackageDiscussionPage(w http.ResponseWriter, r *http.Request, d *packageDetailPageContext) {
	var idx repo.PackageIndex
	if err := s.repoClientset.ForRepoWithName(d.repositoryName).FetchPackageIndex(d.manifestName, &idx); err != nil {
		s.respondAlertAndLog(w, err, "An error occurred fetching versions of "+d.manifestName, "danger")
		return
	}

	if d.manifest == nil {
		d.manifest = &v1alpha1.PackageManifest{}
		if err := s.repoClientset.ForRepoWithName(d.repositoryName).
			FetchPackageManifest(d.manifestName, idx.LatestVersion, d.manifest); err != nil {
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching manifest of %v in version %v in repository %v",
					d.manifest, idx.LatestVersion, d.repositoryName), "danger")
			return
		}
	}

	pkgHref := util.GetPackageHref(d.pkg, d.manifest)

	err := s.templates.pkgDiscussionPageTmpl.Execute(w, s.enrichPage(r, map[string]any{
		"Giscus":             giscus.Client().Config,
		"Package":            d.pkg,
		"Status":             client.GetStatusOrPending(d.pkg),
		"Manifest":           d.manifest,
		"LatestVersion":      idx.LatestVersion,
		"UpdateAvailable":    s.isUpdateAvailableForPkg(r.Context(), d.pkg),
		"ShowDiscussionLink": true,
		"PackageHref":        pkgHref,
		"DiscussionHref":     fmt.Sprintf("%s/discussion", pkgHref),
		"AutoUpdate":         clientutils.AutoUpdateString(d.pkg, "Disabled"),
	}, nil))
	checkTmplError(err, fmt.Sprintf("package-discussion (%s)", d.manifestName))
}

func (s *server) discussionBadge(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkgName"]
	if pkgName == "" {
		pkgName = mux.Vars(r)["manifestName"]
	}

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

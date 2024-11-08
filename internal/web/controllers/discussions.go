package controllers

import (
	"fmt"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/giscus"
	"github.com/glasskube/glasskube/internal/httperror"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/internal/web/types"
	"github.com/glasskube/glasskube/internal/web/util"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/describe"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"os"
)

func GetPackageDiscussion(w http.ResponseWriter, r *http.Request) {
	req := getPackageContext(r).request
	pkg, manifest, err := describe.DescribeInstalledPackage(r.Context(), req.namespace, req.name)
	if err != nil && !errors.IsNotFound(err) {
		responder.SendToast(w, toast.WithErr(
			fmt.Errorf("failed to fetch installed package %v/%v: %w", req.namespace, req.name, err)))
		return
	}

	handlePackageDiscussionPage(w, r, &packageContext{
		request:  req,
		pkg:      pkg,
		manifest: manifest,
	})
}

func PostGiscus(w http.ResponseWriter, r *http.Request) {
	githubUrl := r.FormValue("githubUrl")
	telemetry.SetUserProperty("github_url", githubUrl)
}

func GetClusterPackageDiscussion(w http.ResponseWriter, r *http.Request) {
	req := getPackageContext(r).request
	pkg, manifest, err := describe.DescribeInstalledClusterPackage(r.Context(), req.manifestName)
	if err != nil && !errors.IsNotFound(err) {
		responder.SendToast(w,
			toast.WithErr(fmt.Errorf("failed to fetch installed clusterpackage %v: %w", req.manifestName, err)))
		return
	}

	handlePackageDiscussionPage(w, r, &packageContext{
		request:  req,
		pkg:      pkg,
		manifest: manifest,
	})
}

type discussionPageData struct {
	types.TemplateContextHolder
	packageDetailCommonData
	Giscus         *giscus.GiscusConfig
	DiscussionHref string
}

func handlePackageDiscussionPage(w http.ResponseWriter, r *http.Request, d *packageContext) {
	ctx := r.Context()
	repoClientset := clicontext.RepoClientsetFromContext(ctx)
	repoClient := repoClientset.ForRepoWithName(d.request.repositoryName)
	var idx repo.PackageIndex
	if err := repoClient.FetchPackageIndex(d.request.manifestName, &idx); err != nil {
		responder.SendToast(w,
			toast.WithErr(fmt.Errorf("failed to fetch package index of %v in repo %v: %w",
				d.request.manifestName, d.request.repositoryName, err)))
		return
	}

	if d.manifest == nil {
		d.manifest = &v1alpha1.PackageManifest{}
		if err := repoClientset.ForRepoWithName(d.request.repositoryName).
			FetchPackageManifest(d.request.manifestName, idx.LatestVersion, d.manifest); err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch manifest of %v (%v) in repo %v: %w",
				d.request.manifestName, idx.LatestVersion, d.request.repositoryName, err)))
			return
		}
	}

	pkgHref := util.GetPackageHrefWithFallback(d.pkg, d.manifest)

	autoUpdaterInstalled, err := clientutils.IsAutoUpdaterInstalled(r.Context())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to check whether auto updater is installed: %v\n", err)
	}

	responder.SendPage(w, r, "pages/discussion", responder.ContextualizedTemplate(&discussionPageData{
		packageDetailCommonData: packageDetailCommonData{
			Package:              d.pkg,
			Status:               client.GetStatusOrPending(d.pkg),
			Manifest:             d.manifest,
			PackageManifestUrl:   "", // TODO fix this
			LatestVersion:        idx.LatestVersion,
			UpdateAvailable:      isUpdateAvailableForPkg(ctx, d.pkg),
			ShowDiscussionLink:   true,
			PackageHref:          pkgHref,
			AutoUpdaterInstalled: autoUpdaterInstalled,
		},
		Giscus:         giscus.Client().Config,
		DiscussionHref: fmt.Sprintf("%s/discussion", pkgHref),
	}))
}

type discussionBadgeData struct {
	TotalCount int
}

func GetDiscussionBadge(w http.ResponseWriter, r *http.Request) {
	pkgName := r.PathValue("manifestName")

	var totalCount int
	if counts, err := giscus.Client().GetCountsFor(pkgName); err != nil {
		if !httperror.IsNotFound(err) {
			fmt.Fprintf(os.Stderr, "failed to get discussion counts from giscus: %v\n", err)
		}
	} else {
		totalCount = counts.ReactionCount + counts.TotalCommentCount + counts.TotalReplyCount
	}

	responder.SendComponent(w, r, "components/discussion-badge",
		responder.RawTemplate(discussionBadgeData{TotalCount: totalCount}))
}

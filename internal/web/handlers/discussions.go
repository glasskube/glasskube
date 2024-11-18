package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/glasskube/glasskube/internal/giscus"
	"github.com/glasskube/glasskube/internal/httperror"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/internal/web/types"
	"github.com/glasskube/glasskube/internal/web/util"
	"github.com/glasskube/glasskube/pkg/describe"
	"k8s.io/apimachinery/pkg/api/errors"
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
	pkgDetailCommonData, _, _, _ := resolvePkgDetailCommon(w, r.Context(), d)
	if pkgDetailCommonData == nil {
		return
	}
	pkgHref := util.GetPackageHrefWithFallback(d.pkg, d.manifest)
	responder.SendPage(w, r, "pages/discussion", responder.ContextualizedTemplate(&discussionPageData{
		packageDetailCommonData: *pkgDetailCommonData,
		Giscus:                  giscus.Client().Config,
		DiscussionHref:          fmt.Sprintf("%s/discussion", pkgHref),
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

package plain

import (
	"net/http"
	"net/url"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
)

func (r *Adapter) newManifestRequest(pi *packagesv1alpha1.PackageInfo, urlOrPath string) (*http.Request, error) {
	if parsedUrl, err := url.Parse(urlOrPath); err != nil {
		return nil, err
	} else if parsedUrl.Scheme == "" && parsedUrl.Host == "" {
		if parsedBase, err := url.Parse(pi.Status.ResolvedUrl); err != nil {
			return nil, err
		} else if ref, err := parsedBase.Parse(urlOrPath); err != nil {
			return nil, err
		} else if request, err := clientutils.NewResourcesRequest(ref.String()); err != nil {
			return nil, err
		} else {
			r.repo.ForRepoWithName(pi.Spec.RepositoryName).Authenticate(request)
			return request, nil
		}
	} else {
		return clientutils.NewResourcesRequest(urlOrPath)
	}
}

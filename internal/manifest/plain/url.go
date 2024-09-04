package plain

import (
	"net/url"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
)

func getActualManifestUrl(pi *packagesv1alpha1.PackageInfo, urlOrPath string) (string, error) {
	if parsedUrl, err := url.Parse(urlOrPath); err != nil {
		return "", err
	} else if parsedUrl.Scheme == "" && parsedUrl.Host == "" {
		if parsedBase, err := url.Parse(pi.Status.ResolvedUrl); err != nil {
			return "", err
		} else if ref, err := parsedBase.Parse(urlOrPath); err != nil {
			return "", err
		} else {
			return ref.String(), nil
		}
	} else {
		return urlOrPath, nil
	}
}

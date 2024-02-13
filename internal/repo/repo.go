package repo

import (
	"errors"
	"net/http"
	"net/url"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var defaultRepositoryURL = "https://packages.dl.glasskube.dev/packages/"
var ErrNotFound = errors.New("not found")

func UpdatePackageManifest(pi *packagesv1alpha1.PackageInfo) (err error) {
	var manifest packagesv1alpha1.PackageManifest
	var version string
	if pi.Spec.Version != "" {
		// PackageInfo has explicit version in Spec
		version = pi.Spec.Version
		if err = FetchPackageManifest(pi.Spec.RepositoryUrl, pi.Spec.Name, version, &manifest); err != nil {
			return
		}
	} else {
		version, err = FetchLatestPackageManifest(pi.Spec.RepositoryUrl, pi.Spec.Name, &manifest)
		if err != nil {
			return
		}
	}

	pi.Status.Manifest = &manifest
	pi.Status.Version = version
	return nil
}

func FetchLatestPackageManifest(repoURL, name string, target *packagesv1alpha1.PackageManifest) (version string, err error) {
	var versions PackageIndex
	if err = FetchPackageIndex(repoURL, name, &versions); err != nil {
		if err != ErrNotFound {
			return
		}
		// no versions.yaml file for package in repo. Try versionless manifest
		version = ""
	} else {
		version = versions.LatestVersion
	}
	err = FetchPackageManifest(repoURL, name, version, target)
	return
}

func FetchPackageManifest(repoURL, name, version string, target *packagesv1alpha1.PackageManifest) error {
	if url, err := getPackageManifestURL(repoURL, name, version); err != nil {
		return err
	} else {
		return fetchYAMLOrJSON(url, target)
	}
}

func FetchPackageIndex(repoURL, name string, target *PackageIndex) error {
	if url, err := getPackageIndexURL(repoURL, name); err != nil {
		return err
	} else {
		return fetchYAMLOrJSON(url, target)
	}
}

func FetchPackageRepoIndex(repoURL string, target *PackageRepoIndex) error {
	if url, err := getPackageRepoIndexURL(repoURL); err != nil {
		return err
	} else {
		return fetchYAMLOrJSON(url, target)
	}
}

func fetchYAMLOrJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch " + url + ": " + resp.Status)
	}
	return yaml.NewYAMLOrJSONDecoder(resp.Body, 4096).Decode(target)
}

func getPackageRepoIndexURL(repoURL string) (string, error) {
	return url.JoinPath(getBaseURL(repoURL), "index.yaml")
}

func getPackageIndexURL(repoURL, name string) (string, error) {
	return url.JoinPath(getBaseURL(repoURL), name, "versions.yaml")
}

func getPackageManifestURL(repoURL, name, version string) (string, error) {
	pathSegments := []string{name}
	if version != "" {
		pathSegments = append(pathSegments, version)
	}
	pathSegments = append(pathSegments, "package.yaml")
	return url.JoinPath(getBaseURL(repoURL), pathSegments...)
}

func getBaseURL(explicitRepositoryURL string) string {
	if len(explicitRepositoryURL) > 0 {
		return explicitRepositoryURL
	} else {
		return defaultRepositoryURL
	}
}

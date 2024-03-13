package client

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/httperror"
	"github.com/glasskube/glasskube/internal/repo/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type defaultClient struct {
	defaultRepositoryURL string
	idxMutex             sync.Mutex
	idxUpdate            time.Time
	packageRepoIndex     types.PackageRepoIndex
}

func New(repoURL string) RepoClient {
	return &defaultClient{defaultRepositoryURL: repoURL}
}

// FetchLatestPackageManifest implements repo.RepoClient.
func (c *defaultClient) FetchLatestPackageManifest(repoURL string, name string, target *v1alpha1.PackageManifest) (version string, err error) {
	var versions types.PackageIndex
	if err = c.FetchPackageIndex(repoURL, name, &versions); err != nil {
		return
	} else {
		version = versions.LatestVersion
	}
	err = c.FetchPackageManifest(repoURL, name, version, target)
	return
}

// FetchPackageManifest implements repo.RepoClient.
func (c *defaultClient) FetchPackageManifest(repoURL string, name string, version string, target *v1alpha1.PackageManifest) error {
	if url, err := c.GetPackageManifestURL(repoURL, name, version); err != nil {
		return err
	} else {
		return c.fetchYAMLOrJSON(url, target)
	}
}

// FetchPackageIndex implements repo.RepoClient.
func (c *defaultClient) FetchPackageIndex(repoURL string, name string, target *types.PackageIndex) error {
	if url, err := c.getPackageIndexURL(repoURL, name); err != nil {
		return err
	} else {
		return c.fetchYAMLOrJSON(url, target)
	}
}

// FetchPackageRepoIndex implements repo.RepoClient.
func (c *defaultClient) FetchPackageRepoIndex(repoURL string, target *types.PackageRepoIndex) error {
	if url, err := c.getPackageRepoIndexURL(repoURL); err != nil {
		return err
	} else {
		return c.fetchYAMLOrJSON(url, target)
	}
}

// GetLatestVersion implements repo.RepoClient.
func (c *defaultClient) GetLatestVersion(repoURL string, pkgName string) (string, error) {
	c.idxMutex.Lock()
	defer c.idxMutex.Unlock()
	if len(c.packageRepoIndex.Packages) == 0 || c.idxUpdate.Add(5*time.Minute).Before(time.Now()) {
		if err := c.FetchPackageRepoIndex(repoURL, &c.packageRepoIndex); err != nil {
			return "", err
		}
		c.idxUpdate = time.Now()
	}
	for _, pkg := range c.packageRepoIndex.Packages {
		if pkg.Name == pkgName {
			return pkg.LatestVersion, nil
		}
	}
	return "", nil
}

func (c *defaultClient) fetchYAMLOrJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if err = httperror.CheckResponse(resp); err != nil {
		return fmt.Errorf("failed to fetch %v: %w", url, err)
	}
	return yaml.NewYAMLOrJSONDecoder(resp.Body, 4096).Decode(target)
}

func (c *defaultClient) getPackageRepoIndexURL(repoURL string) (string, error) {
	return url.JoinPath(c.getBaseURL(repoURL), "index.yaml")
}

func (c *defaultClient) getPackageIndexURL(repoURL, name string) (string, error) {
	return url.JoinPath(c.getBaseURL(repoURL), url.PathEscape(name), "versions.yaml")
}

// GetPackageManifestURL implements repo.RepoClient.
func (c *defaultClient) GetPackageManifestURL(repoURL, name, version string) (string, error) {
	pathSegments := []string{url.PathEscape(name), url.PathEscape(version), "package.yaml"}
	return url.JoinPath(c.getBaseURL(repoURL), pathSegments...)
}

func (c *defaultClient) getBaseURL(explicitRepositoryURL string) string {
	if len(explicitRepositoryURL) > 0 {
		return explicitRepositoryURL
	} else {
		return c.defaultRepositoryURL
	}
}

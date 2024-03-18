package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/httperror"
	"github.com/glasskube/glasskube/internal/repo/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type defaultClient struct {
	defaultRepositoryURL string
	maxCacheAge          time.Duration
	cache                sync.Map
	debug                bool
}

type cacheItem struct {
	bytes   []byte
	updated time.Time
	mutex   sync.Mutex
}

func New(repoURL string, maxCacheAge time.Duration) *defaultClient {
	return &defaultClient{defaultRepositoryURL: repoURL, maxCacheAge: maxCacheAge}
}

func NewDebug(repoURL string, maxCacheAge time.Duration) *defaultClient {
	c := New(repoURL, maxCacheAge)
	c.debug = true
	return c
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
	var idx types.PackageRepoIndex
	if err := c.FetchPackageRepoIndex(repoURL, &idx); err != nil {
		return "", err
	}
	for _, pkg := range idx.Packages {
		if pkg.Name == pkgName {
			return pkg.LatestVersion, nil
		}
	}
	return "", nil
}

func (c *defaultClient) fetchYAMLOrJSON(url string, target any) error {
	cached := &cacheItem{}
	if c, hit := c.cache.LoadOrStore(url, cached); hit {
		if c, ok := c.(*cacheItem); ok {
			cached = c
		} else {
			return errors.New("unexpected cache type")
		}
	}

	if cached.updated.Add(c.maxCacheAge).After(time.Now()) {
		if c.debug {
			fmt.Fprintln(os.Stderr, "cache hit", url)
		}
		return yaml.Unmarshal(cached.bytes, target)
	}

	cached.mutex.Lock()
	defer cached.mutex.Unlock()

	// try again after acquiring the mutex
	if cached.updated.Add(c.maxCacheAge).After(time.Now()) {
		if c.debug {
			fmt.Fprintln(os.Stderr, "cache hit (after lock)", url)
		}
		return yaml.Unmarshal(cached.bytes, target)
	}

	if c.debug {
		fmt.Fprintln(os.Stderr, "cache miss", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if err = httperror.CheckResponse(resp); err != nil {
		return fmt.Errorf("failed to fetch %v: %w", url, err)
	}
	if bytes, err := io.ReadAll(resp.Body); err != nil {
		return err
	} else if err := yaml.Unmarshal(bytes, target); err != nil {
		return err
	} else {
		cached.bytes = bytes
		cached.updated = time.Now()
		return nil
	}
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

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
	"github.com/glasskube/glasskube/internal/contenttype"
	"github.com/glasskube/glasskube/internal/httperror"
	"github.com/glasskube/glasskube/internal/repo/client/auth"
	"github.com/glasskube/glasskube/internal/repo/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type defaultClient struct {
	auth.Authenticator
	url         string
	maxCacheAge time.Duration
	cache       sync.Map
	debug       bool
}

type cacheItem struct {
	bytes   []byte
	updated time.Time
	mutex   sync.Mutex
}

func New(url string, authenticator auth.Authenticator, maxCacheAge time.Duration) *defaultClient {
	return &defaultClient{url: url, Authenticator: authenticator, maxCacheAge: maxCacheAge}
}

func NewDebug(url string, authenticator auth.Authenticator, maxCacheAge time.Duration) *defaultClient {
	c := New(url, authenticator, maxCacheAge)
	c.debug = true
	return c
}

var _ RepoClient = &defaultClient{}

// FetchLatestPackageManifest implements repo.RepoClient.
func (c *defaultClient) FetchLatestPackageManifest(name string, target *v1alpha1.PackageManifest) (
	version string, err error,
) {
	var versions types.PackageIndex
	if err = c.FetchPackageIndex(name, &versions); err != nil {
		return
	} else {
		version = versions.LatestVersion
	}
	err = c.FetchPackageManifest(name, version, target)
	return
}

// FetchPackageManifest implements repo.RepoClient.
func (c *defaultClient) FetchPackageManifest(name string, version string,
	target *v1alpha1.PackageManifest) error {
	if url, err := c.GetPackageManifestURL(name, version); err != nil {
		return err
	} else {
		return c.fetchYAMLOrJSON(url, target)
	}
}

// FetchPackageIndex implements repo.RepoClient.
func (c *defaultClient) FetchPackageIndex(name string, target *types.PackageIndex) error {
	if url, err := c.getPackageIndexURL(name); err != nil {
		return err
	} else {
		return c.fetchYAMLOrJSON(url, target)
	}
}

// FetchPackageRepoIndex implements repo.RepoClient.
func (c *defaultClient) FetchPackageRepoIndex(target *types.PackageRepoIndex) error {
	if url, err := c.getPackageRepoIndexURL(); err != nil {
		return err
	} else {
		return c.fetchYAMLOrJSON(url, target)
	}
}

// GetLatestVersion implements repo.RepoClient.
func (c *defaultClient) GetLatestVersion(pkgName string) (string, error) {
	var idx types.PackageRepoIndex
	if err := c.FetchPackageRepoIndex(&idx); err != nil {
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

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	c.Authenticate(request)
	request.Header.Add("Accept", contenttype.MediaTypeJSON)
	request.Header.Add("Accept", contenttype.MediaTypeYAML)
	resp, err := httperror.CheckResponse(http.DefaultClient.Do(request))
	if err != nil {
		return fmt.Errorf("failed to fetch %v: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if err := contenttype.IsJsonOrYaml(resp); err != nil {
		return fmt.Errorf("could not decode %v: %w", url, err)
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

func (c *defaultClient) getPackageRepoIndexURL() (string, error) {
	return url.JoinPath(c.getBaseURL(), "index.yaml")
}

func (c *defaultClient) getPackageIndexURL(name string) (string, error) {
	return url.JoinPath(c.getBaseURL(), url.PathEscape(name), "versions.yaml")
}

// GetPackageManifestURL implements repo.RepoClient.
func (c *defaultClient) GetPackageManifestURL(name, version string) (string, error) {
	pathSegments := []string{url.PathEscape(name), url.PathEscape(version), "package.yaml"}
	return url.JoinPath(c.getBaseURL(), pathSegments...)
}

func (c *defaultClient) getBaseURL() string {
	return c.url
}

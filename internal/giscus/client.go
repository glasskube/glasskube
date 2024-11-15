package giscus

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/glasskube/glasskube/internal/httperror"
)

var client *GiscusClient

type GiscusConfig struct {
	Repo       string
	RepoId     string
	Category   string
	CategoryId string
}

type GiscusClient struct {
	Config      *GiscusConfig
	cache       sync.Map
	maxCacheAge time.Duration
}

type cacheItem struct {
	data    GiscusDiscussionCounts
	updated time.Time
	mutex   sync.Mutex
}

type GiscusDiscussionCounts struct {
	TotalCommentCount int `json:"totalCommentCount"`
	TotalReplyCount   int `json:"totalReplyCount"`
	ReactionCount     int `json:"reactionCount"`
}

type giscusResponse struct {
	Discussion GiscusDiscussionCounts `json:"discussion"`
}

func Client() *GiscusClient {
	if client == nil {
		client = &GiscusClient{
			Config: &GiscusConfig{
				Repo:       "glasskube/glasskube",
				RepoId:     "R_kgDOLDumDw",
				Category:   "Packages",
				CategoryId: "DIC_kwDOLDumD84CcybS",
			},
			maxCacheAge: 5 * time.Minute,
		}
	}
	return client
}

func (c *GiscusClient) GetCountsFor(pkgName string) (*GiscusDiscussionCounts, error) {
	cached := &cacheItem{}
	if c, ok := c.cache.LoadOrStore(pkgName, cached); ok {
		if c, ok := c.(*cacheItem); ok {
			cached = c
		} else {
			return nil, errors.New("unexpected cache type")
		}
	}

	if cached.updated.Add(c.maxCacheAge).After(time.Now()) {
		return &cached.data, nil
	}

	cached.mutex.Lock()
	defer cached.mutex.Unlock()

	// try again after acquiring the mutex
	if cached.updated.Add(c.maxCacheAge).After(time.Now()) {
		return &cached.data, nil
	}

	if u, err := url.Parse("https://giscus.app/api/discussions"); err != nil {
		return nil, err
	} else {
		values := u.Query()
		values.Add("repo", c.Config.Repo)
		values.Add("term", pkgName)
		values.Add("category", c.Config.Category)
		values.Add("strict", "false")
		values.Add("first", "1")
		values.Add("number", "0")
		u.RawQuery = values.Encode()

		resp, err := httperror.CheckResponse(http.Get(u.String()))
		if err != nil {
			return nil, err
		}
		var target giscusResponse
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := io.ReadAll(resp.Body); err != nil {
			return nil, err
		} else if err := json.Unmarshal(bytes, &target); err != nil {
			return nil, err
		} else {
			cached.data = target.Discussion
			cached.updated = time.Now()
			return &target.Discussion, nil
		}
	}
}

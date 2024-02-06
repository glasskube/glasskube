package repo

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

var (
	defaultRepositoryUrl = "https://packages.dl.glasskube.dev/packages/"
)

type PackageRepoIndex struct {
	Packages []PackageTeaser
}

type PackageTeaser struct {
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription,omitempty"`
	IconUrl          string `json:"iconUrl,omitempty"`
}

func FetchPackageManifest(ctx context.Context, pi *packagesv1alpha1.PackageInfo) error {
	log := log.FromContext(ctx)
	url, err := getPackageManifestUrl(*pi)
	if err != nil {
		log.Error(err, "can not get manifest url")
		return err
	}
	log.Info("starting to fetch " + url)
	body, err := doFetch(url)
	if err != nil {
		return err
	}
	var manifest packagesv1alpha1.PackageManifest
	if err = yaml.Unmarshal(body, &manifest); err != nil {
		return err
	}
	pi.Status.Manifest = &manifest
	return nil
}

func FetchPackageRepoIndex(repoUrl string) (*PackageRepoIndex, error) {
	if len(repoUrl) == 0 {
		repoUrl = defaultRepositoryUrl
	}
	url, err := url.JoinPath(repoUrl, "index.yaml")
	if err != nil {
		return nil, err
	}
	body, err := doFetch(url)
	if err != nil {
		return nil, err
	}
	var index PackageRepoIndex
	if err = yaml.Unmarshal(body, &index); err != nil {
		return nil, err
	}
	return &index, nil
}

func doFetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch " + url + " : " + resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func getPackageManifestUrl(pi packagesv1alpha1.PackageInfo) (string, error) {
	var baseUrl string
	if len(pi.Spec.RepositoryUrl) > 0 {
		baseUrl = pi.Spec.RepositoryUrl
	} else {
		baseUrl = defaultRepositoryUrl
	}
	return url.JoinPath(baseUrl, pi.Spec.Name, "package.yaml")
}

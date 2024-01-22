package repo

import (
	"context"
	"errors"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"io"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

var (
	defaultRepositoryUrl = "https://packages.dl.glasskube.dev/packages/"
)

func FetchPackageManifest(ctx context.Context, pi *packagesv1alpha1.PackageInfo) error {
	log := log.FromContext(ctx)
	url, err := getPackageManifestUrl(*pi)
	if err != nil {
		log.Error(err, "can not get manifest url")
		return err
	}
	log.Info("starting to fetch manifest from " + url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("could not fetch package manifest: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var manifest packagesv1alpha1.PackageManifest
	// TODO: Figure out why Helm.Values is not unmarshalled
	if err := yaml.Unmarshal(body, &manifest); err != nil {
		return err
	}
	pi.Status.Manifest = &manifest
	return nil
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

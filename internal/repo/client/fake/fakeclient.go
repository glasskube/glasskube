package fake

import (
	"errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/repo/client/auth"
	"github.com/glasskube/glasskube/internal/repo/types"
)

func ClientsetWithClient(client *fakeClient) *fakeClientset {
	return &fakeClientset{Client: client}
}

func EmptyClientset() *fakeClientset {
	return ClientsetWithClient(EmptyClient())
}

type fakeClientset struct {
	Client *fakeClient
}

// Default implements client.RepoClientset.
func (f *fakeClientset) Default() client.RepoClient {
	return f.Client
}

// Meta implements client.RepoClientset.
func (f *fakeClientset) Meta() client.RepoMetaclient {
	return f.Client
}

// ForRepo implements client.RepoClientset.
func (f *fakeClientset) ForRepo(repo v1alpha1.PackageRepository) client.RepoClient {
	return f.Client
}

// ForRepoWithName implements client.RepoClientset.
func (f *fakeClientset) ForRepoWithName(name string) client.RepoClient {
	return f.Client
}

// ForPackage implements client.RepoClientset.
func (f *fakeClientset) ForPackage(pkg ctrlpkg.Package) client.RepoClient {
	return f.Client
}

var _ client.RepoClientset = &fakeClientset{}

// fakeClient is a mock implementation of RepoClient for use in tests
type fakeClient struct {
	auth.Authenticator
	Packages            map[string]map[string]*v1alpha1.PackageManifest
	PackageRepositories []v1alpha1.PackageRepository
}

func EmptyClient() *fakeClient {
	return &fakeClient{
		Authenticator:       auth.Noop(),
		PackageRepositories: []v1alpha1.PackageRepository{{}},
	}
}

func EmptyClientWithAuth(auth auth.Authenticator) *fakeClient {
	return &fakeClient{
		Authenticator:       auth,
		PackageRepositories: []v1alpha1.PackageRepository{{}},
	}
}

// FetchMetaIndex implements client.RepoMetaclient.
func (f *fakeClient) FetchMetaIndex(target *types.MetaIndex) error {
	panic("unimplemented")
}

// GetReposForPackage implements client.RepoAggregator.
func (f *fakeClient) GetReposForPackage(name string) ([]v1alpha1.PackageRepository, error) {
	return f.PackageRepositories, nil
}

func (f *fakeClient) AddPackage(name, version string, manifest *v1alpha1.PackageManifest) {
	if f.Packages == nil {
		f.Clear()
	}
	if _, ok := f.Packages[name]; !ok {
		f.Packages[name] = map[string]*v1alpha1.PackageManifest{}
	}
	f.Packages[name][version] = manifest
}

func (f *fakeClient) Clear() {
	f.Packages = map[string]map[string]*v1alpha1.PackageManifest{}
}

var _ client.RepoClient = &fakeClient{}

// FetchLatestPackageManifest implements client.RepoClient.
func (f *fakeClient) FetchLatestPackageManifest(name string, target *v1alpha1.PackageManifest) (
	version string, err error,
) {
	if versions, ok := f.Packages[name]; ok {
		for v, m := range versions {
			*target = *m
			return v, nil
		}
	}
	return "", errors.New("not found")
}

// FetchPackageIndex implements client.RepoClient.
func (f *fakeClient) FetchPackageIndex(name string, target *types.PackageIndex) error {
	if versions, ok := f.Packages[name]; ok {
		var result types.PackageIndex
		for v := range versions {
			result.LatestVersion = v
			result.Versions = append(result.Versions, types.PackageIndexItem{Version: v})
		}
		*target = result
		return nil
	}
	return errors.New("not found")
}

// FetchPackageManifest implements client.RepoClient.
func (f *fakeClient) FetchPackageManifest(name string, version string, target *v1alpha1.PackageManifest) error {
	if versions, ok := f.Packages[name]; ok {
		if manifest, ok := versions[version]; ok {
			*target = *manifest
			return nil
		}
	}
	return errors.New("not found")
}

// FetchPackageRepoIndex implements client.RepoClient.
func (f *fakeClient) FetchPackageRepoIndex(target *types.PackageRepoIndex) error {
	var result types.PackageRepoIndex
	for pkg, versions := range f.Packages {
		item := types.PackageRepoIndexItem{Name: pkg}
		for v := range versions {
			item.LatestVersion = v
		}
		result.Packages = append(result.Packages, item)
	}
	*target = result
	return nil
}

// GetLatestVersion implements client.RepoClient.
func (f *fakeClient) GetLatestVersion(pkgName string) (string, error) {
	if versions, ok := f.Packages[pkgName]; ok {
		for v := range versions {
			return v, nil
		}
	}
	return "", errors.New("not found")
}

// GetPackageManifestURL implements client.RepoClient.
func (f *fakeClient) GetPackageManifestURL(name string, version string) (string, error) {
	return "fake url", nil
}

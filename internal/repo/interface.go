package repo

import (
	"time"

	"github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/repo/types"
)

type (
	PackageIndex         = types.PackageIndex
	PackageIndexItem     = types.PackageIndexItem
	PackageRepoIndex     = types.PackageRepoIndex
	PackageRepoIndexItem = types.PackageRepoIndexItem
)

var (
	DefaultClient              = client.New("https://packages.dl.glasskube.dev/packages/", 5*time.Minute)
	FetchLatestPackageManifest = DefaultClient.FetchLatestPackageManifest
	FetchPackageManifest       = DefaultClient.FetchPackageManifest
	FetchPackageIndex          = DefaultClient.FetchPackageIndex
	FetchPackageRepoIndex      = DefaultClient.FetchPackageRepoIndex
	GetLatestVersion           = DefaultClient.GetLatestVersion
	GetPackageManifestURL      = DefaultClient.GetPackageManifestURL
)

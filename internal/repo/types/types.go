package types

import "github.com/glasskube/glasskube/api/v1alpha1"

type PackageIndex struct {
	Versions      []PackageIndexItem `json:"versions" jsonschema:"required"`
	LatestVersion string             `json:"latestVersion" jsonschema:"required"`
}

type PackageIndexItem struct {
	Version string `json:"version" jsonschema:"required"`
}

type PackageRepoIndex struct {
	Packages []PackageRepoIndexItem `json:"packages" jsonschema:"required"`
}

type PackageRepoIndexItem struct {
	Name             string                 `json:"name"`
	ShortDescription string                 `json:"shortDescription,omitempty"`
	IconUrl          string                 `json:"iconUrl,omitempty"`
	LatestVersion    string                 `json:"latestVersion,omitempty"`
	Scope            *v1alpha1.PackageScope `json:"scope,omitempty"`
}

type MetaIndex struct {
	Packages []MetaIndexItem
}

type MetaIndexItem struct {
	PackageRepoIndexItem
	Repos []string `json:"repos,omitempty"`
}

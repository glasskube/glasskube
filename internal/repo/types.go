package repo

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
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription,omitempty"`
	IconUrl          string `json:"iconUrl,omitempty"`
}

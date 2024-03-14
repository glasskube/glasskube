package cliutils

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/releaseinfo"
)

// CheckForUpdate determines the new version by fetching the latest release info.
// If no newer version exists, it returns nil.
func CheckForUpdate() (*string, error) {
	if releaseInfo, err := releaseinfo.FetchLatestRelease(); err != nil {
		return nil, err
	} else if version, err := semver.NewVersion(config.Version); err != nil {
		return nil, err
	} else if latestVersion, err := semver.NewVersion(releaseInfo.Version); err != nil {
		return nil, err
	} else if latestVersion.GreaterThan(version) {
		return &releaseInfo.Version, nil
	} else {
		return nil, nil
	}
}

func UpdateFetch() {
	newerVersion, err := CheckForUpdate()
	if err == nil && newerVersion != nil {
		fmt.Fprintf(os.Stderr, "📣 A newer version of Glasskube is available: %s → %s\n", config.Version, *newerVersion)
		fmt.Fprintf(os.Stderr, "📘 Release notes: https://github.com/glasskube/glasskube/releases/tag/v%v\n\n", *newerVersion)
	}
}

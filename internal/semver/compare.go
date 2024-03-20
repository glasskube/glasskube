package semver

import (
	"strconv"

	"github.com/Masterminds/semver/v3"
)

// isUpgradable checks if latest is greater than installed, according to semver.
// As a fallback if either cannot be parsed as semver, it returns whether they are different.
// Am important deviation from the semver standard is that this function DOES try to interpret
// the version metadata as a number when comparing.
func IsUpgradable(installed, latest string) bool {
	var parsedInstalled, parsedLatest *semver.Version
	var err error
	if parsedInstalled, err = semver.NewVersion(installed); err != nil {
		parsedInstalled = nil
	} else if parsedLatest, err = semver.NewVersion(latest); err != nil {
		parsedLatest = nil
	}
	if parsedLatest != nil && parsedInstalled != nil {
		if parsedLatest.GreaterThan(parsedInstalled) {
			return true
		}
		return isUpgradableMetadata(parsedInstalled, parsedLatest)
	} else {
		return installed != latest
	}
}

func isUpgradableMetadata(installed, latest *semver.Version) bool {
	latestMeta := latest.Metadata()
	installedMeta := installed.Metadata()

	if latestMeta == installedMeta {
		return false
	}

	if latestMeta == "" {
		return false
	} else if installedMeta == "" {
		return true
	}

	if latestMetaInt, err := strconv.Atoi(latestMeta); err != nil {
		return installedMeta != latestMeta
	} else if installedMetaInt, err := strconv.Atoi(installedMeta); err != nil {
		return installedMeta != latestMeta
	} else {
		return installedMetaInt < latestMetaInt
	}
}

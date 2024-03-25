package semver

import (
	"strconv"

	"github.com/Masterminds/semver/v3"
)

// IsUpgradable checks if desired is greater than installed, according to semver.
// As a fallback if either cannot be parsed as semver, it returns whether they are different.
// Am important deviation from the semver standard is that this function DOES try to interpret
// the version metadata as a number when comparing.
func IsUpgradable(installed, desired string) bool {
	var parsedInstalled, parsedDesired *semver.Version
	var err error
	if parsedInstalled, err = semver.NewVersion(installed); err != nil {
		parsedInstalled = nil
	} else if parsedDesired, err = semver.NewVersion(desired); err != nil {
		parsedDesired = nil
	}
	if parsedDesired != nil && parsedInstalled != nil {
		if parsedDesired.GreaterThan(parsedInstalled) {
			return true
		} else if parsedDesired.Equal(parsedInstalled) {
			return isUpgradableMetadata(parsedInstalled, parsedDesired)
		} else {
			return false
		}
	} else {
		return installed != desired
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

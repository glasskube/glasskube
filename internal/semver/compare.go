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
	if parsedInstalled, err := semver.NewVersion(installed); err != nil {
		return installed != desired
	} else if parsedDesired, err := semver.NewVersion(desired); err != nil {
		return installed != desired
	} else {
		return IsVersionUpgradable(parsedInstalled, parsedDesired)
	}
}

func IsVersionUpgradable(installed, desired *semver.Version) bool {
	return desired.GreaterThan(installed) || (desired.Equal(installed) && isUpgradableMetadata(installed, desired))
}

func isUpgradableMetadata(installed, desired *semver.Version) bool {
	desiredMeta := desired.Metadata()
	installedMeta := installed.Metadata()

	if desiredMeta == installedMeta {
		return false
	}

	if desiredMeta == "" {
		return false
	} else if installedMeta == "" {
		return true
	}

	if desiredMetadataInt, err := strconv.Atoi(desiredMeta); err != nil {
		return installedMeta != desiredMeta
	} else if installedMetaInt, err := strconv.Atoi(installedMeta); err != nil {
		return installedMeta != desiredMeta
	} else {
		return installedMetaInt < desiredMetadataInt
	}
}

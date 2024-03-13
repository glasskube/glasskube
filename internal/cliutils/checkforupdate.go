package cliutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/glasskube/glasskube/internal/config"
)

type ReleaseInfo struct {
	Version string `json:"version"`
}

func versionIsGreater(version1, version2 string) bool {
	r := regexp.MustCompile(`(\d+)`)
	v1 := r.FindAllString(version1, -1)
	v2 := r.FindAllString(version2, -1)

	for i := 0; i < len(v1) && i < len(v2); i++ {
		num1, _ := strconv.Atoi(v1[i])
		num2, _ := strconv.Atoi(v2[i])

		if num1 < num2 {
			return true
		} else if num1 > num2 {
			return false
		}
	}

	return len(v1) < len(v2)
}

func CheckForUpdate() (bool, string) {
	url := "https://glasskube.dev/release.json"

	resp, _ := http.Get(url)

	defer func() {
		_ = resp.Body.Close()

	}()

	var releaseInfo ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return false, ""
	}

	if versionIsGreater(config.Version, releaseInfo.Version) {
		return true, releaseInfo.Version
	}

	return false, ""
}

func UpdateFetch() {
	updateAvailable, latestVersion := CheckForUpdate()
	if updateAvailable {
		fmt.Fprintf(os.Stderr, "📣 A newer version of Glasskube is available: %s → %s\n", config.Version, latestVersion)
		fmt.Fprintf(os.Stderr, "📘 Release notes: https://github.com/glasskube/glasskube/releases/tag/v%v\n\n", latestVersion)
	}
}

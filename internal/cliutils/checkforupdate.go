package cliutils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/glasskube/glasskube/internal/config"
)

type ReleaseInfo struct {
	Version string `json:"version"`
}

func checkForUpdate() (bool, string) {

	url := "https://glasskube.dev/release.json"

	resp, err := http.Get(url)
	if err != nil {
		return false, ""
	}
	_ = resp.Body.Close()

	var releaseInfo ReleaseInfo
	err = json.NewDecoder(resp.Body).Decode(&releaseInfo)
	if err != nil {
		return false, ""
	}

	if releaseInfo.Version != config.Version {
		return true, releaseInfo.Version
	}

	return false, ""
}

func updateFetch() {
	updateAvailable, latestVersion := checkForUpdate()

	if updateAvailable {
		fmt.Printf("\n   --------------------------------------------------------------------------------------------------------------- \n\n")
		fmt.Printf("                                           Update available %s → %s\n", config.Version, latestVersion)
		fmt.Printf("                              Please update glasskube to the latest version\n\n")
		fmt.Printf("   --------------------------------------------------------------------------------------------------------------- \n\n")
	}
}

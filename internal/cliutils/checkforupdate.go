package cliutils

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func checkForUpdate() (bool, string, error) {
	url := "https://glasskube.dev/release.json"

	resp, err := http.Get(url)
	if err != nil {
		return false, "", fmt.Errorf("failed to fetch release information: %v", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var releaseInfo ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return false, "", fmt.Errorf("failed to decode release information: %v", err)
	}

	if versionIsGreater(config.Version, releaseInfo.Version) {
		return true, releaseInfo.Version, nil
	}

	return false, "", nil
}

func UpdateFetch() {
	updateAvailable, latestVersion, err := checkForUpdate()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printUpdateMessage := func() {
		fmt.Printf("\n   --------------------------------------------------------------------------------------------------------------- \n\n")
		fmt.Printf("                                           Update available %s → %s\n", config.Version, latestVersion)
		fmt.Printf("                              Please update glasskube to the latest version\n\n")
		fmt.Printf("   --------------------------------------------------------------------------------------------------------------- \n\n")
	}

	if updateAvailable {
		printUpdateMessage()
	}
}

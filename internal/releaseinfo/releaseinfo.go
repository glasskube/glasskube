package releaseinfo

import (
	"encoding/json"
	"net/http"
)

type ReleaseInfo struct {
	Version string `json:"version"`
}

func FetchLatestRelease() (*ReleaseInfo, error) {
	url := "https://glasskube.dev/release.json"

	resp, _ := http.Get(url)

	defer func() {
		_ = resp.Body.Close()
	}()

	var releaseInfo ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return nil, err
	}

	return &releaseInfo, nil
}

package releaseinfo

import (
	"encoding/json"
	"net/http"

	"github.com/glasskube/glasskube/internal/httperror"
)

type ReleaseInfo struct {
	Version string `json:"version"`
}

func FetchLatestRelease() (*ReleaseInfo, error) {
	url := "https://glasskube.dev/release.json"

	resp, err := httperror.CheckResponse(http.Get(url))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var releaseInfo ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return nil, err
	}

	return &releaseInfo, nil
}

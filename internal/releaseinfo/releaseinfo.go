package releaseinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/glasskube/glasskube/internal/httperror"
)

type ReleaseInfo struct {
	Version string `json:"version"`
}

func FetchLatestRelease() (*ReleaseInfo, error) {
	url := "https://glasskube.dev/release.json"

	resp, err := http.Get(url)
	if err != nil {
		if err := httperror.CheckResponse(resp); err != nil {
			switch {
			case httperror.IsNetworkError(err):
				return nil, fmt.Errorf("Network Timeout, check your network %v: %v", url, err)
			case httperror.Is(err, http.StatusServiceUnavailable):
				return nil, fmt.Errorf("Could not connect to the glasskube server. Please check your internet connection and try again. %v: %v", url, err)
			default:
				return nil, fmt.Errorf("Failed to fetch glasskube version %v: %v", url, err)
			}
		}
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

package releaseinfo

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/glasskube/glasskube/internal/httperror"
)

type ReleaseInfo struct {
	Version string `json:"version"`
}

var cachedResponse *ReleaseInfo
var mutex sync.Mutex

func FetchLatestRelease() (*ReleaseInfo, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if cachedResponse == nil {
		resp, err := httperror.CheckResponse(http.Get("https://glasskube.dev/release.json"))
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

		cachedResponse = &releaseInfo
	}

	return &ReleaseInfo{cachedResponse.Version}, nil
}

package util

import (
	"fmt"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/telemetry/annotations"
	"net/http"
	"os"
)

func IsGitopsModeEnabled(req *http.Request) bool {
	coreListers := clicontext.CoreListersFromContext(req.Context())
	if coreListers != nil && coreListers.NamespaceLister != nil {
		if ns, err := (*coreListers.NamespaceLister).Get("glasskube-system"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch glasskube-system namespace: %v\n", err)
			return true
		} else {
			return annotations.IsGitopsModeEnabled(ns.Annotations)
		}
	}
	return false
}

package util

import (
	"fmt"
	"net/http"
	"os"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/telemetry/annotations"
)

func IsGitopsModeEnabled(req *http.Request) bool {
	coreListers := clicontext.CoreListersFromContext(req.Context())
	if coreListers != nil && coreListers.NamespaceLister != nil {
		if ns, err := (*coreListers.NamespaceLister).Get("glasskube-system"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to determine GitOps mode: %v\n", err)
			return true
		} else {
			return annotations.IsGitopsModeEnabled(ns.Annotations)
		}
	}
	return false
}

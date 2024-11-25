package util

import (
	"context"
	"fmt"
	"net/http"
	"os"

	webcontext "github.com/glasskube/glasskube/internal/web/context"

	"github.com/Masterminds/semver/v3"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/web/types"
)

func GetVersionDetails(req *http.Request) types.VersionDetails {
	ctx := req.Context()
	details := types.VersionDetails{}
	operatorVersion, clientVersion, err := getGlasskubeVersions(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check for version mismatch: %v\n", err)
	} else if operatorVersion != nil && clientVersion != nil && !operatorVersion.Equal(clientVersion) {
		details.MismatchWarning = true
	}
	if operatorVersion != nil && clientVersion != nil && !config.IsDevBuild() {
		details.OperatorVersion = operatorVersion.String()
		details.ClientVersion = clientVersion.String()
		details.NeedsOperatorUpdate = operatorVersion.LessThan(clientVersion)
	}
	if config.IsDevBuild() {
		details.OperatorVersion = config.Version // TODO not correct? you could also dev against a "prod" cluster
		details.ClientVersion = config.Version
	}
	return details
}

func getGlasskubeVersions(ctx context.Context) (*semver.Version, *semver.Version, error) {
	if !config.IsDevBuild() {
		coreListers := webcontext.CoreListersFromContext(ctx)
		if operatorVersion, err := clientutils.GetPackageOperatorVersionForLister(coreListers.DeploymentLister); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check package operator version: %v\n", err)
			return nil, nil, err
		} else if parsedOperator, err := semver.NewVersion(operatorVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse operator version %v: %v\n", operatorVersion, err)
			return nil, nil, err
		} else if parsedClient, err := semver.NewVersion(config.Version); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse client version %v: %v\n", config.Version, err)
			return nil, nil, err
		} else {
			return parsedOperator, parsedClient, nil
		}
	}
	return nil, nil, nil
}

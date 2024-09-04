package cliutils

import (
	"context"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/clientutils"

	"github.com/glasskube/glasskube/internal/config"
)

func CheckPackageOperatorVersion(ctx context.Context) error {
	operatorVersion, err := clientutils.GetPackageOperatorVersion(ctx)
	if err != nil {
		return err
	}
	if config.IsDevBuild() && operatorVersion != "" {
		fmt.Fprintf(os.Stderr, "❗ Glasskube CLI version is dev but the operator version is %s\n", operatorVersion[1:])
	} else if operatorVersion[1:] != config.Version {
		fmt.Fprintf(os.Stderr, "❗ Glasskube PackageOperator needs to be updated: %s -> %s\n", operatorVersion[1:], config.Version)
		fmt.Fprintf(os.Stderr, "💡 Please run `glasskube bootstrap` again to update Glasskube PackageOperator\n")
	}
	return nil
}

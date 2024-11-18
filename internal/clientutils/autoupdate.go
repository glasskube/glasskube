package clientutils

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"k8s.io/apimachinery/pkg/api/errors"
)

func AutoUpdateString(pkg ctrlpkg.Package, disabledStr string) string {
	if !pkg.IsNil() {
		if pkg.AutoUpdatesEnabled() {
			return "Enabled"
		}
		return disabledStr
	}
	return ""
}

func IsAutoUpdaterInstalled(ctx context.Context) (bool, error) {
	client := clicontext.PackageClientFromContext(ctx)
	var pkg v1alpha1.ClusterPackage
	if err := client.ClusterPackages().Get(ctx, "glasskube-autoupdater", &pkg); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

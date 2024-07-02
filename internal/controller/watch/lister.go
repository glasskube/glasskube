package watch

import (
	"context"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
)

type PackageLister interface {
	ListPackages(ctx context.Context) ([]ctrlpkg.Package, error)
}

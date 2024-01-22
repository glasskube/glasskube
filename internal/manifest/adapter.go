package manifest

import (
	"context"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/manifest/result"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ManifestAdapter interface {
	ControllerInit(builder *builder.Builder) error
	Reconcile(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*result.ReconcileResult, error)
}

package manifest

import (
	"context"

	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/manifest/result"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ManifestAdapter interface {
	ControllerInit(builder *builder.Builder, client client.Client, scheme *runtime.Scheme) error
	Reconcile(ctx context.Context,
		pkg ctrlpkg.Package,
		manifest *packagesv1alpha1.PackageManifest,
		patches manifestvalues.TargetPatches,
	) (*result.ReconcileResult, error)
}

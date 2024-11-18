package manifest

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/manifest/result"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/resourcepatch"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ManifestAdapter interface {
	ControllerInit(
		builder *builder.Builder,
		client client.Client,
		repo repoclient.RepoClientset,
		scheme *runtime.Scheme,
	) error
	Reconcile(
		ctx context.Context,
		pkg ctrlpkg.Package,
		pi *v1alpha1.PackageInfo,
		patches resourcepatch.TargetPatches,
	) (*result.ReconcileResult, error)
}

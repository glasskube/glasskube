package manifesttransformations

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/resourcepatch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GeneratePatches(tr v1alpha1.TransformationDefinition, resolvedValue any) ([]resourcepatch.TargetPatch, error) {
	result := make([]resourcepatch.TargetPatch, len(tr.Targets))
	for i, t := range tr.Targets {
		if patch, err := resourcepatch.GenerateTargetPatch(t, resolvedValue); err != nil {
			return nil, err
		} else {
			result[i] = *patch
		}
	}
	return result, nil
}

func ResolveAndGeneratePatches(
	ctx context.Context,
	client client.Client,
	pkg ctrlpkg.Package,
	manifest *v1alpha1.PackageManifest,
) (resourcepatch.TargetPatches, error) {
	if len(manifest.Transformations) == 0 {
		return nil, nil
	}
	resolver := NewResolver(client)
	var allPatches []resourcepatch.TargetPatch
	for _, t := range manifest.Transformations {
		if value, err := resolver.Resolve(ctx, pkg, t.Source); err != nil {
			return nil, err
		} else if patches, err := GeneratePatches(t, value); err != nil {
			return nil, err
		} else {
			allPatches = append(allPatches, patches...)
		}
	}
	return allPatches, nil
}

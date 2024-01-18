package flux

import (
	"context"
	"fmt"
	"time"

	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type FluxHelmAdapter struct {
	scheme *runtime.Scheme
}

func NewAdapter(scheme *runtime.Scheme) manifest.ManifestAdapter {
	return &FluxHelmAdapter{scheme: scheme}
}

func (a FluxHelmAdapter) ControllerInit(builder *builder.Builder) {
	sourcev1beta2.AddToScheme(a.scheme)
	builder.Owns(&sourcev1beta2.HelmRepository{})
	builder.Owns(&corev1.Namespace{})
}

func (a FluxHelmAdapter) Reconcile(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*result.ReconcileResult, error) {
	if err := ensureNamespace(ctx, client, pkg, manifest); err != nil {
		return nil, fmt.Errorf("could not ensure namespace: %w", err)
	}
	if err := ensureHelmRepository(ctx, client, pkg, manifest); err != nil {
		return nil, fmt.Errorf("could not ensure helm repository: %w", err)
	}
	// TODO: Ensure HelmRelease
	// TODO: Construct Result from HelmRelease status
	return result.Waiting("TODO"), nil
}

func ensureNamespace(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) error {
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: manifest.DefaultNamespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, client, &namespace, func() error {
		return controllerutil.SetOwnerReference(pkg, &namespace, client.Scheme())
	})
	return err
}

func ensureHelmRepository(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) error {
	repository := sourcev1beta2.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifest.Name,
			Namespace: manifest.DefaultNamespace,
		},
		Spec: sourcev1beta2.HelmRepositorySpec{},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, client, &repository, func() error {
		repository.Spec.URL = manifest.Helm.RepositoryUrl
		repository.Spec.Interval = metav1.Duration{Duration: 5 * time.Minute}
		return controllerutil.SetOwnerReference(pkg, &repository, client.Scheme())
	})
	return err
}

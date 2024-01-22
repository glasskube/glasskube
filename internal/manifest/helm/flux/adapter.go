package flux

import (
	"context"
	"fmt"
	"time"

	helmv1beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type FluxHelmAdapter struct {
	scheme *runtime.Scheme
}

func NewAdapter(scheme *runtime.Scheme) manifest.ManifestAdapter {
	return &FluxHelmAdapter{scheme: scheme}
}

func (a FluxHelmAdapter) ControllerInit(builder *builder.Builder) {
	sourcev1beta2.AddToScheme(a.scheme)
	helmv1beta2.AddToScheme(a.scheme)
	builder.Owns(&sourcev1beta2.HelmRepository{})
	builder.Owns(&helmv1beta2.HelmRelease{})
	builder.Owns(&corev1.Namespace{})
}

func (a FluxHelmAdapter) Reconcile(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*result.ReconcileResult, error) {
	if namespace, err := ensureNamespace(ctx, client, pkg, manifest); err != nil {
		return nil, err
	} else if namespace.Status.Phase == corev1.NamespaceTerminating {
		return result.Waiting("Namespace is still terminating"), nil
	}
	if err := ensureHelmRepository(ctx, client, pkg, manifest); err != nil {
		return nil, err
	}
	if helmRelease, err := ensureHelmRelease(ctx, client, pkg, manifest); err != nil {
		return nil, err
	} else {
		return extractResult(helmRelease), nil
	}
}

func ensureNamespace(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*corev1.Namespace, error) {
	log := log.FromContext(ctx)
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: manifest.DefaultNamespace,
		},
	}
	result, err := controllerutil.CreateOrUpdate(ctx, client, &namespace, func() error {
		if namespace.Status.Phase == corev1.NamespaceTerminating {
			return nil
		} else {
			return controllerutil.SetOwnerReference(pkg, &namespace, client.Scheme())
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not ensure namespace: %w", err)
	} else {
		log.V(1).Info("CreateUrUpdate result: "+string(result), "Namespace", namespace.Name)
		return &namespace, nil
	}
}

func ensureHelmRepository(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) error {
	log := log.FromContext(ctx)
	helmRepository := sourcev1beta2.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifest.Name,
			Namespace: manifest.DefaultNamespace,
		},
	}
	result, err := controllerutil.CreateOrUpdate(ctx, client, &helmRepository, func() error {
		helmRepository.Spec.URL = manifest.Helm.RepositoryUrl
		helmRepository.Spec.Interval = metav1.Duration{Duration: 1 * time.Hour}
		return controllerutil.SetOwnerReference(pkg, &helmRepository, client.Scheme())
	})
	if err != nil {
		return fmt.Errorf("could not ensure helm repository: %w", err)
	} else {
		log.V(1).Info("CreateUrUpdate result: "+string(result), "HelmRepository", helmRepository.Name)
		return err
	}
}

func ensureHelmRelease(ctx context.Context, client client.Client, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*helmv1beta2.HelmRelease, error) {
	log := log.FromContext(ctx)
	helmRelease := helmv1beta2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifest.Name,
			Namespace: manifest.DefaultNamespace,
		},
	}
	result, err := controllerutil.CreateOrUpdate(ctx, client, &helmRelease, func() error {
		helmRelease.Spec.Chart.Spec.Chart = manifest.Helm.ChartName
		helmRelease.Spec.Chart.Spec.Version = manifest.Helm.ChartVersion
		helmRelease.Spec.Chart.Spec.SourceRef.Kind = "HelmRepository"
		helmRelease.Spec.Chart.Spec.SourceRef.Name = manifest.Name
		helmRelease.Spec.Values = &apiextensionsv1.JSON{Raw: manifest.Helm.Values.Raw[:]}
		helmRelease.Spec.Interval = metav1.Duration{Duration: 5 * time.Minute}
		return controllerutil.SetOwnerReference(pkg, &helmRelease, client.Scheme())
	})
	if err != nil {
		return nil, fmt.Errorf("could not ensure helm release: %w", err)
	} else {
		log.V(1).Info("CreateUrUpdate result: "+string(result), "HelmRelease", helmRelease.Name)
		return &helmRelease, nil
	}
}

func extractResult(helmRelease *helmv1beta2.HelmRelease) *result.ReconcileResult {
	if readyCondition := meta.FindStatusCondition(helmRelease.Status.Conditions, "Ready"); readyCondition != nil {
		if readyCondition.Status == metav1.ConditionTrue {
			return result.Ready("flux: " + readyCondition.Message)
		} else if readyCondition.Status == metav1.ConditionFalse {
			return result.Failed("flux: " + readyCondition.Message)
		}
	}
	if reconcilingCondition := meta.FindStatusCondition(helmRelease.Status.Conditions, "Reconciling"); reconcilingCondition != nil {
		return result.Waiting("flux: " + reconcilingCondition.Message)
	} else {
		return result.Waiting("Waiting for HelmRelease reconciliation")
	}
}

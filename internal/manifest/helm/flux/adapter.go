package flux

import (
	"context"
	"fmt"
	"strings"
	"time"

	helmv1beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type FluxHelmAdapter struct {
	client.Client
	*owners.OwnerManager
}

func NewAdapter() manifest.ManifestAdapter {
	return &FluxHelmAdapter{}
}

func (a *FluxHelmAdapter) ControllerInit(buildr *builder.Builder, client client.Client, scheme *runtime.Scheme) error {
	if err := sourcev1beta2.AddToScheme(scheme); err != nil {
		return err
	}
	if err := helmv1beta2.AddToScheme(scheme); err != nil {
		return err
	}
	if a.OwnerManager == nil {
		a.OwnerManager = owners.NewOwnerManager(scheme)
	}
	a.Client = client
	buildr.Owns(&sourcev1beta2.HelmRepository{})
	buildr.Owns(&helmv1beta2.HelmRelease{}, builder.MatchEveryOwner)
	buildr.Owns(&corev1.Namespace{})
	return nil
}

func (a *FluxHelmAdapter) Reconcile(ctx context.Context, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*result.ReconcileResult, error) {
	if namespace, err := a.ensureNamespace(ctx, pkg, manifest); err != nil {
		return nil, err
	} else if namespace.Status.Phase == corev1.NamespaceTerminating {
		return result.Waiting("Namespace is still terminating"), nil
	}
	if err := a.ensureHelmRepository(ctx, pkg, manifest); err != nil {
		return nil, err
	}
	if helmRelease, err := a.ensureHelmRelease(ctx, pkg, manifest); err != nil {
		return nil, err
	} else {
		return extractResult(helmRelease), nil
	}
}

func (a *FluxHelmAdapter) ensureNamespace(ctx context.Context, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*corev1.Namespace, error) {
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: manifest.DefaultNamespace,
		},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("Namespace", namespace.Name)
	result, err := controllerutil.CreateOrUpdate(ctx, a.Client, &namespace, func() error {
		if namespace.Status.Phase == corev1.NamespaceTerminating {
			return nil
		} else {
			return a.SetOwner(pkg, &namespace, owners.BlockOwnerDeletion)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not ensure namespace: %w", err)
	} else {
		log.V(1).Info("ensured Namespace", "result", result)
		return &namespace, nil
	}
}

func (a *FluxHelmAdapter) ensureHelmRepository(ctx context.Context, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) error {
	helmRepository := sourcev1beta2.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifest.Name,
			Namespace: manifest.DefaultNamespace,
		},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("HelmRepository", helmRepository.Name)
	result, err := controllerutil.CreateOrUpdate(ctx, a.Client, &helmRepository, func() error {
		helmRepository.Spec.URL = manifest.Helm.RepositoryUrl
		helmRepository.Spec.Interval = metav1.Duration{Duration: 1 * time.Hour}
		return a.SetOwner(pkg, &helmRepository, owners.BlockOwnerDeletion)
	})
	if err != nil {
		return fmt.Errorf("could not ensure helm repository: %w", err)
	} else {
		log.V(1).Info("ensured HelmRepository", "result", result)
		return err
	}
}

func (a *FluxHelmAdapter) ensureHelmRelease(ctx context.Context, pkg *packagesv1alpha1.Package, manifest *packagesv1alpha1.PackageManifest) (*helmv1beta2.HelmRelease, error) {
	helmRelease := helmv1beta2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifest.Name,
			Namespace: manifest.DefaultNamespace,
		},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("HelmRelease", helmRelease.Name)
	result, err := controllerutil.CreateOrUpdate(ctx, a.Client, &helmRelease, func() error {
		helmRelease.Spec.Chart.Spec.Chart = manifest.Helm.ChartName
		helmRelease.Spec.Chart.Spec.Version = manifest.Helm.ChartVersion
		helmRelease.Spec.Chart.Spec.SourceRef.Kind = "HelmRepository"
		helmRelease.Spec.Chart.Spec.SourceRef.Name = manifest.Name
		if manifest.Helm.Values != nil {
			helmRelease.Spec.Values = &apiextensionsv1.JSON{Raw: manifest.Helm.Values.Raw[:]}
		} else {
			helmRelease.Spec.Values = nil
		}
		helmRelease.Spec.Interval = metav1.Duration{Duration: 5 * time.Minute}
		return a.SetOwner(pkg, &helmRelease, owners.BlockOwnerDeletion)
	})
	if err != nil {
		return nil, fmt.Errorf("could not ensure helm release: %w", err)
	} else {
		log.V(1).Info("ensured HelmRelease", "result", result)
		return &helmRelease, nil
	}
}

func extractResult(helmRelease *helmv1beta2.HelmRelease) *result.ReconcileResult {
	if readyCondition := meta.FindStatusCondition(helmRelease.Status.Conditions, "Ready"); readyCondition != nil {
		if readyCondition.Status == metav1.ConditionTrue {
			return result.Ready("flux: " + readyCondition.Message)
		} else if readyCondition.Status == metav1.ConditionFalse {
			if strings.Contains(readyCondition.Message, "latest generation of object has not been reconciled") {
				return result.Waiting("flux: " + readyCondition.Message)
			} else {
				return result.Failed("flux: " + readyCondition.Message)
			}
		}
	}
	if reconcilingCondition := meta.FindStatusCondition(helmRelease.Status.Conditions, "Reconciling"); reconcilingCondition != nil {
		return result.Waiting("flux: " + reconcilingCondition.Message)
	} else {
		return result.Waiting("Waiting for HelmRelease reconciliation")
	}
}

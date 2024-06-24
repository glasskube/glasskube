package flux

import (
	"context"
	"fmt"
	"strings"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	packagesv1alpha1 "github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/controller/labels"
	"github.com/glasskube/glasskube/internal/controller/owners"
	"github.com/glasskube/glasskube/internal/controller/owners/utils"
	"github.com/glasskube/glasskube/internal/manifest"
	"github.com/glasskube/glasskube/internal/manifest/result"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/names"
	corev1 "k8s.io/api/core/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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
	if err := helmv2.AddToScheme(scheme); err != nil {
		return err
	}
	if a.OwnerManager == nil {
		a.OwnerManager = owners.NewOwnerManager(scheme)
	}
	a.Client = client
	buildr.Owns(&sourcev1beta2.HelmRepository{})
	buildr.Owns(&helmv2.HelmRelease{}, builder.MatchEveryOwner)
	buildr.Owns(&corev1.Namespace{})
	return nil
}

func (a *FluxHelmAdapter) Reconcile(
	ctx context.Context,
	pkg ctrlpkg.Package,
	manifest *packagesv1alpha1.PackageManifest,
	patches manifestvalues.TargetPatches,
) (*result.ReconcileResult, error) {
	log := ctrl.LoggerFrom(ctx)
	var ownedResources []packagesv1alpha1.OwnedResourceRef
	if !pkg.IsNamespaceScoped() {
		if namespace, err := a.ensureNamespace(ctx, pkg, manifest); err != nil {
			return nil, err
		} else {
			if _, err := utils.AddOwnedResourceRef(a.Scheme(), &ownedResources, namespace); err != nil {
				log.Error(err, "could not add Namespace to ownedResources")
			}
			if namespace.Status.Phase == corev1.NamespaceTerminating {
				return result.Waiting("Namespace is still terminating", ownedResources), nil
			}
		}
	}
	if helmRepository, err := a.ensureHelmRepository(ctx, pkg, manifest); err != nil {
		return nil, err
	} else {
		if _, err := utils.AddOwnedResourceRef(a.Scheme(), &ownedResources, helmRepository); err != nil {
			log.Error(err, "could not add HelmRepository to ownedResources")
		}
	}
	if helmRelease, err := a.ensureHelmRelease(ctx, pkg, manifest, patches); err != nil {
		return nil, err
	} else {
		if _, err := utils.AddOwnedResourceRef(a.Scheme(), &ownedResources, helmRelease); err != nil {
			log.Error(err, "could not add HelmRelease to ownedResources")
		}
		return extractResult(helmRelease, ownedResources), nil
	}
}

func (a *FluxHelmAdapter) ensureNamespace(
	ctx context.Context,
	pkg ctrlpkg.Package,
	manifest *packagesv1alpha1.PackageManifest,
) (*corev1.Namespace, error) {
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
			return a.SetOwnerIfManagedOrNotExists(a.Client, ctx, pkg, &namespace)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not ensure namespace: %w", err)
	} else {
		log.V(1).Info("ensured Namespace", "result", result)
		return &namespace, nil
	}
}

func (a *FluxHelmAdapter) ensureHelmRepository(
	ctx context.Context,
	pkg ctrlpkg.Package,
	manifest *packagesv1alpha1.PackageManifest,
) (*sourcev1beta2.HelmRepository, error) {
	var namespace string
	if pkg.IsNamespaceScoped() {
		namespace = pkg.GetNamespace()
	} else {
		namespace = manifest.DefaultNamespace
	}
	helmRepository := sourcev1beta2.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.HelmResourceName(pkg, manifest),
			Namespace: namespace,
		},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("HelmRepository", helmRepository.Name)
	result, err := controllerutil.CreateOrUpdate(ctx, a.Client, &helmRepository, func() error {
		helmRepository.Spec.URL = manifest.Helm.RepositoryUrl
		helmRepository.Spec.Interval = metav1.Duration{Duration: 1 * time.Hour}
		labels.SetManaged(&helmRepository)
		return a.SetOwner(pkg, &helmRepository, owners.BlockOwnerDeletion)
	})
	if err != nil {
		return nil, fmt.Errorf("could not ensure helm repository: %w", err)
	} else {
		log.V(1).Info("ensured HelmRepository", "result", result)
		return &helmRepository, nil
	}
}

func (a *FluxHelmAdapter) ensureHelmRelease(
	ctx context.Context,
	pkg ctrlpkg.Package,
	manifest *packagesv1alpha1.PackageManifest,
	patches manifestvalues.TargetPatches,
) (*helmv2.HelmRelease, error) {
	var namespace string
	if pkg.IsNamespaceScoped() {
		namespace = pkg.GetNamespace()
	} else {
		namespace = manifest.DefaultNamespace
	}
	helmRelease := helmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.HelmResourceName(pkg, manifest),
			Namespace: namespace,
		},
	}
	log := ctrl.LoggerFrom(ctx).WithValues("HelmRelease", helmRelease.Name)
	result, err := controllerutil.CreateOrUpdate(ctx, a.Client, &helmRelease, func() error {
		if helmRelease.Spec.Chart == nil {
			helmRelease.Spec.Chart = &helmv2.HelmChartTemplate{}
		}
		helmRelease.Spec.Chart.Spec.Chart = manifest.Helm.ChartName
		helmRelease.Spec.Chart.Spec.Version = manifest.Helm.ChartVersion
		helmRelease.Spec.Chart.Spec.SourceRef.Kind = "HelmRepository"
		helmRelease.Spec.Chart.Spec.SourceRef.Name = names.HelmResourceName(pkg, manifest)
		if manifest.Helm.Values != nil {
			helmRelease.Spec.Values = &extv1.JSON{Raw: manifest.Helm.Values.Raw[:]}
		} else {
			helmRelease.Spec.Values = nil
		}
		if err := patches.ApplyToHelmRelease(&helmRelease); err != nil {
			return err
		}
		helmRelease.Spec.Interval = metav1.Duration{Duration: 5 * time.Minute}
		labels.SetManaged(&helmRelease)
		return a.SetOwner(pkg, &helmRelease, owners.BlockOwnerDeletion)
	})
	if err != nil {
		return nil, fmt.Errorf("could not ensure helm release: %w", err)
	} else {
		log.V(1).Info("ensured HelmRelease", "result", result)
		return &helmRelease, nil
	}
}

func extractResult(
	helmRelease *helmv2.HelmRelease,
	ownedResources []packagesv1alpha1.OwnedResourceRef,
) *result.ReconcileResult {
	if readyCondition := meta.FindStatusCondition(helmRelease.Status.Conditions, "Ready"); readyCondition != nil {
		message := fmt.Sprintf("flux: %v", readyCondition.Message)
		if readyCondition.Status == metav1.ConditionTrue {
			// TODO: Add HelmRepository, HelmRelease to owned resources
			return result.Ready(message, ownedResources)
		} else if readyCondition.Status == metav1.ConditionFalse {
			if strings.Contains(readyCondition.Message, "latest generation of object has not been reconciled") {
				return result.Waiting(message, ownedResources)
			} else {
				return result.Failed(message, ownedResources)
			}
		}
	}
	reconcilingCondition := meta.FindStatusCondition(helmRelease.Status.Conditions, "Reconciling")
	if reconcilingCondition != nil {
		return result.Waiting("flux: "+reconcilingCondition.Message, ownedResources)
	} else {
		return result.Waiting("Waiting for HelmRelease reconciliation", ownedResources)
	}
}

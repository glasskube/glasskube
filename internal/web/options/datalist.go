package options

import (
	"context"
	"fmt"
	"slices"

	"github.com/glasskube/glasskube/internal/clicontext"
	webcontext "github.com/glasskube/glasskube/internal/web/context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/maputils"
	"github.com/glasskube/glasskube/internal/web/components"
	"github.com/glasskube/glasskube/pkg/manifest"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func GetDatalistOptions(ctx context.Context, ref *v1alpha1.ValueReference, namespaceOptions []string, pkgsOptions []string) (
	*components.PkgConfigInputDatalistOptions, error) {
	datalistOptions := components.PkgConfigInputDatalistOptions{}
	if ref.ConfigMapRef != nil {
		datalistOptions.Namespaces = namespaceOptions
		if nameOptions, err := GetConfigMapNameOptions(ctx, ref.ConfigMapRef.Namespace); err != nil {
			return &datalistOptions, fmt.Errorf("error fetching ConfigMaps of namespace %v: %w", ref.ConfigMapRef.Namespace, err)
		} else {
			datalistOptions.Names = nameOptions
			if keyOptions, err := GetConfigMapKeyOptions(ctx, ref.ConfigMapRef.Namespace, ref.ConfigMapRef.Name); err != nil {
				return &datalistOptions, fmt.Errorf("error fetching keys of ConfigMap %v in namespace %v: %w", ref.ConfigMapRef.Name, ref.ConfigMapRef.Namespace, err)
			} else {
				datalistOptions.Keys = keyOptions
			}
		}
	} else if ref.SecretRef != nil {
		datalistOptions.Namespaces = namespaceOptions
		if nameOptions, err := GetSecretNameOptions(ctx, ref.SecretRef.Namespace); err != nil {
			return &datalistOptions, fmt.Errorf("error fetching Secrets of namespace %v: %w", ref.SecretRef.Namespace, err)
		} else {
			datalistOptions.Names = nameOptions
			if keyOptions, err := GetSecretKeyOptions(ctx, ref.SecretRef.Namespace, ref.SecretRef.Name); err != nil {
				return &datalistOptions, fmt.Errorf("error fetching keys of Secret %v in namespace %v: %w", ref.SecretRef.Name, ref.SecretRef.Namespace, err)
			} else {
				datalistOptions.Keys = keyOptions
			}
		}
	} else if ref.PackageRef != nil {
		datalistOptions.Names = pkgsOptions
		if keyOptions, err := GetPackageValuesOptions(ctx, ref.PackageRef.Name); err != nil {
			return &datalistOptions, fmt.Errorf("error fetching installed manifest of %v", ref.PackageRef.Name)
		} else {
			datalistOptions.Keys = keyOptions
		}
	}
	return &datalistOptions, nil
}

func GetNamespaceOptions(ctx context.Context) ([]string, error) {
	coreListers := webcontext.CoreListersFromContext(ctx)
	if namespaces, err := (*coreListers.NamespaceLister).List(labels.NewSelector()); err != nil {
		return nil, err
	} else {
		return sortedNames(namespaces), nil
	}
}

func GetPackagesOptions(ctx context.Context) ([]string, error) {
	var packages v1alpha1.ClusterPackageList
	pkgClient := clicontext.PackageClientFromContext(ctx)
	if err := pkgClient.ClusterPackages().GetAll(ctx, &packages); err != nil {
		return nil, err
	} else {
		options := make([]string, 0)
		for _, pkg := range packages.Items {
			options = append(options, pkg.Name)
		}
		slices.Sort(options)
		return options, nil
	}
}

func GetConfigMapNameOptions(ctx context.Context, namespace string) ([]string, error) {
	if namespace != "" {
		coreListers := webcontext.CoreListersFromContext(ctx)
		if configMaps, err := (*coreListers.ConfigMapLister).ConfigMaps(namespace).List(labels.NewSelector()); err != nil {
			return nil, err
		} else {
			return sortedNames(configMaps), nil
		}
	}
	return nil, nil
}

func GetConfigMapKeyOptions(ctx context.Context, namespace string, name string) ([]string, error) {
	if namespace != "" && name != "" {
		coreLister := webcontext.CoreListersFromContext(ctx)
		if configMap, err := (*coreLister.ConfigMapLister).ConfigMaps(namespace).Get(name); err != nil {
			return nil, err
		} else {
			return maputils.KeysSorted(configMap.Data), nil
		}
	}
	return nil, nil
}

func GetSecretNameOptions(ctx context.Context, namespace string) ([]string, error) {
	if namespace != "" {
		coreListers := webcontext.CoreListersFromContext(ctx)
		if secrets, err := (*coreListers.SecretLister).Secrets(namespace).List(labels.NewSelector()); err != nil {
			return nil, err
		} else {
			return sortedNames(secrets), nil
		}
	}
	return nil, nil
}

func GetSecretKeyOptions(ctx context.Context, namespace string, name string) ([]string, error) {
	if namespace != "" && name != "" {
		coreListers := webcontext.CoreListersFromContext(ctx)
		if secret, err := (*coreListers.SecretLister).Secrets(namespace).Get(name); err != nil {
			return nil, err
		} else {
			return maputils.KeysSorted(secret.Data), nil
		}
	}
	return nil, nil
}

func GetPackageValuesOptions(ctx context.Context, pkgName string) ([]string, error) {
	if pkgName != "" {
		if refManifest, err := manifest.GetInstalledManifest(ctx, pkgName); err != nil {
			return nil, err
		} else {
			return maputils.KeysSorted(refManifest.ValueDefinitions), nil
		}
	}
	return nil, nil
}

func sortedNames[T v1.Object](objects []T) []string {
	names := make([]string, 0, len(objects))
	for _, obj := range objects {
		names = append(names, obj.GetName())
	}
	slices.Sort(names)
	return names
}

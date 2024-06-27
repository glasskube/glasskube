package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/glasskube/glasskube/pkg/manifest"

	"github.com/glasskube/glasskube/internal/maputils"
	"github.com/glasskube/glasskube/internal/web/components/pkg_config_input"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

const (
	namespaceKey     = "namespace"
	nameKey          = "name"
	keyKey           = "key"
	packageKey       = "package"
	valueKey         = "value"
	refKindKey       = "refKind"
	refKindConfigMap = "ConfigMap"
	refKindSecret    = "Secret"
	refKindPackage   = "Package"
)

func formKey(valueName string, key string) string {
	return fmt.Sprintf("%s[%s]", valueName, key)
}

func extractValues(r *http.Request, manifest *v1alpha1.PackageManifest) (map[string]v1alpha1.ValueConfiguration, error) {
	values := make(map[string]v1alpha1.ValueConfiguration)
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	for valueName, valueDef := range manifest.ValueDefinitions {
		if refKindVal := r.Form.Get(fmt.Sprintf("%s[%s]", valueName, refKindKey)); refKindVal == refKindConfigMap {
			values[valueName] = v1alpha1.ValueConfiguration{
				ValueFrom: &v1alpha1.ValueReference{
					ConfigMapRef: extractObjectKeyValueSource(r, valueName),
				},
			}
		} else if refKindVal == refKindSecret {
			values[valueName] = v1alpha1.ValueConfiguration{
				ValueFrom: &v1alpha1.ValueReference{
					SecretRef: extractObjectKeyValueSource(r, valueName),
				},
			}
		} else if refKindVal == refKindPackage {
			values[valueName] = v1alpha1.ValueConfiguration{
				ValueFrom: &v1alpha1.ValueReference{
					PackageRef: extractPackageValueSource(r, valueName),
				},
			}
		} else if refKindVal == "" {
			formVal := r.Form.Get(valueName)
			if valueDef.Type == v1alpha1.ValueTypeBoolean {
				boolStr := strconv.FormatBool(false)
				if strings.ToLower(formVal) == "on" {
					boolStr = strconv.FormatBool(true)
				}
				values[valueName] = v1alpha1.ValueConfiguration{Value: &boolStr}
			} else {
				values[valueName] = v1alpha1.ValueConfiguration{Value: &formVal}
			}
		} else {
			return nil, fmt.Errorf("cannot extract value %v because of unknown reference kind %v", valueName, refKindVal)
		}
	}
	return values, nil
}

func extractObjectKeyValueSource(r *http.Request, valueName string) *v1alpha1.ObjectKeyValueSource {
	namespaceFormKey := formKey(valueName, namespaceKey)
	nameFormKey := formKey(valueName, nameKey)
	keyFormKey := formKey(valueName, keyKey)
	return &v1alpha1.ObjectKeyValueSource{
		Name:      r.Form.Get(nameFormKey),
		Namespace: r.Form.Get(namespaceFormKey),
		Key:       r.Form.Get(keyFormKey),
	}
}

func extractPackageValueSource(r *http.Request, valueName string) *v1alpha1.PackageValueSource {
	packageFormKey := formKey(valueName, packageKey)
	valueFormKey := formKey(valueName, valueKey)
	return &v1alpha1.PackageValueSource{
		Name:  r.Form.Get(packageFormKey),
		Value: r.Form.Get(valueFormKey),
	}
}

// packageConfigurationInput is like clusterPackageConfigurationInput but for packages
func (s *server) packageConfigurationInput(w http.ResponseWriter, r *http.Request) {
	manifestName := mux.Vars(r)["manifestName"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	selectedVersion := r.FormValue("selectedVersion")
	repositoryName := r.FormValue("repositoryName")
	pkg, manifest, err := describe.DescribeInstalledPackage(r.Context(), namespace, name)
	if err != nil && !errors.IsNotFound(err) {
		err = fmt.Errorf("an error occurred fetching package details of %v: %w", manifestName, err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	s.handleConfigurationInput(w, r, &packageDetailPageContext{
		repositoryName:  repositoryName,
		selectedVersion: selectedVersion,
		manifestName:    manifestName,
		pkg:             pkg,
		manifest:        manifest,
	})
}

// clusterPackageConfigurationInput is a GET endpoint, which returns an html snippet containing an input container.
// The endpoint requires the pkgName query parameter to be set, as well as the valueName query parameter (which holds
// the name of the desired value according to the package value definitions).
// An optional query parameter refKind can be passed to request the snippet in a certain variant, where the accepted
// refKind values are: ConfigMap, Secret, Package. If no refKind is given, the "regular" input is returned.
// In any case, the input container consists of a button where the user can change the type of reference or remove the
// reference, and the actual input field(s).
func (s *server) clusterPackageConfigurationInput(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkgName"]
	selectedVersion := r.FormValue("selectedVersion")
	repositoryName := r.FormValue("repositoryName")
	pkg, manifest, err := describe.DescribeInstalledClusterPackage(r.Context(), pkgName)
	if err != nil && !errors.IsNotFound(err) {
		err = fmt.Errorf("an error occurred fetching package details of %v: %w", pkgName, err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	s.handleConfigurationInput(w, r, &packageDetailPageContext{
		repositoryName:  repositoryName,
		selectedVersion: selectedVersion,
		manifestName:    pkgName,
		pkg:             pkg,
		manifest:        manifest,
	})
}

func (s *server) handleConfigurationInput(w http.ResponseWriter, r *http.Request, d *packageDetailPageContext) {
	if d.manifest == nil {
		d.manifest = &v1alpha1.PackageManifest{}
		if err := s.repoClientset.ForRepoWithName(d.repositoryName).
			FetchPackageManifest(d.manifestName, d.selectedVersion, d.manifest); err != nil {
			// TODO check error handling again?
			s.respondAlertAndLog(w, err,
				fmt.Sprintf("An error occurred fetching manifest of %v in version %v", d.manifestName, d.selectedVersion),
				"danger")
			return
		}
	}

	valueName := mux.Vars(r)["valueName"]
	refKind := r.URL.Query().Get("refKind")
	if valueDefinition, ok := d.manifest.ValueDefinitions[valueName]; ok {
		options := pkg_config_input.PkgConfigInputDatalistOptions{}
		if refKind == refKindConfigMap || refKind == refKindSecret {
			if opts, err := s.getNamespaceOptions(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get namespace options: %v\n", err)
			} else {
				options.Namespaces = opts
			}
		} else {
			if opts, err := s.getPackagesOptions(r.Context()); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get package options: %v\n", err)
			} else {
				options.Names = opts
			}
		}
		input := pkg_config_input.ForPkgConfigInput(
			d.pkg, d.repositoryName, d.selectedVersion, d.manifest, valueName, valueDefinition, nil, &options,
			&pkg_config_input.PkgConfigInputRenderOptions{
				Autofocus:      true,
				DesiredRefKind: &refKind,
			})
		err := s.templates.pkgConfigInput.Execute(w, input)
		checkTmplError(err, fmt.Sprintf("package config input (%s, %s)", d.manifestName, valueName))
	}
}

// namesDatalist is a GET endpoint returning an html datalist, containing options depending on the given valueName,
// kind of reference and namespace. It is only usable for ConfigMap and Secret refs, since packages don't have a
// namespace. In case the refKind is ConfigMap, the datalist contains the config maps of the given namespace; in case
// the refKind is Secret, the datalist contains the secrets of the given namespace; if no namespace is given or an
// error occurs, an empty datalist is returned
func (s *server) namesDatalist(w http.ResponseWriter, r *http.Request) {
	valueName := mux.Vars(r)["valueName"]
	refKind := r.FormValue(refKindKey)
	id := r.FormValue("id")
	nsKey := formKey(valueName, namespaceKey)
	namespace := r.Form.Get(nsKey)
	var options []string
	if refKind == refKindConfigMap {
		if opts, err := s.getConfigMapNameOptions(namespace); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get config map name options: %v\n", err)
		} else {
			options = opts
		}
	} else if refKind == refKindSecret {
		if opts, err := s.getSecretNameOptions(namespace); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get secret name options: %v\n", err)
		} else {
			options = opts
		}
	}
	tmplErr := s.templates.datalistTmpl.Execute(w, map[string]any{
		"Options": options,
		"Id":      id,
	})
	checkTmplError(tmplErr, "names-datalist")
}

func (s *server) keysDatalist(w http.ResponseWriter, r *http.Request) {
	valueName := mux.Vars(r)["valueName"]
	refKind := r.FormValue(refKindKey)
	nsKey := formKey(valueName, namespaceKey)
	nameKey := formKey(valueName, nameKey)
	pkgKey := formKey(valueName, packageKey)
	namespace := r.Form.Get(nsKey)
	name := r.Form.Get(nameKey)
	pkg := r.Form.Get(pkgKey)
	var options []string
	var err error
	if refKind == refKindConfigMap {
		if options, err = s.getConfigMapKeyOptions(namespace, name); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get ConfigMap key options of %v in %v: %v\n", name, namespace, err)
		}
	} else if refKind == refKindSecret {
		if options, err = s.getSecretKeyOptions(namespace, name); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get Secret key options of %v in %v: %v\n", name, namespace, err)
		}
	} else if refKind == refKindPackage {
		if options, err = s.getPackageValuesOptions(r.Context(), pkg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get package value options of %v: %v\n", pkg, err)
		}
	}
	tmplErr := s.templates.datalistTmpl.Execute(w, map[string]any{
		"Options": options,
		"Id":      r.FormValue("id"),
	})
	checkTmplError(tmplErr, "keys-datalist")
}

func (s *server) getDatalistOptions(ctx context.Context, ref *v1alpha1.ValueReference, namespaceOptions []string, pkgsOptions []string) (*pkg_config_input.PkgConfigInputDatalistOptions, error) {
	datalistOptions := pkg_config_input.PkgConfigInputDatalistOptions{}
	if ref.ConfigMapRef != nil {
		datalistOptions.Namespaces = namespaceOptions
		if nameOptions, err := s.getConfigMapNameOptions(ref.ConfigMapRef.Namespace); err != nil {
			return &datalistOptions, fmt.Errorf("error fetching ConfigMaps of namespace %v: %w", ref.ConfigMapRef.Namespace, err)
		} else {
			datalistOptions.Names = nameOptions
			if keyOptions, err := s.getConfigMapKeyOptions(ref.ConfigMapRef.Namespace, ref.ConfigMapRef.Name); err != nil {
				return &datalistOptions, fmt.Errorf("error fetching keys of ConfigMap %v in namespace %v: %w", ref.ConfigMapRef.Name, ref.ConfigMapRef.Namespace, err)
			} else {
				datalistOptions.Keys = keyOptions
			}
		}
	} else if ref.SecretRef != nil {
		datalistOptions.Namespaces = namespaceOptions
		if nameOptions, err := s.getSecretNameOptions(ref.SecretRef.Namespace); err != nil {
			return &datalistOptions, fmt.Errorf("error fetching Secrets of namespace %v: %w", ref.SecretRef.Namespace, err)
		} else {
			datalistOptions.Names = nameOptions
			if keyOptions, err := s.getSecretKeyOptions(ref.SecretRef.Namespace, ref.SecretRef.Name); err != nil {
				return &datalistOptions, fmt.Errorf("error fetching keys of Secret %v in namespace %v: %w", ref.SecretRef.Name, ref.SecretRef.Namespace, err)
			} else {
				datalistOptions.Keys = keyOptions
			}
		}
	} else if ref.PackageRef != nil {
		datalistOptions.Names = pkgsOptions
		if keyOptions, err := s.getPackageValuesOptions(ctx, ref.PackageRef.Name); err != nil {
			return &datalistOptions, fmt.Errorf("error fetching installed manifest of %v", ref.PackageRef.Name)
		} else {
			datalistOptions.Keys = keyOptions
		}
	}
	return &datalistOptions, nil
}

func (s *server) getNamespaceOptions() ([]string, error) {
	if namespaces, err := (*s.namespaceLister).List(labels.NewSelector()); err != nil {
		return nil, err
	} else {
		return sortedNames(namespaces), nil
	}
}

func (s *server) getPackagesOptions(ctx context.Context) ([]string, error) {
	var packages v1alpha1.ClusterPackageList
	if err := s.pkgClient.ClusterPackages().GetAll(ctx, &packages); err != nil {
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

func (s *server) getConfigMapNameOptions(namespace string) ([]string, error) {
	if namespace != "" {
		if configMaps, err := (*s.configMapLister).ConfigMaps(namespace).List(labels.NewSelector()); err != nil {
			return nil, err
		} else {
			return sortedNames(configMaps), nil
		}
	}
	return nil, nil
}

func (s *server) getConfigMapKeyOptions(namespace string, name string) ([]string, error) {
	if namespace != "" && name != "" {
		if configMap, err := (*s.configMapLister).ConfigMaps(namespace).Get(name); err != nil {
			return nil, err
		} else {
			return maputils.KeysSorted(configMap.Data), nil
		}
	}
	return nil, nil
}

func (s *server) getSecretNameOptions(namespace string) ([]string, error) {
	if namespace != "" {
		if secrets, err := (*s.secretLister).Secrets(namespace).List(labels.NewSelector()); err != nil {
			return nil, err
		} else {
			return sortedNames(secrets), nil
		}
	}
	return nil, nil
}

func (s *server) getSecretKeyOptions(namespace string, name string) ([]string, error) {
	if namespace != "" && name != "" {
		if secret, err := (*s.secretLister).Secrets(namespace).Get(name); err != nil {
			return nil, err
		} else {
			return maputils.KeysSorted(secret.Data), nil
		}
	}
	return nil, nil
}

func (s *server) getPackageValuesOptions(ctx context.Context, pkgName string) ([]string, error) {
	if pkgName != "" {
		if refManifest, err := manifest.GetInstalledManifest(ctx, pkgName); err != nil {
			return nil, err
		} else {
			return maputils.KeysSorted(refManifest.ValueDefinitions), nil
		}
	}
	return nil, nil
}

func sortedNames[T metav1.Object](objects []T) []string {
	names := make([]string, 0)
	for _, obj := range objects {
		names = append(names, obj.GetName())
	}
	slices.Sort(names)
	return names
}

package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/web/components"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	opts "github.com/glasskube/glasskube/internal/web/options"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/pkg/describe"
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	formValuePrefix  = "values"
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
	return fmt.Sprintf("%s.%s[%s]", formValuePrefix, valueName, key)
}

// extractValues extracts dynamic package configuration values from the form of the given request, such that installation
// or configuration can be done with the provided values
func extractValues(r *http.Request, manifest *v1alpha1.PackageManifest) (map[string]v1alpha1.ValueConfiguration, error) {
	values := make(map[string]v1alpha1.ValueConfiguration)
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	for valueName, valueDef := range manifest.ValueDefinitions {
		if refKindVal := r.Form.Get(fmt.Sprintf("%s.%s[%s]", formValuePrefix, valueName, refKindKey)); refKindVal == refKindConfigMap {
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
			formVal := r.Form.Get(fmt.Sprintf("%v.%v", formValuePrefix, valueName))
			if valueDef.Type == v1alpha1.ValueTypeBoolean {
				boolStr := strconv.FormatBool(false)
				if strings.ToLower(formVal) == "on" {
					boolStr = strconv.FormatBool(true)
				}
				values[valueName] = v1alpha1.ValueConfiguration{InlineValueConfiguration: v1alpha1.InlineValueConfiguration{Value: &boolStr}}
			} else {
				values[valueName] = v1alpha1.ValueConfiguration{InlineValueConfiguration: v1alpha1.InlineValueConfiguration{Value: &formVal}}
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

// GetPackageConfigurationInput is like GetClusterPackageConfigurationInput but for packages
func GetPackageConfigurationInput(w http.ResponseWriter, r *http.Request) {
	pCtx := getPackageContext(r)
	pkg, manifest, err := describe.DescribeInstalledPackage(r.Context(), pCtx.request.namespace, pCtx.request.name)
	if err != nil && !errors.IsNotFound(err) {
		err = fmt.Errorf("an error occurred fetching package details of %v: %w", pCtx.request.manifestName, err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	handleConfigurationInput(w, r, &packageContext{
		request:  pCtx.request,
		pkg:      pkg,
		manifest: manifest,
	})
}

// GetClusterPackageConfigurationInput is a GET endpoint, which returns an html snippet containing an input container.
// The endpoint requires the pkgName query parameter to be set, as well as the valueName query parameter (which holds
// the name of the desired value according to the package value definitions).
// An optional query parameter refKind can be passed to request the snippet in a certain variant, where the accepted
// refKind values are: ConfigMap, Secret, Package. If no refKind is given, the "regular" input is returned.
// In any case, the input container consists of a button where the user can change the type of reference or remove the
// reference, and the actual input field(s).
func GetClusterPackageConfigurationInput(w http.ResponseWriter, r *http.Request) {
	pCtx := getPackageContext(r)
	pkg, manifest, err := describe.DescribeInstalledClusterPackage(r.Context(), pCtx.request.manifestName)
	if err != nil && !errors.IsNotFound(err) {
		err = fmt.Errorf("an error occurred fetching package details of %v: %w", pCtx.request.manifestName, err)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	handleConfigurationInput(w, r, &packageContext{
		request:  pCtx.request,
		pkg:      pkg,
		manifest: manifest,
	})
}

func handleConfigurationInput(w http.ResponseWriter, r *http.Request, d *packageContext) {
	ctx := r.Context()

	if d.manifest == nil {
		repoClientset := clicontext.RepoClientsetFromContext(ctx)
		d.manifest = &v1alpha1.PackageManifest{}
		if err := repoClientset.ForRepoWithName(d.request.repositoryName).
			FetchPackageManifest(d.request.manifestName, d.request.version, d.manifest); err != nil {
			responder.SendToast(w,
				toast.WithErr(fmt.Errorf("failed to fetch manifest of %v in version %v: %w",
					d.request.manifestName, d.request.version, err)))
			return
		}
	}

	valueName := r.PathValue("valueName")
	refKind := r.URL.Query().Get("refKind")
	if valueDefinition, ok := d.manifest.ValueDefinitions[valueName]; ok {
		options := components.PkgConfigInputDatalistOptions{}
		if refKind == refKindConfigMap || refKind == refKindSecret {
			if opts, err := opts.GetNamespaceOptions(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get namespace options: %v\n", err)
			} else {
				options.Namespaces = opts
			}
		} else {
			if opts, err := opts.GetPackagesOptions(r.Context()); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get package options: %v\n", err)
			} else {
				options.Names = opts
			}
		}
		input := components.ForPkgConfigInput(
			d.pkg, d.request.repositoryName, d.request.version, d.manifest, valueName, valueDefinition, nil, &options,
			&components.PkgConfigInputRenderOptions{
				Autofocus:      true,
				DesiredRefKind: &refKind,
			})
		responder.SendComponent(w, r, "components/pkg-config-input", responder.RawTemplate(input))
	}
}

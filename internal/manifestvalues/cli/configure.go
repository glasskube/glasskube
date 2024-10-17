package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/maputils"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	goldmarkutil "github.com/yuin/goldmark/util"
)

var (
	referenceValueStr   = "Reference value"
	literalValueStr     = "Literal value"
	valueKinds          = []string{referenceValueStr, literalValueStr}
	configMapStr        = "ConfigMap"
	secretStr           = "Secret"
	packageStr          = "Package"
	referenceValueKinds = []string{configMapStr, secretStr, packageStr}
)

var (
	bold  = color.New(color.Bold).SprintfFunc()
	faint = color.New(color.Faint).SprintfFunc()
	green = color.GreenString
	red   = color.RedString
)

type ConfigureOptions struct {
	oldValues map[string]v1alpha1.ValueConfiguration
	UseDefaultValuesOption
}

type ConfigureOption func(*ConfigureOptions)

func WithOldValues(oldValues map[string]v1alpha1.ValueConfiguration) ConfigureOption {
	return func(co *ConfigureOptions) { co.oldValues = oldValues }
}

func WithUseDefaults(opts UseDefaultValuesOption) ConfigureOption {
	return func(co *ConfigureOptions) { co.UseDefaultValuesOption = opts }
}

func Configure(
	manifest v1alpha1.PackageManifest,
	opts ...ConfigureOption,
) (map[string]v1alpha1.ValueConfiguration, error) {
	var options ConfigureOptions
	for _, fn := range opts {
		fn(&options)
	}

	newValues := make(map[string]v1alpha1.ValueConfiguration, len(manifest.ValueDefinitions))
	if len(manifest.ValueDefinitions) > 0 {
		fmt.Fprintf(os.Stderr, "\n%v has %v values for configuration.\n\n",
			manifest.Name, len(manifest.ValueDefinitions))
	}

	// TODO: Preserve the order of value definitions set by the author.
	//  This is currently not possible because the kubernetes-sigs/yaml package
	//  converts everything to an interface{} before converting to the target type,
	//  so even if we would use an alternative map implementation that preserves the
	//  order of keys, they would still be different from the original.
	//  Related issue: https://github.com/kubernetes-sigs/yaml/issues/88
	for i, name := range maputils.KeysSorted(manifest.ValueDefinitions) {
		def := manifest.ValueDefinitions[name]
		var oldValuePtr *v1alpha1.ValueConfiguration
		if oldValue, ok := options.oldValues[name]; ok {
			oldValuePtr = &oldValue
		}
		if options.ShouldUseDefault(name, def) {
			fmt.Fprintf(os.Stderr, "Using default value for %v: %v\n", name, def.DefaultValue)
			newValues[name] = v1alpha1.ValueConfiguration{
				InlineValueConfiguration: v1alpha1.InlineValueConfiguration{Value: util.Pointer(def.DefaultValue)},
			}
		} else {
			if newValue, err := ConfigureSingle(name, def, oldValuePtr); err != nil {
				return nil, err
			} else if newValue != nil {
				newValues[name] = *newValue
			}
		}
		fmt.Fprintf(os.Stderr, "\nProgress: %v%v\n\n",
			green(strings.Repeat("✔", i+1)),
			faint(strings.Repeat("·", len(manifest.ValueDefinitions)-(i+1))),
		)
	}
	return newValues, nil
}

func ConfigureSingle(
	name string,
	def v1alpha1.ValueDefinition,
	oldValue *v1alpha1.ValueConfiguration,
) (*v1alpha1.ValueConfiguration, error) {
	for {
		printHeader(name, def)

		if oldValue != nil {
			fmt.Fprintln(os.Stderr, "Old value:", manifestvalues.ValueAsString(*oldValue))
			if cliutils.YesNoPrompt("Keep?", true) {
				return oldValue, nil
			}
		}

		var newValue v1alpha1.ValueConfiguration
		var err error

		useDefault := len(def.DefaultValue) > 0
		if useDefault {
			fmt.Fprintln(os.Stderr, "Default:", def.DefaultValue)
			useDefault = cliutils.YesNoPrompt("Use default?", true)
		}

		if useDefault {
			newValue.Value = &def.DefaultValue
		} else {
			fmt.Fprintln(os.Stderr, "Would you like to specify a reference value (ConfigMap, Secret, Package) or literal value?")
			var opt string
			if opt, err = getOptionWithDefault(valueKinds, &literalValueStr); err == nil {
				switch opt {
				case referenceValueStr:
					newValue.ValueFrom, err = getReferenceValue()
				case literalValueStr:
					newValue.Value, err = getLiteralValue(name, def)
				default:
					err = fmt.Errorf("invalid option: %v", opt)
				}
			}
		}

		// Skip validation if we have a reference value.
		// They should be resolved all at once at a later time.
		if err == nil && newValue.Value != nil {
			err = manifestvalues.ValidateSingle(name, def, *newValue.Value)
		}

		if err == nil {
			return &newValue, nil
		}

		fmt.Fprintln(os.Stderr, red("Invalid input:"), err)
		if !cliutils.YesNoPrompt("Try again?", true) {
			return nil, err
		}
	}
}

func printHeader(name string, def v1alpha1.ValueDefinition) {
	title := name
	if len(def.Metadata.Label) > 0 {
		title = def.Metadata.Label
	}
	fmt.Fprintln(os.Stderr, bold(title))
	if len(def.Metadata.Description) > 0 {
		printMarkdown(os.Stderr, def.Metadata.Description)
	}
}

func printMarkdown(w io.Writer, text string) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Linkify,
		),
		goldmark.WithRenderer(renderer.NewRenderer(
			renderer.WithNodeRenderers(
				goldmarkutil.Prioritized(cliutils.MarkdownRenderer(), 1000),
			),
		)),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		util.Must(fmt.Fprintln(w, text))
	} else {
		util.Must(fmt.Fprint(w, strings.TrimSpace(buf.String())+"\n\n"))
	}
}

func getLiteralValue(name string, def v1alpha1.ValueDefinition) (*string, error) {
	switch def.Type {
	case v1alpha1.ValueTypeOptions:
		if len(def.Options) == 0 {
			// retry makes no sense in this case, we can return an error
			return nil, fmt.Errorf("%v has no options", name)
		}
		if v, err := getOption(def.Options); err != nil {
			return nil, err
		} else {
			return &v, nil
		}
	default:
		fmt.Fprintln(os.Stderr, "Please enter a value:")
		v := getInput(def.Type)
		return &v, nil
	}
}

func getReferenceValue() (*v1alpha1.ValueReference, error) {
	if opt, err := getOption(referenceValueKinds); err != nil {
		return nil, err
	} else {
		switch opt {
		case configMapStr:
			fmt.Fprintln(os.Stderr, "Specify the namespace and name and key of the ConfigMap data")
			return &v1alpha1.ValueReference{ConfigMapRef: getObjectKeyValueSource()}, nil
		case secretStr:
			fmt.Fprintln(os.Stderr, "Specify the namespace and name and key of the Secret data")
			return &v1alpha1.ValueReference{SecretRef: getObjectKeyValueSource()}, nil
		case packageStr:
			fmt.Fprintln(os.Stderr, "Specify the name and value of the Package")
			return &v1alpha1.ValueReference{PackageRef: getPackageValueSource()}, nil
		default:
			return nil, fmt.Errorf("invalid option: %v (this is a bug)", opt)
		}
	}
}

func getObjectKeyValueSource() *v1alpha1.ObjectKeyValueSource {
	var ref v1alpha1.ObjectKeyValueSource
	ref.Namespace = cliutils.GetInputStr("namespace")
	ref.Name = cliutils.GetInputStr("name")
	ref.Key = cliutils.GetInputStr("key")
	return &ref
}

func getPackageValueSource() *v1alpha1.PackageValueSource {
	var ref v1alpha1.PackageValueSource
	ref.Name = cliutils.GetInputStr("name")
	ref.Value = cliutils.GetInputStr("value")
	return &ref
}

func getInput(t v1alpha1.ValueType) string {
	return cliutils.GetInputStr(string(t))
}

func getOption(options []string) (string, error) {
	return cliutils.GetOption(string(v1alpha1.ValueTypeOptions), options)
}

func getOptionWithDefault(options []string, defaultOption *string) (string, error) {
	return cliutils.GetOptionWithDefault(string(v1alpha1.ValueTypeOptions), options, defaultOption)
}

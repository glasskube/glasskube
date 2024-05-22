package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"sigs.k8s.io/yaml"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/spf13/cobra"
)

type DescribeOutput struct {
	Package         string              `json:"package"`
	Version         string              `json:"version"`
	Status          string              `json:"status"`
	AutoUpdate      string              `json:"autoUpdate,omitempty"`
	Entrypoints     []map[string]string `json:"entrypoints,omitempty"`
	Dependencies    []map[string]string `json:"dependencies,omitempty"`
	References      []map[string]string `json:"references,omitempty"`
	LongDescription string              `json:"longDescription,omitempty"`
	Configuration   map[string]string   `json:"configuration,omitempty"`
}

type DescribeFormat string

func (o *DescribeFormat) String() string {
	return string(*o)
}

func (o *DescribeFormat) Set(value string) error {
	switch value {
	case string(JSON), string(YAML):
		*o = DescribeFormat(value)
		return nil
	default:
		return errors.New(`invalid output format, must be "json" or "yaml"`)
	}
}

func (o *DescribeFormat) Type() string {
	return "string"
}

type DescribeOpt struct {
	DescribeFormat DescribeFormat
}

var describeCmdOpt = DescribeOpt{}

var describeCmd = &cobra.Command{
	Use:               "describe [package-name]",
	Short:             "Describe a package",
	Long:              "Shows additional information about the given package.",
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		pkg, pkgStatus, manifest, latestVersion, err := describe.DescribePackage(cmd.Context(), pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not describe package %v: %v\n", pkgName, err)
			cliutils.ExitWithError()
		}

		if describeCmdOpt.DescribeFormat == DescribeFormat(JSON) {
			printJson(pkg, pkgStatus, manifest, latestVersion)
		} else if describeCmdOpt.DescribeFormat == DescribeFormat(YAML) {
			printYaml(pkg, pkgStatus, manifest, latestVersion)
		} else {

			bold := color.New(color.Bold).SprintFunc()

			fmt.Println(bold("Package:"), nameAndDescription(manifest))
			fmt.Println(bold("Version:"), version(pkg, latestVersion))
			fmt.Println(bold("Status: "), status(pkgStatus))
			if pkg != nil {
				fmt.Println(bold("Auto-Update:"), clientutils.AutoUpdateString(pkg, "Disabled"))
			}

			if len(manifest.Entrypoints) > 0 {
				fmt.Println()
				fmt.Println(bold("Entrypoints:"))
				printEntrypoints(manifest)
			}

			if len(manifest.Dependencies) > 0 {
				fmt.Println()
				fmt.Println(bold("Dependencies:"))
				printDependencies(manifest)
			}

			fmt.Println()
			fmt.Printf("%v \n", bold("References:"))
			printReferences(pkg, manifest, latestVersion)

			trimmedDescription := strings.TrimSpace(manifest.LongDescription)
			if len(trimmedDescription) > 0 {
				fmt.Println()
				fmt.Println(bold("Long Description:"))
				printMarkdown(os.Stdout, trimmedDescription)
			}

			if pkg != nil && len(pkg.Spec.Values) > 0 {
				fmt.Println()
				fmt.Println(bold("Configuration:"))
				printValueConfigurations(os.Stdout, pkg.Spec.Values)
			}
		}
	},
}

func printEntrypoints(manifest *v1alpha1.PackageManifest) {
	for _, i := range manifest.Entrypoints {
		var messageParts []string
		if i.Name != "" {
			messageParts = append(messageParts, fmt.Sprint("Name: ", i.Name))
		}
		if i.ServiceName != "" {
			messageParts = append(messageParts, fmt.Sprintf("Remote: %s:%v", i.ServiceName, i.Port))
		}
		var localUrl string
		if i.Scheme != "" {
			localUrl += i.Scheme + "://localhost:"
		} else {
			localUrl += "http://localhost:"
		}
		if i.LocalPort != 0 {
			localUrl += fmt.Sprint(i.LocalPort)
		} else {
			localUrl += fmt.Sprint(i.Port)
		}
		localUrl += "/"
		messageParts = append(messageParts, fmt.Sprint("Local: ", localUrl))
		entrypointMsg := strings.Join(messageParts, ", ")
		fmt.Fprintf(os.Stderr, " * %s\n", entrypointMsg)
	}
}

func printDependencies(manifest *v1alpha1.PackageManifest) {
	for _, dep := range manifest.Dependencies {
		fmt.Printf(" * %v", dep.Name)
		if len(dep.Version) > 0 {
			fmt.Printf(" (%v)", dep.Version)
		}
		fmt.Println()
	}
}

func printReferences(pkg *v1alpha1.Package, manifest *v1alpha1.PackageManifest, latestVersion string) {
	version := latestVersion
	if pkg != nil {
		version = pkg.Spec.PackageInfo.Version
	}
	if url, err := repo.GetPackageManifestURL("", manifest.Name, version); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not get package manifest url: %v\n", err)
	} else {
		fmt.Printf(" * Glasskube Package Manifest: %v\n", url)
	}
	for _, ref := range manifest.References {
		fmt.Printf(" * %v: %v\n", ref.Label, ref.Url)
	}
}

func printValueConfigurations(w io.Writer, values map[string]v1alpha1.ValueConfiguration) {
	for name, value := range values {
		fmt.Fprintf(w, " * %v: %v\n", name, manifestvalues.ValueAsString(value))
	}
}

func printMarkdown(w io.Writer, text string) {
	md := goldmark.New(
		goldmark.WithRenderer(renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(cliutils.MarkdownRenderer(), 1000),
			),
		)),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		fmt.Fprintln(w, text)
	} else {
		fmt.Fprintln(w, strings.TrimSpace(buf.String()))
	}
}

func status(pkgStatus *client.PackageStatus) string {
	if pkgStatus != nil {
		switch pkgStatus.Status {
		case string(condition.Ready):
			return pkgStatus.Status
		case string(condition.Failed):
			return pkgStatus.Status
		default:
			return pkgStatus.Status
		}
	} else {
		return "Not installed"
	}
}

func version(pkg *v1alpha1.Package, latestVersion string) string {
	if len(latestVersion) == 0 {
		if pkg.Spec.PackageInfo.Version != pkg.Status.Version {
			return fmt.Sprintf("%v (desired: %v)", pkg.Status.Version, pkg.Spec.PackageInfo.Version)
		} else {
			return pkg.Status.Version
		}
	} else {
		return latestVersion
	}
}

func nameAndDescription(manifest *v1alpha1.PackageManifest) string {
	var bld strings.Builder
	_, _ = bld.WriteString(manifest.Name)
	trimmedDescription := strings.TrimSpace(manifest.ShortDescription)
	if len(trimmedDescription) > 0 {
		_, _ = bld.WriteString(" — ")
		_, _ = bld.WriteString(trimmedDescription)
	}
	return bld.String()
}

func getEntrypoints(manifest *v1alpha1.PackageManifest) []map[string]string {
	entrypoints := []map[string]string{}
	for _, i := range manifest.Entrypoints {
		entrypoint := make(map[string]string)
		if i.Name != "" {
			entrypoint["name"] = i.Name
		}
		if i.ServiceName != "" {
			entrypoint["remote"] = fmt.Sprintf("%s:%v", i.ServiceName, i.Port)
		}
		var localUrl string
		if i.Scheme != "" {
			localUrl += i.Scheme + "://localhost:"
		} else {
			localUrl += "http://localhost:"
		}
		if i.LocalPort != 0 {
			localUrl += fmt.Sprint(i.LocalPort)
		} else {
			localUrl += fmt.Sprint(i.Port)
		}
		localUrl += "/"
		entrypoint["local"] = localUrl
		entrypoints = append(entrypoints, entrypoint)
	}
	return entrypoints
}

func getDependencies(manifest *v1alpha1.PackageManifest) []map[string]string {
	dependencies := []map[string]string{}
	for _, dep := range manifest.Dependencies {
		dependency := make(map[string]string)
		dependency["name"] = dep.Name
		if len(dep.Version) > 0 {
			dependency["version"] = dep.Version
		}
		dependencies = append(dependencies, dependency)
	}
	return dependencies
}

func getReferences(manifest *v1alpha1.PackageManifest, latestVersion string) []map[string]string {
	references := []map[string]string{}
	if url, err := repo.GetPackageManifestURL("", manifest.Name, latestVersion); err == nil {
		references = append(references, map[string]string{
			"label": "Glasskube Package Manifest",
			"url":   url,
		})
	}
	for _, ref := range manifest.References {
		references = append(references, map[string]string{
			"label": ref.Label,
			"url":   ref.Url,
		})
	}
	return references
}

func getConfigurations(pkg *v1alpha1.Package) map[string]string {
	configuration := make(map[string]string)
	if pkg != nil && len(pkg.Spec.Values) > 0 {
		for name, value := range pkg.Spec.Values {
			configuration[name] = manifestvalues.ValueAsString(value)
		}
	}
	return configuration
}

func printJson(pkg *v1alpha1.Package, pkgStatus *client.PackageStatus, manifest *v1alpha1.PackageManifest, latestVersion string) {
	output := DescribeOutput{
		Package:         nameAndDescription(manifest),
		Version:         version(pkg, latestVersion),
		Status:          status(pkgStatus),
		AutoUpdate:      clientutils.AutoUpdateString(pkg, "Disabled"),
		Entrypoints:     getEntrypoints(manifest),
		Dependencies:    getDependencies(manifest),
		References:      getReferences(manifest, latestVersion),
		LongDescription: manifest.LongDescription,
		Configuration:   getConfigurations(pkg),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	err := enc.Encode(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling data to JSON: %v\n", err)
		cliutils.ExitWithError()
	}
}

func printYaml(pkg *v1alpha1.Package, pkgStatus *client.PackageStatus, manifest *v1alpha1.PackageManifest, latestVersion string) {
	output := DescribeOutput{
		Package:         nameAndDescription(manifest),
		Version:         version(pkg, latestVersion),
		Status:          status(pkgStatus),
		AutoUpdate:      clientutils.AutoUpdateString(pkg, "Disabled"),
		Entrypoints:     getEntrypoints(manifest),
		Dependencies:    getDependencies(manifest),
		References:      getReferences(manifest, latestVersion),
		LongDescription: manifest.LongDescription,
		Configuration:   getConfigurations(pkg),
	}

	data, err := yaml.Marshal(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling data to YAML: %v\n", err)
		cliutils.ExitWithError()
	}
	fmt.Println(string(data))
}

func init() {
	RootCmd.AddCommand(describeCmd)
	describeCmd.PersistentFlags().VarP((&describeCmdOpt.DescribeFormat), "output", "o", "output format (json, yaml, etc.)")
}

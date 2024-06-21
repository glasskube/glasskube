package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/manifestvalues"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/describe"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/yaml"
)

var describeCmdOptions = struct {
	repository string
	OutputOptions
}{}

var describeCmd = &cobra.Command{
	Use:               "describe [package-name]",
	Short:             "Describe a package",
	Long:              "Shows additional information about the given package.",
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		pkgName := args[0]
		repoClient := cliutils.RepositoryClientset(ctx)

		latestManifest, latestVersion, err :=
			describe.DescribeLatestVersion(ctx, describeCmdOptions.repository, pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not get latest info for %v: %v\n", pkgName, err)
			cliutils.ExitWithError()
		}

		repos, err := repoClient.Meta().GetReposForPackage(pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not get repos for %v: %v\n", pkgName, err)
			cliutils.ExitWithError()
		}

		pkg, manifest, err := describe.DescribeInstalledPackage(ctx, pkgName)
		if err != nil && !apierrors.IsNotFound(err) {
			// Unhandled error -> exit
			fmt.Fprintf(os.Stderr, "❌ Could not describe package %v: %v\n", pkgName, err)
			cliutils.ExitWithError()
		} else if err != nil {
			// package not installed -> use latest manifest from repo
			manifest = latestManifest
		}

		bold := color.New(color.Bold).SprintFunc()

		if describeCmdOptions.Output == OutputFormatJSON {
			printJSON(ctx, pkg, manifest, latestVersion, repos)
		} else if describeCmdOptions.Output == OutputFormatYAML {
			printYAML(ctx, pkg, manifest, latestVersion, repos)
		} else {
			fmt.Println(bold("Package:"), nameAndDescription(manifest))
			fmt.Println(bold("Version:"), version(pkg, latestVersion))
			fmt.Println(bold("Status: "), status(pkg))
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
			fmt.Println(bold("Package repositories:"))
			printRepositories(pkg, repos)

			fmt.Println()
			fmt.Printf("%v \n", bold("References:"))
			printReferences(ctx, pkg, manifest)

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

func printRepositories(pkg *v1alpha1.ClusterPackage, repos []v1alpha1.PackageRepository) {
	for _, repo := range repos {
		fmt.Fprintf(os.Stderr, " * %v", repo.Name)
		if isInstalledFrom(pkg, repo) {
			fmt.Fprintln(os.Stderr, " (installed)")
		} else {
			fmt.Fprintln(os.Stderr)
		}
	}
}

func repositoriesAsMap(pkg *v1alpha1.ClusterPackage, repos []v1alpha1.PackageRepository) []map[string]any {
	repositories := make([]map[string]any, 0, len(repos))
	for _, repo := range repos {
		repositories = append(repositories, map[string]any{
			"name":      repo.Name,
			"installed": isInstalledFrom(pkg, repo),
		})
	}
	return repositories
}

func isInstalledFrom(pkg *v1alpha1.ClusterPackage, repo v1alpha1.PackageRepository) bool {
	return pkg != nil &&
		(repo.Name == pkg.Spec.PackageInfo.RepositoryName ||
			(len(pkg.Spec.PackageInfo.RepositoryName) == 0 && repo.IsDefaultRepository()))
}

func printReferences(ctx context.Context, pkg *v1alpha1.ClusterPackage, manifest *v1alpha1.PackageManifest) {
	repo := cliutils.RepositoryClientset(ctx)
	var repoClient repoclient.RepoClient
	if pkg != nil {
		repoClient = repo.ForPackage(pkg)
		if url, err := repoClient.GetPackageManifestURL(manifest.Name, pkg.Spec.PackageInfo.Version); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not get package manifest url: %v\n", err)
		} else {
			fmt.Printf(" * Glasskube Package Manifest: %v\n", url)
		}
	}
	for _, ref := range manifest.References {
		fmt.Printf(" * %v: %v\n", ref.Label, ref.Url)
	}
}

func referencesAsMap(
	ctx context.Context,
	pkg *v1alpha1.ClusterPackage,
	manifest *v1alpha1.PackageManifest,
) []map[string]string {
	references := []map[string]string{}
	for _, ref := range manifest.References {
		reference := make(map[string]string)
		reference["label"] = ref.Label
		reference["url"] = ref.Url
		references = append(references, reference)
	}
	if pkg != nil {
		repo := cliutils.RepositoryClientset(ctx)
		repoClient := repo.ForPackage(pkg)
		if url, err := repoClient.GetPackageManifestURL(manifest.Name, pkg.Spec.PackageInfo.Version); err == nil {
			reference := make(map[string]string)
			reference["label"] = "Glasskube Package Manifest"
			reference["url"] = url
			references = append(references, reference)
		}
	}
	return references
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

func status(pkg *v1alpha1.ClusterPackage) string {
	pkgStatus := client.GetStatusOrPending(pkg)
	if pkgStatus != nil {
		switch pkgStatus.Status {
		case string(condition.Ready):
			return color.GreenString(pkgStatus.Status)
		case string(condition.Failed):
			return color.RedString(pkgStatus.Status)
		default:
			return pkgStatus.Status
		}
	} else {
		return color.New(color.Faint).Sprint("Not installed")
	}
}

func version(pkg *v1alpha1.ClusterPackage, latestVersion string) string {
	if pkg != nil {
		var parts []string
		if len(pkg.Status.Version) > 0 {
			parts = append(parts, pkg.Status.Version)
		}
		if pkg.Spec.PackageInfo.Version != pkg.Status.Version {
			parts = append(parts, fmt.Sprintf("(desired: %v)", pkg.Spec.PackageInfo.Version))
		}
		if semver.IsUpgradable(pkg.Spec.PackageInfo.Version, latestVersion) {
			parts = append(parts, fmt.Sprintf("(latest: %v)", latestVersion))
		}
		return strings.Join(parts, " ")
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

func createOutputStructure(
	ctx context.Context,
	pkg *v1alpha1.ClusterPackage,
	manifest *v1alpha1.PackageManifest,
	latestVersion string,
	repos []v1alpha1.PackageRepository,
) map[string]interface{} {
	data := map[string]interface{}{
		"packageName":      manifest.Name,
		"shortDescription": manifest.ShortDescription,
		"latestVersion":    latestVersion,
		"status":           "Not Installed",
		"entrypoints":      manifest.Entrypoints,
		"dependencies":     manifest.Dependencies,
		"longDescription":  strings.TrimSpace(manifest.LongDescription),
		"repositories":     repositoriesAsMap(pkg, repos),
		"references":       referencesAsMap(ctx, pkg, manifest),
	}
	if pkg != nil {
		data["desiredVersion"] = pkg.Spec.PackageInfo.Version
		data["configuration"] = pkg.Spec.Values
		data["version"] = pkg.Status.Version
		data["autoUpdate"] = pkg.AutoUpdatesEnabled()
		data["isUpgradable"] = semver.IsUpgradable(pkg.Spec.PackageInfo.Version, latestVersion)
		data["status"] = client.GetStatusOrPending(pkg).Status
	}
	return data
}

func printJSON(ctx context.Context,
	pkg *v1alpha1.ClusterPackage,
	manifest *v1alpha1.PackageManifest,
	latestVersion string,
	repos []v1alpha1.PackageRepository) {
	output := createOutputStructure(ctx, pkg, manifest, latestVersion, repos)
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not marshal JSON output: %v\n", err)
		cliutils.ExitWithError()
	}
	fmt.Println(string(jsonOutput))
}

func printYAML(ctx context.Context,
	pkg *v1alpha1.ClusterPackage,
	manifest *v1alpha1.PackageManifest,
	latestVersion string,
	repos []v1alpha1.PackageRepository) {
	output := createOutputStructure(ctx, pkg, manifest, latestVersion, repos)
	yamlOutput, err := yaml.Marshal(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not marshal YAML output: %v\n", err)
		cliutils.ExitWithError()
	}
	fmt.Println(string(yamlOutput))
}

func init() {
	describeCmd.Flags().StringVar(&describeCmdOptions.repository, "repository", describeCmdOptions.repository,
		"specify the name of the package repository used to use when the package is not installed")
	RootCmd.AddCommand(describeCmd)
	describeCmdOptions.OutputOptions.AddFlagsToCommand(describeCmd)
}

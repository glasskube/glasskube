package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/releaseinfo"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type bootstrapOptions struct {
	url                     string
	bootstrapType           bootstrap.BootstrapType
	latest                  bool
	disableTelemetry        bool
	force                   bool
	createDefaultRepository bool
	yes                     bool
	dryRun                  bool
	OutputOptions
}

var bootstrapCmdOptions = bootstrapOptions{
	bootstrapType:           bootstrap.BootstrapTypeAio,
	createDefaultRepository: true,
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap Glasskube in a Kubernetes cluster",
	Long: "Bootstraps Glasskube in a Kubernetes cluster, " +
		"thereby installing the Glasskube operator and checking if the installation was successful.",
	Args:   cobra.NoArgs,
	PreRun: cliutils.SetupClientContext(false, util.Pointer(true)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)
		client := bootstrap.NewBootstrapClient(cfg)
		ctx := cmd.Context()

		var installedVersion, targetVersion *semver.Version
		if installedVersionRaw, err := clientutils.GetPackageOperatorVersion(ctx); err != nil {
			if !apierrors.IsNotFound(err) {
				fmt.Fprintf(os.Stderr, "could not determine installed version: %v\n", err)
				cliutils.ExitWithError()
			}
		} else if installedVersion, err = semver.NewVersion(installedVersionRaw); err != nil {
			fmt.Fprintf(os.Stderr, "could not parse installed version: %v\n", err)
			cliutils.ExitWithError()
		}
		if bootstrapCmdOptions.url == "" {
			version := config.Version
			if bootstrapCmdOptions.latest {
				if releaseInfo, err := releaseinfo.FetchLatestRelease(); err != nil {
					fmt.Fprintf(os.Stderr, "could not determine latest version: %v\n", err)
					cliutils.ExitWithError()
				} else {
					version = releaseInfo.Version
				}
			}
			var err error
			if targetVersion, err = semver.NewVersion(version); err != nil {
				fmt.Fprintf(os.Stderr, "could not parse target version: %v\n", err)
				cliutils.ExitWithError()
			}
		}

		if bootstrapCmdOptions.dryRun {
			fmt.Fprintln(os.Stderr,
				"🔎 Dry-run mode is enabled. "+
					"Nothing will be changed in your cluster, but validations will still be run.")
		}

		verifyLegalUpdate(ctx, installedVersion, targetVersion)

		manifests, err := client.Bootstrap(ctx, bootstrapCmdOptions.asBootstrapOptions())
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			cliutils.ExitWithError()
		}
		if err := printBootstrap(
			manifests,
			bootstrapCmdOptions.Output,
		); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred in printing : %v\n", err)
			cliutils.ExitWithError()
		}
	},
}

func (o bootstrapOptions) asBootstrapOptions() bootstrap.BootstrapOptions {
	return bootstrap.BootstrapOptions{
		Type:                    o.bootstrapType,
		Url:                     o.url,
		Latest:                  o.latest,
		DisableTelemetry:        o.disableTelemetry,
		Force:                   o.force,
		CreateDefaultRepository: o.createDefaultRepository,
		DryRun:                  o.dryRun,
	}
}

func printBootstrap(manifests []unstructured.Unstructured, output OutputFormat) error {
	if output != "" {
		if err := convertAndPrintManifests(manifests, output); err != nil {
			return err
		}
	}
	return nil
}

func verifyLegalUpdate(ctx context.Context, installedVersion, targetVersion *semver.Version) {
	breakings := map[*semver.Version]string{
		semver.New(0, 10, 0, "", ""): "In release v0.10.0, Packages are renamed to ClusterPackages and " +
			"Packages are reintroduced as namespaced resources.\n" +
			"Glasskube must be uninstalled completely, to perform this update.",
	}
	currentContext := color.New(color.Bold).Sprint(clicontext.RawConfigFromContext(ctx).CurrentContext)

	if installedVersion != nil && targetVersion != nil {
		for version, msg := range breakings {
			if installedVersion.LessThan(version) && !targetVersion.LessThan(version) {
				fmt.Fprintf(os.Stderr,
					"❗ Upgrade from version v%v to v%v is not possible due to a breaking change in v%v\n\n"+
						"Details: %v\n\n"+
						"For more information, please refer to our documentation: "+
						"https://glasskube.dev/docs/getting-started/upgrading/#%v\n",
					installedVersion, targetVersion, version, msg, version)
				cliutils.ExitWithError()
			}
		}
		if installedVersion.GreaterThan(targetVersion) {
			fmt.Fprintf(os.Stderr,
				"⚠️  Glasskube is already installed in this cluster in the newer version v%v. "+
					"You are about to install version v%v. This could lead to a broken cluster!\n",
				installedVersion, targetVersion)
			if !bootstrapCmdOptions.yes &&
				!cliutils.YesNoPrompt("Are you sure that you want to downgrade glasskube in this cluster?", false) {
				cancel()
			}
		} else if installedVersion.LessThan(targetVersion) {
			fmt.Fprintf(os.Stderr, "Glasskube will be updated to version v%v in cluster %v.\n",
				targetVersion, currentContext)
			if !bootstrapCmdOptions.yes && !cliutils.YesNoPrompt("Continue?", true) {
				cancel()
			}
		} else {
			fmt.Fprintf(os.Stderr,
				"⚠️  Glasskube is already installed in this cluster (%v) in version v%v. "+
					"You are about to bootstrap this version again.\n",
				currentContext, installedVersion)
			if !bootstrapCmdOptions.yes && !cliutils.YesNoPrompt("Continue?", true) {
				cancel()
			}
		}
	} else if installedVersion != nil && targetVersion == nil {
		fmt.Fprintf(os.Stderr,
			"⚠️  Glasskube is currently installed in this cluster (%v) in version v%v. "+
				"The version you are about to install is unknown. "+
				"Please make sure the versions are compatible, this action could lead to a broken "+
				"cluster!\n",
			currentContext, installedVersion)
		if !bootstrapCmdOptions.yes && !cliutils.YesNoPrompt("Continue?", false) {
			cancel()
		}
	} else {
		fmt.Fprintf(os.Stderr, "Glasskube will be installed in context %s.\n", currentContext)
		if !bootstrapCmdOptions.yes && !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}
	}
}

func convertAndPrintManifests(
	objs []unstructured.Unstructured,
	output OutputFormat,
) error {
	switch output {
	case OutputFormatJSON:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		err := enc.Encode(objs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error marshaling data to JSON: %v\n", err)
			cliutils.ExitWithError()
		}
	case OutputFormatYAML:
		for i, obj := range objs {
			yamlData, err := yaml.Marshal(&obj)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error marshaling data to YAML: %v\n", err)
				cliutils.ExitWithError()
			}

			if i > 0 {
				fmt.Println("---")
			}

			fmt.Println(string(yamlData))
		}
	default:
		return fmt.Errorf("unsupported output format: %v", output)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().StringVarP(&bootstrapCmdOptions.url, "url", "u", "", "URL to fetch the Glasskube operator from")
	bootstrapCmd.Flags().VarP(&bootstrapCmdOptions.bootstrapType, "type", "t", `Type of manifest to use for bootstrapping`)
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.latest, "latest", config.IsDevBuild(),
		"Fetch and bootstrap the latest version")
	bootstrapCmdOptions.OutputOptions.AddFlagsToCommand(bootstrapCmd)
	bootstrapCmd.Flags().BoolVarP(&bootstrapCmdOptions.force, "force", "f", bootstrapCmdOptions.force,
		"Do not bail out if pre-checks fail")
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.disableTelemetry, "disable-telemetry", false, "Disable telemetry")
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.createDefaultRepository, "create-default-repository",
		bootstrapCmdOptions.createDefaultRepository,
		"Toggle creation of the default glasskube package repository")
	bootstrapCmd.PersistentFlags().BoolVar(&bootstrapCmdOptions.dryRun, "dry-run", false,
		"Do not make any changes but run all validations")
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.yes, "yes", false, "Skip confirmation prompt")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "type")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "latest")
}

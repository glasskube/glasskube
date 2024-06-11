package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	js "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type bootstrapOptions struct {
	url                     string
	bootstrapType           bootstrap.BootstrapType
	latest                  bool
	disableTelemetry        bool
	force                   bool
	createDefaultRepository bool
	yes                     bool
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

		currentContext := clicontext.RawConfigFromContext(ctx).CurrentContext

		if !bootstrapCmdOptions.yes {
			confirmMessage := fmt.Sprintf("Glasskube will be installed in context %s.\nContinue? ", currentContext)
			if !cliutils.YesNoPrompt(confirmMessage, true) {
				fmt.Println("Operation stopped")
				cliutils.ExitWithError()
			}
		}

		installedVersion, err := clientutils.GetPackageOperatorVersion(cmd.Context())
		if err != nil {
			IsBootstrapped, err := bootstrap.IsBootstrapped(cmd.Context(), cfg)
			if err != nil && !IsBootstrapped {
				fmt.Printf("error : %v\n", err)
				cliutils.ExitWithError()
			}
		}

		var desiredVersion string
		if bootstrapCmdOptions.url == "" {
			desiredVersion = config.Version
		} else {
			desiredVersion = ""
		}
		if !semver.IsUpgradable(installedVersion, desiredVersion) &&
			installedVersion != "" &&
			installedVersion[1:] != desiredVersion {
			if !cliutils.YesNoPrompt(fmt.Sprintf("Glasskube is already installed in this cluster "+
				"in the newer version %v. You are about to install version %v. This could lead "+
				"to a broken cluster!\nAre you sure that you want to downgrade glasskube "+
				"in this cluster?", installedVersion, desiredVersion), false) {
				fmt.Println("Operation stopped")
				cliutils.ExitWithError()
			}

		}
		manifests, err := client.Bootstrap(
			cmd.Context(),
			bootstrapCmdOptions.asBootstrapOptions(),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			cliutils.ExitWithError()
		}
		if err := bootstrapCmdOptions.printBootsrap(manifests, bootstrapCmdOptions.Output); err != nil {
			fmt.Fprintf(os.Stderr, "\n Error occured in printing : %v\n", err)
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
	}
}

func (o bootstrapOptions) printBootsrap(manifests []unstructured.Unstructured, output OutputFormat) error {
	if o.Output != "" {
		if err := convertAndPrintManifests(manifests, output); err != nil {
			return err
		}
	}
	return nil
}

func convertAndPrintManifests(
	objs []unstructured.Unstructured,
	output OutputFormat,
) error {
	scheme := runtime.NewScheme()
	var opt bool
	var div string
	switch output {
	case "json":
		opt = false
		div = ","
	case "yaml":
		opt = true
		div = "---"
	}
	serializer := js.NewSerializerWithOptions(
		js.DefaultMetaFactory, scheme, scheme,
		js.SerializerOptions{Yaml: opt, Pretty: true, Strict: true},
	)

	for _, obj := range objs {
		runtimeObj := &unstructured.Unstructured{}
		obj.DeepCopyInto(runtimeObj)

		if err := serializer.Encode(runtimeObj, os.Stdout); err != nil {
			return fmt.Errorf("failed to serialize object %v: %v", obj.GetName(), err)
		}
		fmt.Println(div)
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
	bootstrapCmd.Flags().BoolVar(&bootstrapCmdOptions.yes, "yes", false, "Skip confirmation prompt")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "type")
	bootstrapCmd.MarkFlagsMutuallyExclusive("url", "latest")
}

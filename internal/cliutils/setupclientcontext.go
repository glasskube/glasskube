package cliutils

import (
	"context"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/pkg/bootstrap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/spf13/cobra"
)

func SetupClientContext(requireBootstrapped bool, skipUpdateCheck *bool) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg, rawCfg := RequireConfig(config.Kubeconfig)
		telemetry.InitClient(cfg)
		if requireBootstrapped {
			RequireBootstrapped(cmd.Context(), cfg, rawCfg)
		}
		if ctx, err := clicontext.SetupContext(cmd.Context(), cfg, rawCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up the client:\n\n%v\n", err)
			ExitWithError()
		} else {
			cmd.SetContext(ctx)
		}
		if !*skipUpdateCheck {
			if err := CheckPackageOperatorVersion(cmd.Context()); err != nil {
				fmt.Fprintf(os.Stderr, "Error checking PackageOperator version:\n\n%v\n", err)
			}
		}
	}
}

var bootstrapMessage = `
You're almost there!

Glasskube is not yet installed in your current context %s, but you can do so now.
This will bootstrap Glasskube in your cluster using an all-in-one configuration.
If your use-case requires a slim configuration or custom manifest, please use the "glasskube bootstrap" command.

For further information on bootstrapping, please consult the docs: https://glasskube.dev/docs/getting-started/bootstrap
If you need any help or run into issues, don't hesitate to contact us:
Github: https://github.com/glasskube/glasskube
Discord: https://discord.gg/SxH6KUCGH7

Do you want to install Glasskube in your current context (%s)?`

func RequireBootstrapped(ctx context.Context, cfg *rest.Config, rawCfg *api.Config) {
	ok, err := bootstrap.IsBootstrapped(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error validating Glasskube:\n\n%v\n", err)
		ExitWithError()
	}
	if !ok {
		yes := YesNoPrompt(fmt.Sprintf(bootstrapMessage, rawCfg.CurrentContext, rawCfg.CurrentContext), false)
		if !yes {
			fmt.Fprint(os.Stderr, "Execution cancelled – Glasskube is not yet bootstrapped.\n")
			ExitWithError()
		}
		client := bootstrap.NewBootstrapClient(cfg)
		if _, err := client.Bootstrap(ctx, bootstrap.DefaultOptions()); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			ExitWithError()
		} else {
			fmt.Fprintf(os.Stderr, "\n\nCongrats, Glasskube is all set up! Have fun managing packages!\n\n")
		}
	}
}

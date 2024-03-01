package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/spf13/cobra"
)

var installCmdOptions = struct {
	Version           string
	EnableAutoUpdates bool
	NoWait            bool
}{}

var installCmd = &cobra.Command{
	Use:               "install [package-name]",
	Short:             "Install a package",
	Long:              `Install a package.`,
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		config := client.RawConfigFromContext(ctx)
		client := client.FromContext(ctx)
		packageName := args[0]

		// Instantiate installer
		installer := install.NewInstaller(client).WithStatusWriter(statuswriter.Spinner())

		if installCmdOptions.Version == "" && !installCmdOptions.EnableAutoUpdates {
			fmt.Fprintf(os.Stderr, "Version not specified. The latest version of %v will be installed.\n", packageName)

			if !cliutils.YesNoPrompt("Would you like to enable automatic updates?", true) {
				var packageIndex repo.PackageIndex
				if err := repo.FetchPackageIndex("", packageName, &packageIndex); err != nil {
					fmt.Fprintf(os.Stderr, "❗ Error: Could not fetch package metadata: %v\n", err)
					if !cliutils.YesNoPrompt("Continue anyways? (Automatic updates will be enabled)", false) {
						cancel()
					}
				} else {
					installCmdOptions.Version = packageIndex.LatestVersion
				}
			}
		}

		var msg string
		if installCmdOptions.Version != "" {
			msg = fmt.Sprintf("%v (version %v) will be installed in your current cluster (%v).",
				packageName, installCmdOptions.Version, config.CurrentContext)
		} else {
			msg = fmt.Sprintf("%v will be installed in your current cluster (%v).",
				packageName, config.CurrentContext)
		}

		if !cliutils.YesNoPrompt(fmt.Sprintf("%v Continue?", msg), true) {
			cancel()
		}

		// Non-blocking install if --no-wait is used
		if installCmdOptions.NoWait {
			go func() {
				if err := installer.Install(ctx, packageName, installCmdOptions.Version); err != nil {
					fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Installation of %v started in the background.\n", packageName)
			}()
			return
		}

		// Blocking install
		status, err := installer.InstallBlocking(ctx, packageName, installCmdOptions.Version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
			os.Exit(1)
		}
		if status != nil {
			switch status.Status {
			case string(condition.Ready):
				fmt.Printf("✅ %v is now installed in %v.\n", packageName, config.CurrentContext)
			default:
				fmt.Printf("❌ %v installation has status %v, reason: %v\nMessage: %v\n",
					packageName, status.Status, status.Reason, status.Message)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Installation status unknown - no error and no status have been observed (this is a bug).")
			os.Exit(1)
		}
	},
}

func cancel() {
	fmt.Fprintln(os.Stderr, "❌ Installation cancelled.")
	os.Exit(1)
}

func completeAvailablePackageNames(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var index repo.PackageRepoIndex
	err := repo.FetchPackageRepoIndex("", &index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching package repository index: %v\n", err)
		return nil, cobra.ShellCompDirectiveError
	}
	names := make([]string, 0, len(index.Packages))
	for _, pkg := range index.Packages {
		if toComplete == "" || strings.HasPrefix(pkg.Name, toComplete) {
			names = append(names, pkg.Name)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func completeAvailablePackageVersions(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	packageName := args[0]
	var packageIndex repo.PackageIndex
	if err := repo.FetchPackageIndex("", packageName, &packageIndex); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	versions := make([]string, 0, len(packageIndex.Versions))
	for _, version := range packageIndex.Versions {
		if toComplete == "" || strings.HasPrefix(version.Version, toComplete) {
			versions = append(versions, version.Version)
		}
	}
	return versions, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	installCmd.PersistentFlags().StringVarP(&installCmdOptions.Version, "version", "v", "",
		"install a specific version")
	_ = installCmd.RegisterFlagCompletionFunc("version", completeAvailablePackageVersions)
	installCmd.PersistentFlags().BoolVar(&installCmdOptions.EnableAutoUpdates, "enable-auto-updates", false,
		"enable automatic updates for this package")
	installCmd.MarkFlagsMutuallyExclusive("version", "enable-auto-updates")
	installCmd.PersistentFlags().BoolVar(&installCmdOptions.NoWait, "no-wait", false, "perform non-blocking install")
	RootCmd.AddCommand(installCmd)
}

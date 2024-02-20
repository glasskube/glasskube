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
		client := client.FromContext(ctx)
		packageName := args[0]

		if installCmdOptions.Version == "" && !installCmdOptions.EnableAutoUpdates {
			fmt.Fprintf(os.Stderr, "Version not specified. The latest version of %v will be installed.\n", packageName)

			if !cliutils.YesNoPrompt("Would you like to enable automatic updates?", false) {
				var packageIndex repo.PackageIndex
				if err := repo.FetchPackageIndex("", packageName, &packageIndex); err != nil {
					fmt.Fprintf(os.Stderr, "❗ Error: Could not fetch package metadata: %v\n", err)
					if !cliutils.YesNoPrompt("Continue anyways? (Automatic updates will be enabled)", false) {
						fmt.Fprintln(os.Stderr, "❌ Installation cancelled.")
						os.Exit(1)
					}
				} else {
					installCmdOptions.Version = packageIndex.LatestVersion
				}
			}
		}

		status, err := install.NewInstaller(client).
			WithStatusWriter(statuswriter.Spinner()).
			InstallBlocking(ctx, packageName, installCmdOptions.Version)

		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
			os.Exit(1)
		}
		if status != nil {
			switch (*status).Status {
			case string(condition.Ready):
				fmt.Printf("✅ %v installed successfully.\n", packageName)
			default:
				fmt.Printf("❌ %v installation has status %v, reason: %v\nMessage: %v\n",
					packageName, (*status).Status, (*status).Reason, (*status).Message)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Installation status unknown - no error and no status have been observed (this is a bug).")
			os.Exit(1)
		}
	},
}

func completeAvailablePackageNames(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	var index repo.PackageRepoIndex
	err := repo.FetchPackageRepoIndex("", &index)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}
	names := make([]string, 0, len(index.Packages))
	for _, pkg := range index.Packages {
		if toComplete == "" || strings.HasPrefix(pkg.Name, toComplete) {
			names = append(names, pkg.Name)
		}
	}
	return names, 0
}

func completeAvailablePackageVersions(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	packageName := args[0]
	var packageIndex repo.PackageIndex
	if err := repo.FetchPackageIndex("", packageName, &packageIndex); err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}
	versions := make([]string, 0, len(packageIndex.Versions))
	for _, version := range packageIndex.Versions {
		if toComplete == "" || strings.HasPrefix(version.Version, toComplete) {
			versions = append(versions, version.Version)
		}
	}
	return versions, 0
}

func init() {
	installCmd.PersistentFlags().StringVarP(&installCmdOptions.Version, "version", "v", "",
		"install a specific version")
	_ = installCmd.RegisterFlagCompletionFunc("version", completeAvailablePackageVersions)
	installCmd.PersistentFlags().BoolVar(&installCmdOptions.EnableAutoUpdates, "enable-auto-updates", false,
		"enable automatic updates for this package")
	installCmd.MarkFlagsMutuallyExclusive("version", "enable-auto-updates")
	RootCmd.AddCommand(installCmd)
}

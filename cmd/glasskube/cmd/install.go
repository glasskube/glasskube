package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	clientadapter "github.com/glasskube/glasskube/internal/adapter/goclient"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/dependency"
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
	Yes               bool
}{}

var installCmd = &cobra.Command{
	Use:               "install [package-name]",
	Short:             "Install a package",
	Long:              `Install a package.`,
	Args:              cobra.ExactArgs(1),
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		config := client.RawConfigFromContext(ctx)
		pkgClient := client.FromContext(ctx)
		dm := dependency.NewDependencyManager(clientadapter.NewPackageClientAdapter(pkgClient))
		installer := install.NewInstaller(pkgClient).WithStatusWriter(statuswriter.Spinner())
		bold := color.New(color.Bold).SprintFunc()
		packageName := args[0]

		if installCmdOptions.Version == "" {
			var packageIndex repo.PackageIndex
			if err := repo.FetchPackageIndex("", packageName, &packageIndex); err != nil {
				fmt.Fprintf(os.Stderr, "❗ Error: Could not fetch package metadata: %v\n", err)
				os.Exit(1)
			}
			installCmdOptions.Version = packageIndex.LatestVersion
			fmt.Fprintf(os.Stderr, "Version not specified. The latest version %v of %v will be installed.\n",
				installCmdOptions.Version, packageName)
		}

		installationPlan := []dependency.Requirement{
			{PackageWithVersion: dependency.PackageWithVersion{Name: packageName, Version: installCmdOptions.Version}},
		}

		var manifest v1alpha1.PackageManifest
		if err := repo.FetchPackageManifest("", packageName, installCmdOptions.Version, &manifest); err != nil {
			fmt.Fprintf(os.Stderr, "❗ Error: Could not fetch package manifest: %v\n", err)
			os.Exit(1)
		} else if validationResult, err :=
			dm.Validate(ctx, &manifest, installCmdOptions.Version); err != nil {
			fmt.Fprintf(os.Stderr, "❗ Error: Could not validate dependencies: %v\n", err)
			os.Exit(1)
		} else if len(validationResult.Conflicts) > 0 {
			fmt.Fprintf(os.Stderr, "❗ Error: %v cannot be installed due to conflicts: %v\n",
				packageName, validationResult.Conflicts)
			os.Exit(1)
		} else if len(validationResult.Requirements) > 0 {
			installationPlan = append(installationPlan, validationResult.Requirements...)
		}

		if !installCmdOptions.EnableAutoUpdates && !installCmdOptions.Yes {
			if cliutils.YesNoPrompt("Would you like to enable automatic updates?", false) {
				installCmdOptions.EnableAutoUpdates = true
			}
		}

		fmt.Fprintln(os.Stderr, bold("Summary:"))
		fmt.Fprintf(os.Stderr, " * The following packages will be installed in your cluster (%v):\n", config.CurrentContext)
		for i, p := range installationPlan {
			fmt.Fprintf(os.Stderr, "    %v. %v (version %v)\n", i+1, p.Name, p.Version)
		}
		if installCmdOptions.EnableAutoUpdates {
			fmt.Fprintln(os.Stderr, " * Automatic updates will be", bold("enabled"))

		} else {
			fmt.Fprintln(os.Stderr, " * Automatic updates will be", bold("not enabled"))
		}

		if !installCmdOptions.Yes && !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}

		if installCmdOptions.NoWait {
			if err := installer.Install(
				ctx, packageName, installCmdOptions.Version, nil, installCmdOptions.EnableAutoUpdates); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr,
				"☑️  %v is being installed in the background.\n"+
					"💡 Run \"glasskube describe %v\" to get the current status",
				packageName, packageName)
		} else {
			status, err := installer.InstallBlocking(ctx, packageName, installCmdOptions.Version, nil,
				installCmdOptions.EnableAutoUpdates)
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
	installCmd.PersistentFlags().BoolVar(&installCmdOptions.NoWait, "no-wait", false, "perform non-blocking install")
	installCmd.PersistentFlags().BoolVarP(&installCmdOptions.Yes, "yes", "y", false, "do not ask for any confirmation")
	installCmd.MarkFlagsMutuallyExclusive("version", "enable-auto-updates")
	RootCmd.AddCommand(installCmd)
}

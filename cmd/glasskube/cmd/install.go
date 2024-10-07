package cmd

import (
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/namespaces"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/dependency"
	"github.com/glasskube/glasskube/internal/manifestvalues/cli"
	"github.com/glasskube/glasskube/internal/maputils"
	"github.com/glasskube/glasskube/internal/repo"
	repoclient "github.com/glasskube/glasskube/internal/repo/client"
	repotypes "github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var installCmdOptions = struct {
	cli.ValuesOptions
	Version           string
	Repository        string
	EnableAutoUpdates bool
	NoWait            bool
	Yes               bool
	OutputOptions
	NamespaceOptions
	DryRunOptions
}{
	ValuesOptions: cli.NewOptions(),
}

var installCmd = &cobra.Command{
	Use:               "install <package-name> [<name>]",
	Short:             "Install a package",
	Long:              `Install a package.`,
	Args:              cobra.RangeArgs(1, 2),
	PreRun:            cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	ValidArgsFunction: completeAvailablePackageNames,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		config := clicontext.RawConfigFromContext(ctx)
		pkgClient := clicontext.PackageClientFromContext(ctx)
		dm := cliutils.DependencyManager(ctx)
		valueResolver := cliutils.ValueResolver(ctx)
		repoClientset := cliutils.RepositoryClientset(ctx)
		installer := install.NewInstaller(pkgClient)
		cs := clicontext.KubernetesClientFromContext(ctx)

		opts := metav1.CreateOptions{}
		if installCmdOptions.DryRun {
			opts.DryRun = []string{metav1.DryRunAll}
			fmt.Fprintln(os.Stderr,
				"🔎 Dry-run mode is enabled. Nothing will be changed.")
		}

		if !rootCmdOptions.NoProgress {
			installer.WithStatusWriter(statuswriter.Spinner())
		}

		bold := color.New(color.Bold).SprintFunc()
		packageName := args[0]
		pkgBuilder := client.PackageBuilder(packageName)
		var repoClient repoclient.RepoClient

		if len(installCmdOptions.Repository) > 0 {
			repoClient = repoClientset.ForRepoWithName(installCmdOptions.Repository)
			pkgBuilder.WithRepositoryName(installCmdOptions.Repository)
		} else {
			repos, err := repoClientset.Meta().GetReposForPackage(packageName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❗ Error: could not collect repository list: %v\n", err)
			}
			switch len(repos) {
			case 0:
				fmt.Fprintf(os.Stderr, "❗ Error: %v is not available\n", packageName)
				cliutils.ExitWithError()
			case 1:
				repoClient = repoClientset.ForRepo(repos[0])
				pkgBuilder.WithRepositoryName(repos[0].Name)
			default:
				names := make([]string, len(repos))
				for i := range repos {
					names[i] = repos[i].Name
				}
				for {
					fmt.Fprintf(os.Stderr,
						"%v is available from %v repositories. Please select the one to install from.\n",
						packageName, len(names))
					if repoName, err := cliutils.GetOption("", names); err != nil {
						fmt.Fprintf(os.Stderr, "invalid input: %v\n", err)
					} else {
						repoClient = repoClientset.ForRepoWithName(repoName)
						pkgBuilder.WithRepositoryName(repoName)
						break
					}
				}
			}
		}

		if installCmdOptions.Version == "" {
			var packageIndex repo.PackageIndex
			if err := repoClient.FetchPackageIndex(packageName, &packageIndex); err != nil {
				fmt.Fprintf(os.Stderr, "❗ Error: Could not fetch package metadata: %v\n", err)
				cliutils.ExitWithError()
			}
			installCmdOptions.Version = packageIndex.LatestVersion
			fmt.Fprintf(os.Stderr, "Version not specified. The latest version %v of %v will be installed.\n",
				installCmdOptions.Version, packageName)
		} else if !strings.HasPrefix(installCmdOptions.Version, "v") {
			installCmdOptions.Version = "v" + installCmdOptions.Version
		}

		pkgBuilder.WithVersion(installCmdOptions.Version)

		var manifest v1alpha1.PackageManifest
		if err := repoClient.FetchPackageManifest(packageName, installCmdOptions.Version, &manifest); err != nil {
			fmt.Fprintf(os.Stderr, "❗ Error: Could not fetch package manifest: %v\n", err)
			cliutils.ExitWithError()
		}

		installationPlan := []dependency.Requirement{}
		if manifest.Scope.IsCluster() {
			if len(args) != 1 {
				fmt.Fprintf(os.Stderr,
					"❌ %v has scope Cluster. Specifying an instance name for a ClusterPackage is not possible\n",
					packageName)
				cliutils.ExitWithError()
			}
			installationPlan = append(installationPlan,
				dependency.Requirement{PackageWithVersion: dependency.PackageWithVersion{
					Name:    packageName,
					Version: installCmdOptions.Version,
				}},
			)
		} else {
			var name string
			if len(args) != 2 {
				if installCmdOptions.Yes {
					fmt.Fprintf(os.Stderr, "Name not specified. Using default name: %v\n", packageName)
					name = packageName
				} else {
					fmt.Fprintf(os.Stderr, "%v has scope Namespaced. Please enter a name (default %v):\n", packageName, packageName)
					name = cliutils.GetInputStr("name")
					if name == "" {
						name = packageName
					}
				}
			} else {
				name = args[1]
			}
			ns := installCmdOptions.GetActualNamespace(ctx)
			pkgBuilder.WithName(name).WithNamespace(ns)
			installationPlan = append(installationPlan,
				dependency.Requirement{PackageWithVersion: dependency.PackageWithVersion{
					Name:    fmt.Sprintf("%v of type %v in namespace %v", name, packageName, ns),
					Version: installCmdOptions.Version,
				}},
			)
		}

		if installCmdOptions.IsValuesSet() {
			if values, err := installCmdOptions.ParseValues(&manifest, nil); err != nil {
				fmt.Fprintf(os.Stderr, "❌ invalid values in command line flags: %v\n", err)
				cliutils.ExitWithError()
			} else {
				pkgBuilder.WithValues(values)
			}
		} else {
			if values, err := cli.Configure(manifest, cli.WithUseDefaults(installCmdOptions.UseDefault)); err != nil {
				cancel()
			} else {
				pkgBuilder.WithValues(values)
			}
		}

		if !installCmdOptions.EnableAutoUpdates && !installCmdOptions.Yes {
			if cliutils.YesNoPrompt("Would you like to enable automatic updates?", false) {
				installCmdOptions.EnableAutoUpdates = true
			}
		}
		if installCmdOptions.EnableAutoUpdates {
			ok, err := clientutils.IsAutoUpdaterInstalled(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "! Error: Could not check whether glasskube-autoupdater is installed: %v\n", err)
			}
			if !ok {
				fmt.Fprintf(os.Stderr, "Please install glasskube-autoupdater for automatic updates to be applied.\n")
			}
		}

		pkgBuilder.WithAutoUpdates(installCmdOptions.EnableAutoUpdates)

		pkg := pkgBuilder.Build(manifest.Scope)

		if validationResult, err :=
			dm.Validate(ctx, pkg.GetName(), pkg.GetNamespace(), &manifest, installCmdOptions.Version); err != nil {
			fmt.Fprintf(os.Stderr, "❗ Error: Could not validate dependencies: %v\n", err)
			cliutils.ExitWithError()
		} else if len(validationResult.Conflicts) > 0 {
			fmt.Fprintf(os.Stderr, "❗ Error: %v cannot be installed due to conflicts: %v\n",
				packageName, validationResult.Conflicts)
			cliutils.ExitWithError()
		} else if len(validationResult.Requirements) > 0 {
			installationPlan = append(installationPlan, validationResult.Requirements...)
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

		createNamespace := false
		if installCmdOptions.NamespaceOptions.Namespace != "" {
			if ok, err := namespaces.Exists(ctx, cs, installCmdOptions.NamespaceOptions.Namespace); !ok {
				fmt.Fprintf(os.Stderr, " * Namespace %v does not exist and will be created\n",
					installCmdOptions.NamespaceOptions.Namespace)
				createNamespace = true
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred in the Namespace check:\n\n%v\n", err)
				cliutils.ExitWithError()
			}
		}

		if len(pkg.GetSpec().Values) > 0 {
			fmt.Fprintln(os.Stderr, bold("Configuration:"))
			printValueConfigurations(os.Stderr, pkg.GetSpec().Values)
			if _, err := valueResolver.Resolve(ctx, pkg.GetSpec().Values); err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Some values can not be resolved: %v\n", err)
			}
		}

		if !installCmdOptions.Yes && !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}

		if createNamespace {
			ns := &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: installCmdOptions.NamespaceOptions.Namespace,
				},
			}
			_, err := cs.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred in creating the Namespace:\n\n%v\n", err)
				cliutils.ExitWithError()
			}
		}
		if installCmdOptions.NoWait {
			if err := installer.Install(ctx, pkg, opts); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
				cliutils.ExitWithError()
			}
			fmt.Fprintf(os.Stderr,
				"☑️  %v is being installed in the background.\n"+
					"💡 Run \"glasskube describe %v\" to get the current status\n",
				packageName, packageName)
		} else {
			status, err := installer.InstallBlocking(ctx, pkg, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred during installation:\n\n%v\n", err)
				cliutils.ExitWithError()
			}
			if status != nil {
				switch status.Status {
				case string(condition.Ready):
					fmt.Fprintf(os.Stderr, "✅ %v is now installed in %v.\n", packageName, config.CurrentContext)
				default:
					fmt.Fprintf(os.Stderr, "❌ %v installation has status %v, reason: %v\nMessage: %v\n",
						packageName, status.Status, status.Reason, status.Message)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Installation status unknown - no error and no status have been observed (this is a bug).")
				cliutils.ExitWithError()
			}
		}
		if installCmdOptions.OutputOptions.Output != "" {
			output, err := clientutils.Format(installCmdOptions.Output.OutputFormat(), installCmdOptions.ShowAll, pkg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❗ Error: %v\n", err)
				cliutils.ExitWithError()
			}
			fmt.Println(output)
		}
	},
}

func cancel() {
	fmt.Fprintf(os.Stderr, "❌ Operation cancelled.")
	cliutils.ExitWithError()
}

func completeAvailablePackageNames(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	cfg, rawCfg := cliutils.RequireConfig(config.Kubeconfig)
	ctx := util.Must(clicontext.SetupContext(cmd.Context(), cfg, rawCfg))
	repoClient := cliutils.RepositoryClientset(ctx)
	var index repotypes.MetaIndex
	err := repoClient.Meta().FetchMetaIndex(&index)
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
	cfg, rawCfg := cliutils.RequireConfig(config.Kubeconfig)
	ctx := util.Must(clicontext.SetupContext(cmd.Context(), cfg, rawCfg))
	repoClient := cliutils.RepositoryClientset(ctx)
	repos, err := repoClient.Meta().GetReposForPackage(packageName)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	versionsMap := make(map[string]struct{})
	for _, r := range repos {
		var packageIndex repo.PackageIndex
		if err := repoClient.ForRepo(r).FetchPackageIndex(packageName, &packageIndex); err != nil {
			continue
		}
		for _, version := range packageIndex.Versions {
			if toComplete == "" || strings.HasPrefix(version.Version, toComplete) {
				versionsMap[version.Version] = struct{}{}
			}
		}
	}
	return maputils.KeysSorted(versionsMap), cobra.ShellCompDirectiveNoFileComp
}

func init() {
	installCmd.PersistentFlags().StringVarP(&installCmdOptions.Version, "version", "v", "",
		"Install a specific version")
	_ = installCmd.RegisterFlagCompletionFunc("version", completeAvailablePackageVersions)
	installCmd.PersistentFlags().BoolVar(&installCmdOptions.EnableAutoUpdates, "enable-auto-updates", false,
		"Enable automatic updates for this package")
	installCmd.PersistentFlags().StringVar(&installCmdOptions.Repository, "repository", installCmdOptions.Repository,
		"Specify the name of the package repository to install this package from")
	installCmd.PersistentFlags().BoolVar(&installCmdOptions.NoWait, "no-wait", false, "Perform non-blocking install")
	installCmd.PersistentFlags().BoolVarP(&installCmdOptions.Yes, "yes", "y", false, "Do not ask for any confirmation")
	installCmdOptions.ValuesOptions.AddFlagsToCommand(installCmd)
	installCmdOptions.OutputOptions.AddFlagsToCommand(installCmd)
	installCmdOptions.NamespaceOptions.AddFlagsToCommand(installCmd)
	installCmdOptions.DryRunOptions.AddFlagsToCommand(installCmd)
	installCmd.MarkFlagsMutuallyExclusive("version", "enable-auto-updates")
	installCmd.MarkFlagsMutuallyExclusive("no-wait", "dry-run")
	RootCmd.AddCommand(installCmd)
}

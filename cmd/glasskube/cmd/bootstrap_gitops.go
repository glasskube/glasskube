package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	client2 "github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type bootstrapGitopsOptions struct {
	repo string
}

var bootstrapGitopsCmdOptions = bootstrapGitopsOptions{}

var bootstrapGitopsCmd = &cobra.Command{
	Use:    "gitops",
	Short:  "Bootstrap Glasskube with a GitOps tool",
	PreRun: cliutils.SetupClientContext(false, util.Pointer(true)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)
		client := bootstrap.NewBootstrapClient(cfg)
		ctx := cmd.Context()

		var installedVersion *semver.Version
		if installedVersionRaw, err := clientutils.GetPackageOperatorVersion(ctx); err != nil {
			if !apierrors.IsNotFound(err) {
				fmt.Fprintf(os.Stderr, "could not determine installed version: %v\n", err)
				cliutils.ExitWithError()
			}
		} else if installedVersion, err = semver.NewVersion(installedVersionRaw); err != nil {
			fmt.Fprintf(os.Stderr, "could not parse installed version: %v\n", err)
			cliutils.ExitWithError()
		} else if installedVersion != nil {
			fmt.Fprintf(os.Stderr, "gitops bootstrapping on a bootstrapped cluster is not supported yet. \n")
			// TODO cliutils.ExitWithError()
		}

		// verifyLegalUpdate(ctx, installedVersion, targetVersion)
		currentContext := color.New(color.Bold).Sprint(clicontext.RawConfigFromContext(ctx).CurrentContext)
		fmt.Fprintf(os.Stderr, "Glasskube and ArgoCD will be installed in context %s.\n", currentContext)
		if !bootstrapCmdOptions.yes && !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}

		// regularly bootstrap in the cluster
		// TODO without a git client, this means it only supports GitHub for now
		basePath := strings.ReplaceAll(bootstrapGitopsCmdOptions.repo, "github.com", "raw.githubusercontent.com")
		manifestsPath := fmt.Sprintf("%v/bootstrap/glasskube/manifests.yaml", basePath)
		_, err := client.Bootstrap(ctx, bootstrap.BootstrapOptions{
			Url:        manifestsPath,
			Type:       bootstrap.BootstrapTypeAio,
			GitopsMode: true,
			Force:      true, // TODO check why this is necessary (something fails for PackageRepository CRD)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			cliutils.ExitWithError()
		}

		cliutils.SetupClientContext(true, util.Pointer(true))(cmd, make([]string, 0))
		// install argo-cd package
		// TODO use the version of argo given in /packages/argo-cd/clusterpackage.yaml instead
		argoCdPkg := client2.PackageBuilder("argo-cd").
			WithRepositoryName("glasskube").
			WithVersion("v2.11.5+1").
			BuildClusterPackage()

		if _, err := install.NewInstaller(cliutils.PackageClient(cmd.Context())).
			InstallBlocking(cmd.Context(), argoCdPkg, metav1.CreateOptions{}); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred installing argo-cd package:\n%v\n", err)
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "argo-cd package has been installed\n")

		// apply bootstrap/glasskube-application.yaml into the cluster
		// TODO without a git client, this means it only supports GitHub for now
		appPath := fmt.Sprintf("%v/bootstrap/glasskube-application.yaml", basePath)
		if objs, err := clientutils.FetchResources(appPath); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred fetching the bootstrap application:\n%v\n", err)
			cliutils.ExitWithError()
		} else if client, err := dynamic.NewForConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred initializing the dynamic client:\n%v\n", err)
			cliutils.ExitWithError()
		} else {
			for _, obj := range objs {
				if _, err = client.Resource(schema.GroupVersionResource{
					Group:    "argoproj.io",
					Version:  "v1alpha1",
					Resource: "applications",
				}).Namespace(obj.GetNamespace()).Apply(ctx, obj.GetName(), &obj, metav1.ApplyOptions{
					FieldManager: "glasskube", // TODO not sure if correct
				}); err != nil {
					fmt.Fprintf(os.Stderr, "\nAn error occurred applying the bootstrap application:\n%v\n", err)
					cliutils.ExitWithError()
				}
			}
		}

		// wait for argo to be ready, then print...
		// TODO

		fmt.Fprintf(os.Stderr, "DONE\n")

	},
}

func init() {
	bootstrapGitopsCmd.Flags().StringVar(&bootstrapGitopsCmdOptions.repo, "repo", "", "URL of the GitOps Repository")
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type bootstrapGitOptions struct {
	url string
}

var bootstrapGitCmdOptions = bootstrapGitOptions{}

var bootstrapGitCmd = &cobra.Command{
	Use:    "git",
	Short:  "Bootstrap Glasskube with a GitOps tool",
	PreRun: cliutils.SetupClientContext(false, util.Pointer(true)),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cliutils.RequireConfig(config.Kubeconfig)
		bootstrapClient := bootstrap.NewBootstrapClient(cfg)
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
			cliutils.ExitWithError()
		}

		currentContext := color.New(color.Bold).Sprint(clicontext.RawConfigFromContext(ctx).CurrentContext)
		fmt.Fprintf(os.Stderr, "Glasskube and ArgoCD will be installed in cluster %s.\n", currentContext)
		if !bootstrapCmdOptions.yes && !cliutils.YesNoPrompt("Continue?", true) {
			cancel()
		}

		// regularly bootstrap in the cluster
		// without a git client, this means it only supports GitHub for now
		basePath := strings.ReplaceAll(bootstrapGitCmdOptions.url, "github.com", "raw.githubusercontent.com")
		basePath = basePath + "/main"
		manifestsPath := fmt.Sprintf("%v/bootstrap/glasskube/glasskube.yaml", basePath)
		_, err := bootstrapClient.Bootstrap(ctx, bootstrap.BootstrapOptions{
			Url:        manifestsPath,
			Type:       bootstrap.BootstrapTypeAio,
			GitopsMode: true,
			Force:      true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "\nInstalling argo-cd...")
		cliutils.SetupClientContext(true, util.Pointer(true))(cmd, make([]string, 0))

		// get defined argo-cd version from repo:
		argocdPath := fmt.Sprintf("%v/packages/argo-cd/clusterpackage.yaml", basePath)
		argocdVersion := "v2.11.5+1"
		if objs, err := clientutils.FetchResources(argocdPath); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred fetching argo-cd version (will use %v instead): %v\n",
				argocdVersion, err)
		} else {
			for _, obj := range objs {
				if v, ok, err :=
					unstructured.NestedString(obj.Object, "spec", "packageInfo", "version"); !ok || err != nil {
					fmt.Fprintf(os.Stderr, "\nAn error occurred trying to read the argo-cd version "+
						"(will use %v instead): %v\n", argocdVersion, err)
				} else {
					argocdVersion = v
					break
				}
			}
		}
		// install argo-cd package
		argoCdPkg := client.PackageBuilder("argo-cd").
			WithRepositoryName("glasskube").
			WithVersion(argocdVersion).
			BuildClusterPackage()
		if _, err := install.NewInstaller(cliutils.PackageClient(cmd.Context())).
			InstallBlocking(cmd.Context(), argoCdPkg, metav1.CreateOptions{}); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred installing argo-cd:\n%v\n", err)
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "\nargo-cd package has been installed!\n\nApplying the glasskube argo-cd application...")

		// apply bootstrap/glasskube-application.yaml into the cluster
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

		fmt.Fprintf(os.Stderr, "\nglasskube argo-cd application applied successfully!\n\n")
		fmt.Fprintf(os.Stderr, "You have successfully installed Glasskube and ArgoCD in this cluster.\n")
		fmt.Fprintf(os.Stderr, "Right now, ArgoCD is starting up and will soon sync with your GitOps repo. ")
		fmt.Fprintf(os.Stderr, "Please note, that this might take a couple of minutes!\n\n")
		fmt.Fprintf(os.Stderr, "Run glasskube serve to open the Glasskube UI and either open the ArgoCD UI there, "+
			"or with the command `glasskube open argo-cd` – but of course you can also use the Argo CLI.\n"+
			"Follow the ArgoCD docs to get and reset the password to log in: "+
			"https://argo-cd.readthedocs.io/en/stable/getting_started/#4-login-using-the-cli)")
	},
}

func init() {
	bootstrapGitCmd.Flags().StringVar(&bootstrapGitCmdOptions.url, "url", "", "URL of the GitOps Repository")
}

package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

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
	"github.com/glasskube/glasskube/pkg/open"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

type bootstrapGitOptions struct {
	url      string
	branch   string
	username string
	token    string
}

var bootstrapGitCmdOptions = bootstrapGitOptions{}

const doneMessage = `

You have successfully installed Glasskube and ArgoCD in this cluster.
You can now open the ArgoCD UI with

glasskube open argo-cd

username: admin
password: %v

Please reset the ArgoCD admin password! https://argo-cd.readthedocs.io/en/stable/getting_started/#4-login-using-the-cli
`

var bootstrapGitCmd = &cobra.Command{
	Use:    "git",
	Short:  "Bootstrap Glasskube with a GitOps tool",
	PreRun: cliutils.SetupClientContext(false, util.Pointer(true)),
	Run: func(cmd *cobra.Command, args []string) {
		if (bootstrapGitCmdOptions.username == "" && bootstrapGitCmdOptions.token != "") ||
			(bootstrapGitCmdOptions.username != "" && bootstrapGitCmdOptions.token == "") {
			fmt.Fprintf(os.Stderr, "For private repos please provide both username and token!\n")
			cliutils.ExitWithError()
		}
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

		// build url to fetch glasskube manifests from
		bootstrapGitCmdOptions.url, _ = strings.CutSuffix(bootstrapGitCmdOptions.url, ".git")
		baseURL, err := url.Parse(bootstrapGitCmdOptions.url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse URL: %v\n", err)
			cliutils.ExitWithError()
		}
		baseURL.Host = "raw.githubusercontent.com" // without a git client, this means it only supports GitHub for now
		if bootstrapGitCmdOptions.token != "" {
			baseURL.User = url.UserPassword(bootstrapGitCmdOptions.username, bootstrapGitCmdOptions.token)
		}
		branchPath := "main"
		if bootstrapGitCmdOptions.branch != "" {
			branchPath = bootstrapGitCmdOptions.branch
		}
		baseURL = baseURL.JoinPath(branchPath)

		manifestsURL := baseURL.JoinPath("bootstrap", "glasskube", "glasskube.yaml")
		// regularly bootstrap in the cluster
		_, err = bootstrapClient.Bootstrap(ctx, bootstrap.BootstrapOptions{
			Url:        manifestsURL.String(),
			Type:       bootstrap.BootstrapTypeAio,
			GitopsMode: true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred during bootstrap:\n%v\n", err)
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "\nInstalling argo-cd...")
		cliutils.SetupClientContext(true, util.Pointer(true))(cmd, make([]string, 0))

		// get defined argo-cd version and repo from repo:
		argocdPath := baseURL.JoinPath("packages", "argo-cd", "clusterpackage.yaml").String()
		argocdVersion := "v2.11.7+1"
		argocdRepo := "glasskube"
		if objs, err := clientutils.FetchResourcesFromUrl(argocdPath); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred fetching argo-cd clusterpackage "+
				"(will use %v from repo %v instead): %v\n", argocdVersion, argocdRepo, err)
		} else if len(objs) != 1 {
			fmt.Fprintf(os.Stderr, "\nUnexpectedly found %v objects in %v – there should only be one. Aborting.\n",
				len(objs), argocdPath)
			cliutils.ExitWithError()
		} else {
			obj := objs[0]
			if v, ok, err :=
				unstructured.NestedString(obj.Object, "spec", "packageInfo", "version"); !ok || err != nil {
				fmt.Fprintf(os.Stderr, "\nAn error occurred trying to read the argo-cd version "+
					"(will use %v instead): %v\n", argocdVersion, err)
			} else {
				argocdVersion = v
			}
			if repo, ok, err :=
				unstructured.NestedString(obj.Object, "spec", "packageInfo", "repositoryName"); !ok || err != nil {
				fmt.Fprintf(os.Stderr, "\nAn error occurred trying to read the argo-cd repositoryName "+
					"(will use %v instead): %v\n", argocdRepo, err)
			} else {
				argocdRepo = repo
			}
		}
		// install argo-cd package
		argoCdPkg := client.PackageBuilder("argo-cd").
			WithRepositoryName(argocdRepo).
			WithVersion(argocdVersion).
			BuildClusterPackage()
		if _, err := install.NewInstaller(cliutils.PackageClient(cmd.Context())).
			InstallBlocking(cmd.Context(), argoCdPkg, metav1.CreateOptions{}); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred installing argo-cd:\n%v\n", err)
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "\nargo-cd package has been installed!\n\nApplying the glasskube argo-cd application...")

		// apply the repository secret if the repo is not public
		argocdNamespace := "argocd"
		clientset := kubernetes.NewForConfigOrDie(cfg)
		if bootstrapGitCmdOptions.token != "" {
			repoSecret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "gitops-repository",
					Labels: map[string]string{
						"argocd.argoproj.io/secret-type": "repository",
					},
				},
				StringData: map[string]string{
					"type":     "git",
					"url":      bootstrapGitCmdOptions.url,
					"username": bootstrapGitCmdOptions.username,
					"password": bootstrapGitCmdOptions.token,
				},
			}
			if _, err := clientset.CoreV1().Secrets(argocdNamespace).
				Create(ctx, repoSecret, metav1.CreateOptions{}); err != nil {
				fmt.Fprintf(os.Stderr, "\nAn error occurred setting up the repository secret:\n%v\n", err)
				cliutils.ExitWithError()
			}
		}

		// apply bootstrap/glasskube-application.yaml into the cluster
		// appPath := fmt.Sprintf("%v/bootstrap/glasskube-application.yaml", baseURL.String())
		appPath := baseURL.JoinPath("bootstrap", "glasskube-application.yaml").String()
		if objs, err := clientutils.FetchResourcesFromUrl(appPath); err != nil {
			fmt.Fprintf(os.Stderr, "\nAn error occurred fetching the bootstrap application:\n%v\n", err)
			cliutils.ExitWithError()
		} else {
			// re-initialize the rest mapper because with argo-cd installed, new kinds will be available
			if err := bootstrapClient.InitRestMapper(); err != nil {
				fmt.Fprintf(os.Stderr, "\nAn error occurred setting up the restmapper: %v\n", err)
				cliutils.ExitWithError()
			}
			for _, obj := range objs {
				gvk := obj.GroupVersionKind()
				mapping, err := bootstrapClient.Mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
				if err != nil {
					fmt.Fprintf(os.Stderr, "\nAn error occurred preparing %v/%v:\n%v\n", obj.GetNamespace(), obj.GetName(), err)
					cliutils.ExitWithError()
				}
				if _, err := bootstrapClient.Client.Resource(mapping.Resource).
					Namespace(obj.GetNamespace()).Apply(ctx, obj.GetName(), &obj, metav1.ApplyOptions{
					FieldManager: "glasskube", // TODO not sure if correct
				}); err != nil {
					fmt.Fprintf(os.Stderr, "\nAn error occurred applying the bootstrap application:\n%v\n", err)
					cliutils.ExitWithError()
				}
			}
		}

		// try to open argo-cd every 5 seconds until it succeeds in order to check whether it is ready
		fmt.Fprintf(os.Stderr, "\nglasskube argo-cd application applied successfully!\n\nWaiting for argo-cd to be ready...")
		argoCheckCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		tick := time.NewTicker(5 * time.Second)
		tickC := tick.C
		defer tick.Stop()
	WaitLoop:
		for {
			select {
			case <-argoCheckCtx.Done():
				fmt.Fprintf(os.Stderr, "\nargo-cd has not become ready in the specified timeout: %v\n", ctx.Err())
				break WaitLoop
			case <-tickC:
				if ready, _ := open.NewOpener().HasReadyPod(argoCheckCtx, argoCdPkg, "", 0); ready {
					break WaitLoop
				}
				// else: generally assume an error means not ready yet
			}
		}

		// get initial admin password to print it
		var adminPwStr string
		if adminPw, err := clientset.CoreV1().Secrets(argocdNamespace).
			Get(ctx, "argocd-initial-admin-secret", metav1.GetOptions{}); err != nil {
			fmt.Fprintf(os.Stderr, "\nCould not obtain initial ArgoCD password: %v\n\n", err)
		} else {
			adminPwStr = string(adminPw.Data["password"])
		}

		fmt.Fprintf(os.Stderr, doneMessage, adminPwStr)
	},
}

func init() {
	bootstrapGitCmd.Flags().StringVar(&bootstrapGitCmdOptions.url, "url", "",
		"URL of the GitOps Repository")
	bootstrapGitCmd.Flags().StringVar(&bootstrapGitCmdOptions.branch, "branch", "",
		"Branch of the GitOps Repository to use (default is main)")
	bootstrapGitCmd.Flags().StringVar(&bootstrapGitCmdOptions.username, "username", "",
		"Username to use for authentication")
	bootstrapGitCmd.Flags().StringVar(&bootstrapGitCmdOptions.token, "token", os.Getenv("GITHUB_TOKEN"),
		"Token to use for authentication. If not set, it will use the GITHUB_TOKEN environment variable")
}

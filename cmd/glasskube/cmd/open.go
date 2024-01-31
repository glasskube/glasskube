package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/repo"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/list"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [package-name] [entrypoint]",
	Short: "Open the Web UI of a package",
	Long: `Open the Web UI of a package.
If the package manifest has more than one entrypoint, specify the name of the entrypoint to open.`,
	Args:   cobra.RangeArgs(1, 2),
	PreRun: cliutils.SetupClientContext(true),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]
		var entrypointName string
		if len(args) == 2 {
			entrypointName = args[1]
		}
		pkgClient := client.FromContext(cmd.Context())
		pkg, err := list.Get(pkgClient, cmd.Context(), pkgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred retreiving the package:\n\n%v\n", err)
			os.Exit(1)
			return
		}

		pkgInfo := &v1alpha1.PackageInfo{
			Spec: v1alpha1.PackageInfoSpec{
				Name: pkg.Name,
			},
		}
		err = repo.FetchPackageManifest(cmd.Context(), pkgInfo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred retreiving the package manifest:\n\n%v\n", err)
			os.Exit(1)
			return
		}

		stopCh := make(chan os.Signal, 1)
		signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

		commands := make([]*exec.Cmd, 0)
		var browserUrl string
		manifest := pkgInfo.Status.Manifest
		namespace := manifest.DefaultNamespace
		for _, ep := range manifest.Entrypoints {
			kubectl := exec.Command(
				"kubectl", "-n", namespace, "port-forward",
				fmt.Sprintf("service/%s", ep.ServiceName),
				fmt.Sprintf("%d:%d", ep.Port, ep.Port))
			kubectl.Stdout = os.Stdout
			kubectl.Stderr = os.Stderr
			fmt.Printf("Running %v\n", kubectl)

			err = kubectl.Start()
			if err != nil {
				fmt.Printf("Error starting kubectl port-forward: %v\n", err)
				os.Exit(1)
			}
			commands = append(commands, kubectl)

			go func() {
				err = kubectl.Wait()
				if err != nil {
					fmt.Printf("Error waiting for kubectl port-forward to exit: %v\n", err)
				}
			}()

			if len(manifest.Entrypoints) == 1 || ep.Name == entrypointName {
				browserUrl = fmt.Sprintf("http://localhost:%d", ep.Port)
			}
		}

		if len(browserUrl) > 0 {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("%s is now reachable at %s\n", pkgName, browserUrl)
			_ = cliutils.OpenInBrowser(browserUrl)
		}

		<-stopCh

		for _, cmd := range commands {
			err = cmd.Process.Signal(os.Interrupt)
			if err != nil {
				fmt.Printf("Error sending interrupt signal to kubectl port-forward: %v\n", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}

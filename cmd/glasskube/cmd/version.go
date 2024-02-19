package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var versioncmd = &cobra.Command{
	Use:    "version",
	Short:  "Print the version of glasskube and package-operator",
	Long:   `Print the version of glasskube and package-operator`,
	PreRun: cliutils.SetupClientContext(false),
	Run: func(cmd *cobra.Command, args []string) {
		glasskubeVersion := config.Version
		fmt.Fprintf(os.Stderr, "glasskube: v%s\n", glasskubeVersion)
		operatorVersion, err := getPackageOperatorVersion(cmd.Context())
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ no deployments found in the glasskube-system namespace\n")
		} else {
			fmt.Fprintf(os.Stderr, "package-operator: %s\n", operatorVersion)
		}
	},
}

func init() {
	RootCmd.AddCommand(versioncmd)
}

func getPackageOperatorVersion(ctx context.Context) (string, error) {
	config := client.ConfigFromContext(ctx)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	namespace := "glasskube-system"
	deploymentName := "glasskube-controller-manager"
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	containers := deployment.Spec.Template.Spec.Containers
	for _, container := range containers {
		if container.Name == "manager" {
			ref, err := name.ParseReference(container.Image)
			if err != nil {
				return "", err
			}
			return ref.Identifier(), nil
		}
	}
	return "", nil
}

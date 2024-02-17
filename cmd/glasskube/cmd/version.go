package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var versioncmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of glasskube and package-operator",
	Long:  `Print the version of glasskube and package-operator`,
	Run: func(cmd *cobra.Command, args []string) {
		glasskubeVersion := config.Version
		fmt.Fprintf(os.Stderr, "glasskube: v%s\n", glasskubeVersion)
		operatorVersion, err := getPackageOperatorVersion()
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

func getPackageOperatorVersion() (string, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", home+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	namespace := "glasskube-system"
	deploymentName := "glasskube-controller-manager"
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
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
			parts := strings.Split(ref.Identifier(), ":")
			return parts[0], nil
		}
	}
	return "", nil
}

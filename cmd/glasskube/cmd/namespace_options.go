package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceOptions struct {
	Namespace string
}

func (opt *NamespaceOptions) AddFlagsToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opt.Namespace, "namespace", "n", opt.Namespace, "Namespace for resources")
}

func (opt *NamespaceOptions) GetActualNamespace(ctx context.Context) string {
	if opt.Namespace != "" {
		return opt.Namespace
	} else {
		rawConfig := clicontext.RawConfigFromContext(ctx)
		if current, ok := rawConfig.Contexts[rawConfig.CurrentContext]; ok && current.Namespace != "" {
			return current.Namespace
		} else {
			return "default"
		}
	}
}

func DeleteNamespace(ctx context.Context, namespace string) {
	client := cliutils.PackageClient(ctx)
	config := clicontext.ConfigFromContext(ctx)

	var pkgs v1alpha1.PackageList
	if err := client.Packages(namespace).GetAll(ctx, &pkgs); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error listing packages in namespace: %v\n", err)
		cliutils.ExitWithError()
	}

	var namespacePackages []string
	for _, pkg := range pkgs.Items {
		if pkg.Namespace == namespace {
			namespacePackages = append(namespacePackages, pkg.Name)
		}
	}

	if len(namespacePackages) > 0 {
		fmt.Printf("❌ Namespace %s cannot be deleted because it contains other packages: %v\n",
			namespace, strings.Join(namespacePackages, ", "))
		cliutils.ExitWithError()
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ error creating Kubernetes client: %v\n", err)
		cliutils.ExitWithError()
	}

	err = clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ error deleting namespace: %v\n", err)
		cliutils.ExitWithError()
	}

	fmt.Printf("Namespace %s has been deleted.\n", namespace)
}

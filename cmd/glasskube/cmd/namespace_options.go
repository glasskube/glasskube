package cmd

import (
	"context"
	"strings"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/kubeconfig"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceOptions struct {
	Namespace string
}

func (opt *NamespaceOptions) AddFlagsToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opt.Namespace, "namespace", "n", opt.Namespace, "Namespace for resources")
	_ = cmd.RegisterFlagCompletionFunc("namespace", completeNamespaces)
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

func completeNamespaces(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) (names []string, dir cobra.ShellCompDirective) {
	ctx := cmd.Context()
	dir = cobra.ShellCompDirectiveNoFileComp
	if config, _, err := kubeconfig.New(config.Kubeconfig); err != nil {
		dir |= cobra.ShellCompDirectiveError
	} else if client, err := kubernetes.NewForConfig(config); err != nil {
		dir |= cobra.ShellCompDirectiveError
	} else if nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{}); err != nil {
		dir |= cobra.ShellCompDirectiveError
	} else {
		for _, ns := range nsList.Items {
			if toComplete == "" || strings.HasPrefix(toComplete, ns.Name) {
				names = append(names, ns.Name)
			}
		}
	}
	return
}

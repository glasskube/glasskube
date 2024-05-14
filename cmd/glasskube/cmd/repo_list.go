package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/condition"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var repoListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Print a list of package repositories installed in the current cluster",
	PreRun:  cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client := cliutils.PackageClient(ctx)

		var repos v1alpha1.PackageRepositoryList
		if err := client.PackageRepositories().GetAll(ctx, &repos); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error listing package repository: %v\n", err)
			cliutils.ExitWithError()
		}

		util.SortBy(repos.Items, func(repo v1alpha1.PackageRepository) string { return repo.Name })

		_ = cliutils.PrintTable(os.Stdout, repos.Items,
			[]string{"NAME", "URL", "AUTHENTICATION", "STATUS", "MESSAGE"},
			func(repo v1alpha1.PackageRepository) []string {
				condition := meta.FindStatusCondition(repo.Status.Conditions, string(condition.Ready))
				authType := "None"
				if repo.Spec.Auth != nil {
					if repo.Spec.Auth.Basic != nil {
						authType = "Basic"
					} else if repo.Spec.Auth.Bearer != nil {
						authType = "Bearer"
					}
				}
				status := "Unknown"
				message := ""
				if condition != nil {
					if condition.Status == metav1.ConditionTrue {
						status = "Ready"
					} else if condition.Status == metav1.ConditionFalse {
						status = "Not ready"
					}
					message = condition.Message
				}
				return []string{
					repo.Name,
					repo.Spec.Url,
					authType,
					status,
					message,
				}
			})

		cliutils.ExitSuccess()
	},
}

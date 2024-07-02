package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var repoAddCmdOptions = repoOptions{}

var repoAddCmd = &cobra.Command{
	Use:    "add <name> <url>",
	Short:  "Add a package repository to the current cluster",
	Args:   cobra.ExactArgs(2),
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var defaultRepo *v1alpha1.PackageRepository

		ctx := cmd.Context()
		client := cliutils.PackageClient(ctx)
		repoName := args[0]
		repoAddCmdOptions.Url = args[1]

		if err := repoAddCmdOptions.Normalize(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			cliutils.ExitWithError()
		}

		repo := v1alpha1.PackageRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: repoName,
			},
			Spec: v1alpha1.PackageRepositorySpec{
				Url: repoAddCmdOptions.Url,
			},
		}

		repo.Spec.Auth = repoAddCmdOptions.SetAuth()

		if repoAddCmdOptions.Default {
			defaultRepo, err = cliutils.GetDefaultRepo(ctx)

			if errors.Is(err, cliutils.NoDefaultRepo) {
				repo.SetDefaultRepository()
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "❌ error getting the default package repository: %v\n", err)
				cliutils.ExitWithError()
			} else if defaultRepo.Name != repoName {
				defaultRepo.SetDefaultRepositoryBool(false)
				if err := client.PackageRepositories().Update(ctx, defaultRepo); err != nil {
					fmt.Fprintf(os.Stderr, "❌ error updating current default package repository: %v\n", err)
					cliutils.ExitWithError()
				}
				repo.SetDefaultRepository()
			}
		}
		if err := client.PackageRepositories().Create(ctx, &repo, metav1.CreateOptions{}); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error creating package repository: %v\n", err)

			if repoAddCmdOptions.Default && defaultRepo != nil && defaultRepo.Name != repoName {
				defaultRepo.SetDefaultRepositoryBool(true)
				if err := client.PackageRepositories().Update(ctx, defaultRepo); err != nil {
					fmt.Fprintf(os.Stderr, "❌ error rolling back to default package repository: %v\n", err)
				}
			}

			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "✅ package repository %v added\n", repoName)
		cliutils.ExitSuccess()
	},
}

func init() {
	repoAddCmdOptions.BindToCmdFlags(repoAddCmd, false)
}

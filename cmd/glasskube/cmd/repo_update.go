package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/spf13/cobra"
)

var repoUpdateCmdOptions = repoOptions{}

var repoUpdateCmd = &cobra.Command{
	Use:    "update <name>",
	Short:  "Update a package repository for the current cluster",
	Args:   cobra.ExactArgs(1),
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var repo v1alpha1.PackageRepository
		var defaultRepo *v1alpha1.PackageRepository

		ctx := cmd.Context()
		client := cliutils.PackageClient(ctx)
		repoName := args[0]

		if err := repoUpdateCmdOptions.Normalize(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			cliutils.ExitWithError()
		}

		if err := client.PackageRepositories().Get(ctx, repoName, &repo); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error getting the package repository: %v\n", err)
			cliutils.ExitWithError()
		}

		if repoUpdateCmdOptions.Auth == repoNoAuth {
			repo.Spec.Auth = nil
		} else {
			repo.Spec.Auth = repoUpdateCmdOptions.SetAuth()
		}

		if repoUpdateCmdOptions.Url != "" {
			repo.Spec.Url = repoUpdateCmdOptions.Url
		}

		if repoUpdateCmdOptions.Default {
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

		if err := client.PackageRepositories().Update(ctx, &repo); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error updating the package repository: %v\n", err)

			if repoUpdateCmdOptions.Default && defaultRepo != nil && defaultRepo.Name != repoName {
				defaultRepo.SetDefaultRepositoryBool(true)
				if err := client.PackageRepositories().Update(ctx, defaultRepo); err != nil {
					fmt.Fprintf(os.Stderr, "❌ error rolling back to default package repository: %v\n", err)
				}
			}
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "✅ package repository %v updated\n", repoName)
		cliutils.ExitSuccess()
	},
}

func init() {
	repoUpdateCmdOptions.BindToCmdFlags(repoUpdateCmd, true)
}

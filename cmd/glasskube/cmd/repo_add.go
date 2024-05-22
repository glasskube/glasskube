package cmd

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var repoAddCmdOptions = repoOptions{}

var repoAddCmd = &cobra.Command{
	Use:    "add [name] [url]",
	Short:  "Add a package repository to the current cluster",
	Args:   cobra.ExactArgs(2),
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client := cliutils.PackageClient(ctx)
		repoName := args[0]
		repoUrl := args[1]

		if err := repoAddCmdOptions.Normalize(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			cliutils.ExitWithError()
		}

		repo := v1alpha1.PackageRepository{
			ObjectMeta: metav1.ObjectMeta{
				Name: repoName,
			},
			Spec: v1alpha1.PackageRepositorySpec{
				Url: repoUrl,
			},
		}

		switch repoAddCmdOptions.Auth {
		case repoAddBasicAuth:
			if len(repoAddCmdOptions.Username) == 0 {
				fmt.Fprintln(os.Stderr, "Basic authentication was requested. Please enter a username:")
				for {
					username := cliutils.GetInputStr("username")
					if len(username) > 0 {
						repoAddCmdOptions.Username = username
						break
					}
				}
			}
			if len(repoAddCmdOptions.Password) == 0 {
				fmt.Fprintln(os.Stderr, "Basic authentication was requested. Please enter a password:")
				for {
					password := cliutils.GetInputStr("password")
					if len(password) > 0 {
						repoAddCmdOptions.Password = password
						break
					}
				}
			}
			repo.Spec.Auth = &v1alpha1.PackageRepositoryAuthSpec{
				Basic: &v1alpha1.PackageRepositoryBasicAuthSpec{
					Username: &repoAddCmdOptions.Username,
					Password: &repoAddCmdOptions.Password,
				},
			}
		case repoAddBearerAuth:
			if len(repoAddCmdOptions.Token) == 0 {
				fmt.Fprintln(os.Stderr, "Bearer authentication was requested. Please enter a token:")
				for {
					token := cliutils.GetInputStr("token")
					if len(token) > 0 {
						repoAddCmdOptions.Token = token
						break
					}
				}
			}
			repo.Spec.Auth = &v1alpha1.PackageRepositoryAuthSpec{
				Bearer: &v1alpha1.PackageRepositoryBearerAuthSpec{
					Token: &repoAddCmdOptions.Token,
				},
			}
		}

		if repoAddCmdOptions.Default {
			repo.SetDefaultRepository()
		}
		if err := client.PackageRepositories().Create(ctx, &repo); err != nil {
			fmt.Fprintf(os.Stderr, "❌ error creating package repository: %v\n", err)
			cliutils.ExitWithError()
		}

		fmt.Fprintf(os.Stderr, "✅ package repository %v added\n", repoName)
		cliutils.ExitSuccess()
	},
}

func init() {
	repoAddCmdOptions.BindToCmdFlags(repoAddCmd)
}

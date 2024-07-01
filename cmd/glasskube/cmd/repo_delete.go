package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var repoDeleteCmd = &cobra.Command{
	Use:    "delete [repositoryName]",
	Short:  "Delete a repository",
	Args:   cobra.ExactArgs(1),
	PreRun: cliutils.SetupClientContext(true, &rootCmdOptions.SkipUpdateCheck),
	Run: func(cmd *cobra.Command, args []string) {
		repositoryName := args[0]
		ctx := cmd.Context()
		deleteRepository(ctx, repositoryName)
	},
}

func deleteRepository(ctx context.Context, repoName string) {
	client := cliutils.PackageClient(ctx)

	var targetRepo v1alpha1.PackageRepository
	if err := client.PackageRepositories().Get(ctx, repoName, &targetRepo); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error listing package repository: %v\n", err)
		cliutils.ExitWithError()
	}

	var pkgs v1alpha1.PackageList
	if err := client.Packages("").GetAll(ctx, &pkgs); err != nil {
		fmt.Fprintf(os.Stderr, "Could not list packages: %v", err)
		cliutils.ExitWithError()
	}

	var clpkgs v1alpha1.ClusterPackageList
	if err := client.ClusterPackages().GetAll(ctx, &clpkgs); err != nil {
		fmt.Fprintf(os.Stderr, "Could not list Cluster packages: %v", err)
		cliutils.ExitWithError()
	}

	repoPackages := getPackagesFromRepo(clpkgs, pkgs, repoName)
	if len(repoPackages) > 0 {
		fmt.Printf("Repository %s cannot be deleted, because the following packages are installed from this repository: %v\n",
			repoName, strings.Join(repoPackages, ", "))
		cliutils.ExitWithError()
	}

	if !cliutils.YesNoPrompt(fmt.Sprintf("Repository %s will now be deleted. Do you want to continue?", repoName), false) {
		fmt.Println("❌ Repository Deletion Cancelled")
		cliutils.ExitWithError()
	}
	err := client.PackageRepositories().Delete(ctx, &targetRepo, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Error deleting repository:", err)
		cliutils.ExitWithError()
	}

	fmt.Printf("Repository %s has been deleted.\n", repoName)
}

func getPackagesFromRepo(clpkgs v1alpha1.ClusterPackageList, pkgs v1alpha1.PackageList, repoName string) []string {
	var repoPackages []string
	for _, clpkg := range clpkgs.Items {
		if clpkg.Spec.PackageInfo.RepositoryName == repoName {
			repoPackages = append(repoPackages, clpkg.Name)
		}
	}

	for _, pkg := range pkgs.Items {
		if pkg.Spec.PackageInfo.RepositoryName == repoName {
			repoPackages = append(repoPackages, pkg.Name)
		}
	}
	return repoPackages
}

func init() {
	repoCmd.AddCommand(repoDeleteCmd)
}

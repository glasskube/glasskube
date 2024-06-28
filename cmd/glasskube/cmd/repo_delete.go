package cmd

import (
	"context"
	"fmt"
	"os"

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

	var repos v1alpha1.PackageRepositoryList
	if err := client.PackageRepositories().GetAll(ctx, &repos); err != nil {
		fmt.Fprintf(os.Stderr, "❌ error listing package repository: %v\n", err)
		cliutils.ExitWithError()
	}

	var targetRepo *v1alpha1.PackageRepository
	for i := range repos.Items {
		if repos.Items[i].Name == repoName {
			targetRepo = &repos.Items[i]
			break
		}
	}

	if targetRepo == nil {
		fmt.Fprintf(os.Stderr, "❌ repository %s not found\n", repoName)
		cliutils.ExitWithError()
	}

	var pkgs v1alpha1.PackageList
	_ = client.Packages("").GetAll(ctx, &pkgs)

	var clpkgs v1alpha1.ClusterPackageList
	_ = client.ClusterPackages().GetAll(ctx, &clpkgs)

	repoPackages := getPackagesFromRepo(clpkgs, pkgs, repoName)
	if len(repoPackages) > 0 {
		fmt.Printf("Repository %s cannot be deleted, because the following packages are installed from this repository: %v\n",
			repoName, repoPackages)
		cliutils.ExitWithError()
	}

	if !cliutils.YesNoPrompt("Do you want to continue?", false) {
		fmt.Println("❌ Repository Deletion Cancelled")
		cliutils.ExitWithError()
	}
	err := client.PackageRepositories().Delete(ctx, targetRepo, metav1.DeleteOptions{})
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

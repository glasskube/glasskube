package cmd

import (
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage package repositories in the current cluster",
}

func init() {
	repoCmd.AddCommand(repoListCmd, repoAddCmd, repoUpdateCmd)
	RootCmd.AddCommand(repoCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	RootCmd = cobra.Command{
		Use:     "glasskube",
		Version: "0.0.0",
		Short:   "Kubernetes Package Management the easy way 🔥",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("glasskube cli stub")
		},
	}
)

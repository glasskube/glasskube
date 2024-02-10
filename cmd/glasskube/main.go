package main

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/cmd/glasskube/cmd"
	"github.com/glasskube/glasskube/internal/cliutils"
)

func main() {
	cliutils.UpdateFetch()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

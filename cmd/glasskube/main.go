package main

import (
	"fmt"
	"os"

	"github.com/glasskube/glasskube/cmd/glasskube/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

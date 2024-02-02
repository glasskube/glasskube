package cliutils

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/glasskube/glasskube/pkg/list"
)

func PrintPackageTable(
	w io.Writer,
	packages []*list.PackageTeaserWithStatus,
	cols []string,
	getColsOfRow func(pkg *list.PackageTeaserWithStatus) []string,
) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	sep := "\t"
	fmt.Fprintf(tw, "%s\n", strings.Join(cols, sep))
	var sb strings.Builder
	for _, pkg := range packages {
		colsOfRow := getColsOfRow(pkg)
		for _, col := range colsOfRow {
			sb.WriteString(col)
			sb.WriteString(sep)
		}
		fmt.Fprintf(tw, "%s\n", sb.String())
		sb.Reset()
	}
	_ = tw.Flush()
}

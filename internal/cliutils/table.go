package cliutils

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/glasskube/glasskube/pkg/list"
)

var tabSep = "\t"

func PrintPackageTable(
	w io.Writer,
	packages []*list.PackageWithStatus,
	cols []string,
	getColsOfRow func(pkg *list.PackageWithStatus) []string,
) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(cols, tabSep))
	var sb strings.Builder
	for _, pkg := range packages {
		colsOfRow := getColsOfRow(pkg)
		if len(colsOfRow) != len(cols) {
			return fmt.Errorf("column mapping func returned %v columns instead of %v", len(colsOfRow), len(cols))
		}
		for _, col := range colsOfRow {
			sb.WriteString(col)
			sb.WriteString(tabSep)
		}
		fmt.Fprintln(tw, sb.String())
		sb.Reset()
	}
	return tw.Flush()
}

package cliutils

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const tabSep = "\t"

func PrintTable[T any](
	w io.Writer,
	rowItems []T,
	columns []string,
	getColsOfRow func(item T) []string,
) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, strings.Join(columns, tabSep)); err != nil {
		return err
	}
	var sb strings.Builder
	for _, pkg := range rowItems {
		colsOfRow := getColsOfRow(pkg)
		if len(colsOfRow) != len(columns) {
			return fmt.Errorf("column mapping func returned %v columns instead of %v", len(colsOfRow), len(columns))
		}
		for _, col := range colsOfRow {
			sb.WriteString(col)
			sb.WriteString(tabSep)
		}
		if _, err := fmt.Fprintln(tw, sb.String()); err != nil {
			return err
		}
		sb.Reset()
	}
	return tw.Flush()
}

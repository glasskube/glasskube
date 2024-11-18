package cliutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"golang.org/x/term"
)

var (
	bold      = color.New(color.Bold)
	italic    = color.New(color.Italic)
	underline = color.New(color.Underline)
	magenta   = color.New(color.FgMagenta)
)

type listContext struct {
	ordered bool
	index   int
}

type markdownRenderer struct {
	listContext *listContext
}

func MarkdownRenderer() renderer.NodeRenderer {
	return &markdownRenderer{}
}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *markdownRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)

	// inlines
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
}

func (*markdownRenderer) renderDocument(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderHeading(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		bold.SetWriter(writer)
	} else {
		bold.UnsetWriter(writer)
		if _, err := fmt.Fprintln(writer); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderBlockquote(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		italic.SetWriter(writer)
	} else {
		italic.UnsetWriter(writer)
	}
	return ast.WalkContinue, nil
}

func (r *markdownRenderer) renderCodeBlock(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if err := r.writeLines(writer, source, n); err != nil {
			return ast.WalkStop, err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (r *markdownRenderer) renderFencedCodeBlock(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if err := r.writeLines(writer, source, n); err != nil {
			return ast.WalkStop, err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderHTMLBlock(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *markdownRenderer) renderList(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.listContext = &listContext{
			ordered: n.(*ast.List).IsOrdered(),
		}
	} else {
		r.listContext = nil
		if _, err := fmt.Fprintln(writer); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (r *markdownRenderer) renderListItem(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	n.Parent().(*ast.List).IsOrdered()
	if entering {
		if r.listContext.ordered {
			if _, err := fmt.Fprintf(writer, " %v) ", r.listContext.index+1); err != nil {
				return ast.WalkStop, err
			}
		} else {
			if _, err := fmt.Fprint(writer, " * "); err != nil {
				return ast.WalkStop, err
			}
		}
		r.listContext.index++
	} else {
		if _, err := fmt.Fprintln(writer); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderParagraph(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if _, err := fmt.Fprint(writer, "\n\n"); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderTextBlock(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if n.NextSibling() != nil && n.FirstChild() != nil {
			if _, err := fmt.Fprintln(writer); err != nil {
				return ast.WalkStop, err
			}
		}
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderThematicBreak(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w, _, _ := term.GetSize(int(os.Stdout.Fd()))
		if w <= 0 {
			w = 40
		}
		_, _ = writer.WriteString(strings.Repeat("─", max(w, 40)) + "\n\n")
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderAutoLink(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	al := n.(*ast.AutoLink)
	if entering {
		underline.SetWriter(writer)
		if _, err := writer.Write(al.URL(source)); err != nil {
			return ast.WalkStop, err
		}
	} else {
		underline.UnsetWriter(writer)
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderCodeSpan(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		magenta.SetWriter(writer)
	} else {
		magenta.UnsetWriter(writer)
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderEmphasis(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	em := n.(*ast.Emphasis)
	style := italic
	if em.Level >= 2 {
		style = bold
	}
	if entering {
		style.SetWriter(writer)
	} else {
		style.UnsetWriter(writer)
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderImage(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderLink(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	l := n.(*ast.Link)
	if !entering {
		url := string(l.Destination)
		if !strings.Contains(url, "://") {
			url = "https://" + url
		}
		if _, err := fmt.Fprintf(writer, " (%v)", underline.Sprint(url)); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderRawHTML(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderText(
	writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		t := n.(*ast.Text)
		if _, err := writer.Write(t.Segment.Value(source)); err != nil {
			return ast.WalkStop, err
		}
		if t.HardLineBreak() {
			if _, err := fmt.Fprintln(writer); err != nil {
				return ast.WalkStop, err
			}
		}
		if t.SoftLineBreak() {
			_, _ = writer.WriteRune(' ')
		}
	}
	return ast.WalkContinue, nil
}

func (*markdownRenderer) renderString(
	writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if _, err := writer.Write(node.(*ast.String).Value); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (r *markdownRenderer) writeLines(w util.BufWriter, source []byte, n ast.Node) error {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		if _, err := w.Write(line.Value(source)); err != nil {
			return err
		}
	}
	return nil
}

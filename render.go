package markv

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/russross/blackfriday"
	"io"
	"regexp"
	"strings"
)

type Render struct{}

func (r *Render) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	w.Write([]byte(dfs(node, 0, 0)))
	return blackfriday.Terminate
}

func (r *Render) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	w.Write(ast.Literal)
}

func (r *Render) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	w.Write(ast.Literal)
}

func literal(node *blackfriday.Node) string {
	return string(node.Literal)
}

func oneline(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, "")
}

func dfs(node *blackfriday.Node, indent, index int) string {
	switch node.Type {
	case blackfriday.Document, blackfriday.Paragraph:
		rc := ""
		for c := node.FirstChild; c != nil; c = c.Next {
			rc += dfs(c, indent, index)
		}
		rc += "\n"
		return rc
	case blackfriday.Heading:
		m := []color.Attribute{
			color.FgRed, color.FgYellow, color.FgGreen,
			color.FgCyan, color.FgBlue, color.FgMagenta}
		rc := strings.Repeat(" ", node.Level-1)
		for c := node.FirstChild; c != nil; c = c.Next {
			rc += dfs(c, indent, index)
		}
		rc += "\n"
		return color.New(m[node.Level-1], color.Bold).Sprint(rc)
	case blackfriday.Text:
		return oneline(literal(node))
	case blackfriday.CodeBlock:
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		table := tablewriter.NewWriter(writer)

		table.Append([]string{literal(node)})
		table.Render()
		writer.Flush()

		rc := ""
		rc += b.String()
		return color.New(color.FgHiYellow).Sprint(rc)
	case blackfriday.Table:
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		table := tablewriter.NewWriter(writer)

		headNode := node.FirstChild
		var keys []string
		for c := headNode.FirstChild.FirstChild; c != nil; c = c.Next {
			keys = append(keys, dfs(c, indent, index))
		}
		table.SetHeader(keys)

		bodyNode := node.LastChild
		for r := bodyNode.FirstChild; r != nil; r = r.Next {
			var row []string
			for c := r.FirstChild; c != nil; c = c.Next {
				row = append(row, dfs(c, indent, index))
			}
			table.Append(row)
		}

		table.Render()
		writer.Flush()
		rc := b.String()
		return rc
	case blackfriday.HorizontalRule:
		return color.New(color.FgGreen, color.Bold).Sprint(
			"♦━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━♦\n")
	case blackfriday.Emph:
		return color.New(color.Italic).Sprint(dfs(node.FirstChild, indent, index))
	case blackfriday.Strong:
		return color.New(color.Bold).Sprint(dfs(node.FirstChild, indent, index))
	case blackfriday.BlockQuote:
		// TODO: check newline ascii code
		return "\n⦀ " + dfs(node.FirstChild, indent, index) + "\n"
	case blackfriday.List:
		// TODO: we must consider node.ListFlags.
		rc := ""
		i := 0
		for c := node.FirstChild; c != nil; c = c.Next {
			rc += dfs(c, indent+1, i)
			i++
		}
		return rc
	case blackfriday.Item:
		// TODO: we must consider indent.
		rc := strings.Repeat("  ", indent)
		if node.ListFlags&blackfriday.ListTypeOrdered != 0 {
			rc += color.New(color.FgRed).Sprintf("%d. ", index+1)
		} else {
			rc += color.New(color.FgRed).Sprint("- ")
		}
		for c := node.FirstChild; c != nil; c = c.Next {
			rc += dfs(c, indent, index)
		}
		return rc
	case blackfriday.TableCell:
		rc := dfs(node.FirstChild, 0, 0)
		return rc
	case blackfriday.HTMLBlock:
		return color.New(color.FgHiGreen).Sprint(literal(node) + "\n")
	case blackfriday.HTMLSpan:
		return color.New(color.FgHiGreen).Sprint(literal(node))
	case blackfriday.Link:
		return color.New(color.FgHiMagenta, color.Underline).Sprint(dfs(node.FirstChild, indent, index))
	case blackfriday.Image:
		return "\n" + RenderImage(string(node.LinkData.Destination))
	case blackfriday.Code:
		return color.New(color.FgYellow, color.Bold).Sprint(literal(node))
	}
	return fmt.Sprintf("Unknown type %#v is found.\n", node.Type)
}

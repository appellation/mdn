package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func parseContent(rd io.Reader, build *strings.Builder) error {
	doc, err := html.Parse(rd)
	if err != nil {
		return err
	}

	parseNode(doc, build, false)
	return nil
}

func parseNode(n *html.Node, build *strings.Builder, append bool) {
	if n == nil {
		return
	}

	switch n.Type {
	case html.DoctypeNode, html.DocumentNode:
		iterateNode(n, build, append)
	case html.ElementNode:
	attrItr:
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				switch attr.Val {
				case "hidden", "bc-data":
					append = false
					break attrItr
				}
			}
		}

		switch n.DataAtom {
		case atom.Article:
			append = true
		case atom.Table, atom.Pre:
			append = false
		}

		iterateNode(n, build, append)
	case html.TextNode:
		d := strings.TrimSpace(n.Data)
		if d != "" && append {
			build.WriteString(d + " ")
		}
	}
}

func iterateNode(n *html.Node, build *strings.Builder, append bool) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c, build, append)
	}
}

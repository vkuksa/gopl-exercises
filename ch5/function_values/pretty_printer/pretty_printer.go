// Exercise 5.7 from gopl.io

package pretty_printer

import (
	"bytes"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

var depth int
var buf bytes.Buffer

func Outline(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	forEachNode(doc, startElement, endElement)

	return buf.Bytes(), nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		attr := fmt.Sprintf("%*s<%s", depth*2, "", n.Data)

		for _, a := range n.Attr {
			attr += string(" " + a.Key + "=\"" + a.Val + "\"")
		}

		if n.FirstChild != nil {
			attr += ">\n"
			depth++
		} else {
			attr += "/>\n"
		}

		buf.WriteString(attr)
	} else if n.Type == html.CommentNode {
		buf.WriteString(fmt.Sprintf("%*s<!--%s-->\n", (depth+1)*2, "", n.Data))
	} else if n.Type == html.TextNode {
		buf.WriteString(fmt.Sprintf("%*s%s\n", (depth+1)*2, "", n.Data))
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		if n.FirstChild != nil {
			depth--
			buf.WriteString(fmt.Sprintf("%*s</%s>\n", depth*2, "", n.Data))
		}
	}
}

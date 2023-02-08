// Exercise 7.4: The strings.NewReader function returns a value that satisfies the io.Reader
// interface (and others) by reading from its argument, a string. Implement a simple version of
// NewReader yourself, and use it to make the HTML parser (§5.2) take input from a string.

package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/html"
)

const (
	sample = ` <!DOCTYPE html>
	<html>
	<body>
	
	<h1>My First Heading</h1>
	<p>My first paragraph.</p>
	
	</body>
	</html> `
)

type StringReader struct {
	s string
	i int64 // current reading index
}

func (r *StringReader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return n, nil
}

func NewReader(s string) *StringReader {
	return &StringReader{s, 0}
}

func main() {
	reader := NewReader(sample)

	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var depth int
	var traverseSubnodes func(n *html.Node)
	traverseSubnodes = func(n *html.Node) {
		depth++
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverseSubnodes(c)
		}

		if n.Type == html.ElementNode {
			fmt.Printf("%*s<!--%s-->\n", (depth+1)*2, "", n.Data)
		}

		depth--
	}
	traverseSubnodes(doc)
}

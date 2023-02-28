// Exercise 7.17: Extend xmlselect so that elements may be selected not just by name, but by
// their attributes too, in the manner of CSS, so that, for instance, an element like
// <div id="page" class="wide"> could be selected by a matching id or class as well as its
// name.

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type match struct {
	name  *xml.Name
	attrs *[]xml.Attr
}

func main() {
	dec := xml.NewDecoder(os.Stdin)
	var stack []match // stack of element names
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, match{&tok.Name, &tok.Attr}) // push

		case xml.EndElement:
			stack = stack[:len(stack)-1] // pop
		case xml.CharData:
			if matchesInput(stack, os.Args[1:]) {
				printStack(stack, tok)
			}
		}
	}
}

func matchesInput(x []match, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0].name.Local == y[0] {
			y = y[1:]
		} else if len(*x[0].attrs) > 0 {
			for _, a := range *x[0].attrs {
				if a.Value == y[0] {
					y = y[1:]
				}
			}
		}
		x = x[1:]
	}
	return false
}

func printStack(stack []match, tok xml.Token) {
	var matches string
	for _, m := range stack {
		matches += ("<" + m.name.Local)
		for _, attr := range *m.attrs {
			matches += fmt.Sprintf(" %s=%s", attr.Name.Local, attr.Value)
		}
		matches += ">\t"
	}
	fmt.Printf("%s: %s\n", matches, tok)

}

//!-

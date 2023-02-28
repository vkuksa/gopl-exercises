// Exercise 7.18: Using the token-based decoder API, write a program that will read an arbitrary
// XML document and construct a tree of generic nodes that represents it. Nodes are of two
// kinds: CharData nodes represent text strings, and Element nodes represent named elements
// and their attributes. Each element node has a slice of child nodes.

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Node interface{} // CharData or *Element

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func (el Element) String() string {
	result := el.Type.Local
	for _, attr := range el.Attr {
		result += fmt.Sprintf(" %s=%s", attr.Name.Local, attr.Value)
	}
	return result
}

func main() {
	root := Element{xml.Name{}, nil, nil}
	stack := []*Element{&root}

	dec := xml.NewDecoder(os.Stdout)
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlparser: %v\n", err)
			os.Exit(1)
		}
		switch tok := t.(type) {
		case xml.StartElement:
			elem := Element{tok.Name, tok.Copy().Attr, make([]Node, 0)}

			// Add to a children collection of last stack element
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, &elem)

			// Make element itself as last element of stack
			stack = append(stack, &elem)
		case xml.EndElement:
			elem := Element{tok.Name, nil, nil}

			// Remove StartElement pair of this element from stack
			stack = stack[:len(stack)-1]

			// Add to a children collection of last stack element
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, &elem)
		case xml.CharData, xml.Comment:
			elem := CharData(string(t.(xml.CharData)))

			// Add to a children collection of last stack element
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, elem)
		case xml.ProcInst, xml.Directive:
			continue
		}
	}

	printTree(&root)
}

func printTree(root Node) {
	depth := 0
	var traverseAndPrintNode func(el Node)
	traverseAndPrintNode = func(el Node) {
		switch el := el.(type) {
		case CharData:
			fmt.Printf("%*s%s\n", depth*2, "", string(el))
		case *Element:
			fmt.Printf("%*s<%s>\n", depth*2, "", el)

			depth++
			for _, child := range el.Children {
				traverseAndPrintNode(child)
			}
			depth--
		default:
			panic("Invalid element type given for tree print")
		}
	}

	for _, child := range root.(*Element).Children {
		traverseAndPrintNode(child)
	}
}

//!-

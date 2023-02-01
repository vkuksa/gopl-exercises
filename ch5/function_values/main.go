// Exercises 5.8, 5.9 from gopl.io

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func expand(s string, f func(string) string) string {
	fmt.Println(s)
	return strings.ReplaceAll(s, "‘‘$foo’’", f("foo"))
}

func getModifiers() []func(str string) string {
	return []func(str string) string{
		func(str string) string {
			return str + "1"
		},
		func(str string) string {
			return str + "2"
		},
		func(str string) string {
			return str + "3"
		},
	}
}

func main() {
	for _, f := range getModifiers() {
		fmt.Println(expand(os.Args[1], f))
	}
}

func printElementById(url, id string) {
	resp, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	element := ElementByID(doc, os.Args[2])
	fmt.Println(element)
}

func ElementByID(doc *html.Node, id string) *html.Node {
	var result *html.Node = nil
	findElement := func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					result = n
					return false
				}
			}
		}

		return true
	}

	forEachNode(doc, findElement, findElement)

	return result
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) bool {
	if !pre(n) {
		return false
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		proceed := forEachNode(c, pre, post)

		if !proceed {
			return false
		}
	}

	return post(n)
}

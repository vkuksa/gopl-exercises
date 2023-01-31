// Exercise 5.1, 5.2, 5.3, 5.4 from gopl.io

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func parseHtmlFrom(url string) *html.Node {
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(err)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func visit(n *html.Node, callback func(*html.Node)) {
	if n == nil {
		return
	}

	visit(n.FirstChild, callback)
	callback(n)
	visit(n.NextSibling, callback)
}

func main() {
	doc := parseHtmlFrom(os.Args[1])

	// 5.1
	var links []string
	findLinks := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		}
	}
	visit(doc, findLinks)

	for _, link := range links {
		fmt.Println(link)
	}

	// 5.2
	elements := make(map[string]int)
	countElements := func(n *html.Node) {
		if n.Type == html.ElementNode {
			elements[n.Data]++
		}
	}
	visit(doc, countElements)
	fmt.Println(elements)

	// 5.3
	printContent := func(n *html.Node) {
		if n.Type == html.TextNode && n.Parent.Data != "script" && n.Parent.Data != "style" {
			fmt.Println(n.Data)
		}
	}
	visit(doc, printContent)

	// 5.4
	findLinksExt := func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						fmt.Println(a.Val)
					}
				}
			} else if n.Data == "img" ||
				n.Data == "script" ||
				n.Data == "audio" ||
				n.Data == "embed" ||
				n.Data == "iframe" ||
				n.Data == "input" ||
				n.Data == "source" ||
				n.Data == "track" ||
				n.Data == "video" {
				for _, a := range n.Attr {
					if a.Key == "src" {
						fmt.Println(a.Val)
					}
				}
			}
		}

	}
	visit(doc, findLinksExt)
}

// Exercises 5.15, 5.16, 5.17 from gopl.io

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func joinVariadic(sep string, vals ...string) string {
	if len(vals) < 1 {
		return ""
	}
	result := vals[0]
	for i := 1; i < len(vals); i++ {
		result += (sep + vals[i])
	}
	return result
}

func max(vals ...int) int {
	if len(vals) < 1 {
		return MaxInt
	}
	max := vals[0]
	for i := 1; i < len(vals); i++ {
		if vals[i] > max {
			max = vals[i]
		}
	}

	return max
}

func min(vals ...int) int {
	if len(vals) < 1 {
		return MinInt
	}
	min := vals[0]
	for i := 1; i < len(vals); i++ {
		if vals[i] < min {
			min = vals[i]
		}
	}

	return min
}

func maxSafe(val int, vals ...int) int {
	max := val
	for _, val := range vals {
		if val > max {
			max = val
		}
	}

	return max
}

func minSafe(val int, vals ...int) int {
	min := val
	for _, val := range vals {
		if val < min {
			min = val
		}
	}

	return min
}

func ElementsByTagName(doc *html.Node, names ...string) []*html.Node {
	elements := make([]*html.Node, 0)

	collect := func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, name := range names {
				if n.Data == name {
					elements = append(elements, n)
				}
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, collect)
	}
	return elements
}

func forEachNode(n *html.Node, apply func(n *html.Node)) {
	apply(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, apply)
	}
}

func main() {
	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatalf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalf("parsing %s as HTML: %v", url, err)
	}

	headings := ElementsByTagName(doc, "h1", "h2", "h3", "h4")
	for _, heading := range headings {
		fmt.Println(heading.Data)
	}
}

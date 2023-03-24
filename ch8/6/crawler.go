// Exercise 8.6: Add depth-limiting to the concurrent crawler. That is, if the user sets -depth=3,
// then only URLs reachable by at most three links will be fetched.

package main

import (
	"flag"
	"fmt"
	"gopl-exercises/ch5/links"
	"log"
	"math"
	"os"
)

type Link struct {
	url   string
	depth int
}
type Links struct {
	urls  []string
	depth int
}

func crawl(url string) (results []string) {
	fmt.Println(url)
	results, err := links.Extract(url, make(chan struct{}))
	if err != nil {
		log.Print(err)
	}

	return
}

var maxDepth = flag.Int("depth", math.MaxInt, "depth of crawl limit")

func init() {
	flag.Parse()
}

// !+
func main() {
	worklist := make(chan Links)   // lists of URLs, may have duplicates
	unseenLinks := make(chan Link) // de-duplicated URLs

	// Init list with command line arguments
	go func() { worklist <- Links{os.Args[1:], 0} }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link.url)
				go func(depth int) { worklist <- Links{foundLinks, depth} }(link.depth + 1)
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		if list.depth <= *maxDepth {
			for _, link := range list.urls {
				if !seen[link] {
					seen[link] = true
					unseenLinks <- Link{link, list.depth}
				}
			}
		}
	}
}

//!-

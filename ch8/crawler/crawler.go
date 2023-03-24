// Exercise 8.10: HTTP requests may be cancelled by closing the optional Cancel channel in the
// http.Request struct. Modify the web crawler of Section 8.6 to support cancellation.

package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopl-exercises/ch5/links"
)

var done = make(chan struct{})

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url, done)
	if err != nil {
		log.Print(err)
	}
	return list
}

// !+
func main() {
	wg := &sync.WaitGroup{}

	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add command-line arguments to worklist.
	wg.Add(1)
	go func() {
		worklist <- os.Args[1:]
		wg.Done()
	}()

	wg.Add(20)
	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case link := <-unseenLinks:
					foundLinks := crawl(link)

					wg.Add(1)
					go func() {
						defer wg.Done()
						select {
						case <-done:
							return
						case worklist <- foundLinks:
						}
					}()
				}
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	wg.Add(1)
	go func() {
		defer wg.Done()
		seen := make(map[string]bool)
		for {
			select {
			case <-done:
				return
			case list := <-worklist:
				for _, link := range list {
					if !seen[link] {
						seen[link] = true

						wg.Add(1)
						go func(link string) {
							defer wg.Done()

							select {
							case <-done:
								return
							case unseenLinks <- link:
							}
						}(link)
					}
				}
			}
		}
	}()

	os.Stdin.Read(make([]byte, 1)) // read a single byte
	close(done)
	wg.Wait()
}

//!-

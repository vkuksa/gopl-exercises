// Exercise 8.7: Write a concurrent program that creates a local mirror of a web site, fetching
// each reachable page and writing it to a directory on the local disk. Only pages within the
// original domain (for instance, golang.org) should be fetched. URLs within mirrored pages
// should be altered as needed so that they refer to the mirrored page, not the original

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"

	"golang.org/x/net/html"
)

const (
	MaxWorkers = 8
)

type DocInfo struct {
	root *html.Node
	path string
}

func scrapHtmlTree(url string) (doc *html.Node, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Returned non OK status")
		return
	}

	doc, err = html.Parse(resp.Body)
	return
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
	var wg sync.WaitGroup

	urlsToScrap := make(chan *url.URL)
	docsToWrite := make(chan *DocInfo)

	seen := make(map[string]bool)
	var seenMutex sync.Mutex

	var workingHostname, saveRoot string
	if url, err := url.Parse(os.Args[1]); err == nil {
		workingHostname = url.Hostname()
		log.Println("working over hostname: " + workingHostname)

		saveRoot = filepath.Join(os.TempDir(), "mirror", workingHostname)
		if err := os.MkdirAll(saveRoot, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
			log.Fatalf("%s save root folder creation failed: %s", saveRoot, err.Error())
		}
		log.Println("save root folder created: " + saveRoot)

		wg.Add(1)
		go func() {
			urlsToScrap <- url
		}()
	} else {
		log.Fatalf("Couldn't extract hostname from %s: %s", os.Args[1], err.Error())
	}

	// Scrapping docs
	for i := 0; i < MaxWorkers; i++ {
		go func() {
			for u := range urlsToScrap {
				func(url *url.URL) {
					defer func() {
						wg.Done()
					}()

					seenMutex.Lock()
					if seen[url.Path] {
						log.Println("scrapping was already done: " + url.String())
						seenMutex.Unlock()
						return
					}
					seen[url.Path] = true
					seenMutex.Unlock()

					log.Println("scrapping: " + url.String())

					dom, err := scrapHtmlTree(url.String())
					if err != nil {
						log.Println("scrapping " + url.String() + ": " + err.Error())
						return
					}

					findLinks := func(n *html.Node) {
						if n.Type == html.ElementNode && n.Data == "a" {
							for _, a := range n.Attr {
								if a.Key == "href" {
									// Ignorring errors from parsing hostname
									if url, err := url.Parse(a.Val); err == nil && url.Hostname() == workingHostname {
										wg.Add(1)
										go func() {
											urlsToScrap <- url
										}()
									}
								}
							}
						}
					}
					log.Println("processing dom: " + url.String())
					visit(dom, findLinks)

					wg.Add(1)
					go func() {
						docsToWrite <- &DocInfo{dom, url.Path}
					}()
				}(u)
			}
		}()
	}

	for i := 0; i < MaxWorkers; i++ {
		go func() {
			for d := range docsToWrite {
				func(doc *DocInfo) {
					defer wg.Done()

					docDir, docFile := path.Split(doc.path)
					saveDir := saveRoot + docDir
					saveFilename := saveDir + docFile
					if docFile == "" {
						saveFilename += "index.html"
					}
					log.Println("saving file: " + saveFilename)

					if err := os.Mkdir(saveDir, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
						log.Printf("%s folder creation failed: %s", saveDir, err.Error())
						return
					}

					file, err := os.OpenFile(saveFilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
					if err != nil && !errors.Is(err, os.ErrExist) {
						log.Printf("%s file creation failed: %s", saveFilename, err.Error())
						return
					}

					html.Render(file, doc.root)
					file.Close()
				}(d)
			}
		}()
	}

	wg.Wait()
}

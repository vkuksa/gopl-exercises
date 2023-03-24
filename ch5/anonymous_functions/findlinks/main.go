// Exercise 5.13, 5.14 from gopl.io

package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopl-exercises/ch5/links"
)

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func retrievePageContent(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Couldn't retrieve page content")
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	return buf.Bytes(), nil
}

func writeDataIntoFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	file.Write(data)
	file.Close()
	return nil
}

func makeLocalCopiesOf(links []string) {
	tempdir := filepath.Join(os.TempDir(), "/findlinks/")
	err := os.Mkdir(tempdir, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Printf("%s folder creation failed: %s", tempdir, err.Error())
		return
	}

	for _, link := range links {
		data, err := retrievePageContent(link)
		if err != nil {
			log.Printf("%s page content retrieval: %s", tempdir, err.Error())
			continue
		}

		u, err := url.Parse(link)
		hostname := strings.TrimPrefix(u.Hostname(), "www.")
		filename := tempdir + "/" + hostname + "/index.html"
		err = os.Mkdir(tempdir+"/"+hostname, os.ModePerm)
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Printf("%s folder creation failed: %s", tempdir, err.Error())
			continue
		}

		err = writeDataIntoFile(filename, data)
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Printf("opening file %s failed: %s", filename, err.Error())
			continue
		}

	}
}

func crawl(link string) []string {
	fmt.Println(link)
	list, err := links.Extract(link, make(chan struct{}))
	if err != nil {
		log.Print(err)
	}

	makeLocalCopiesOf(list)

	return list
}

func main() {
	breadthFirst(crawl, os.Args[1:])
}

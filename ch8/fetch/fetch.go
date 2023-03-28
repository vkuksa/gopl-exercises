// Exercise 8.11: Following the approach of mirroredQuery in Section 8.4.4, implement a vari-
// ant of fetch that requests several URLs concurrently. As soon as the first response arrives,
// cancel the other requests.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// !+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string, cancel chan struct{}) (filename string, written int64, err error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	req.Cancel = cancel

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	var f *os.File
	defer func() {
		resp.Body.Close()

		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()

	filename = path.Base(resp.Request.URL.Path)
	if filename == "/" {
		filename = "index.html"
	}
	f, err = os.Create(filename)
	if err != nil {
		return
	}
	written, err = io.Copy(f, resp.Body)
	return
}

//!-

func main() {
	cancel := make(chan struct{})
	received := make(chan struct{})

	for _, url := range os.Args[1:] {
		go func(url string) {
			filename, written, err := fetch(url, cancel)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			} else {
				close(cancel)
				fmt.Printf("%s => %s (%d bytes).\n", url, filename, written)
				received <- struct{}{}
			}
		}(url)
	}

	select {
	case <-time.After(10 * time.Second):
		close(cancel)
		fmt.Println("Polling time exceeded")
	case <-received:
		os.Exit(0)
	}

	os.Exit(1)
}

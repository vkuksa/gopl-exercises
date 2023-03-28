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

var done = make(chan struct{})

type FetchInfo struct {
	url      string
	filename string
	n        int64
	err      error
}

// !+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) FetchInfo {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	req.Cancel = done

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return FetchInfo{url, "", 0, err}
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return FetchInfo{url, "", 0, err}
	}
	n, err := io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return FetchInfo{url, local, n, err}
}

//!-

func main() {
	response := make(chan FetchInfo)

	for _, url := range os.Args[1:] {
		go func(url string) {
			info := fetch(url)
			if info.err != nil {
				fmt.Fprintf(os.Stderr, "fetch %s: %v\n", info.url, info.err)
			} else {
				response <- info
			}
		}(url)
	}

	select {
	case info := <-response:
		fmt.Printf("%s => %s (%d bytes).\n", info.url, info.filename, info.n)
	case <-time.After(10 * time.Second):
		fmt.Println("Polling time exceeded")
	}

	close(done)
}

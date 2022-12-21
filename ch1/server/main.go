// Exercise 1.12 from gopl.io

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	drawer "gobook/ch1/lissajous"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var size int
		// for _, parameter range r.URL.Query() {
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			size = 100
		}

		fmt.Printf("Size: %d", size)
		drawer.Lissajous(w, size)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// !+handler
// handler echoes the HTTP request.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

//!-handler

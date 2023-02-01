package pretty_printer_test

import (
	"bytes"
	"testing"

	. "gopl-exercises/ch5/function_values/pretty_printer"

	"golang.org/x/net/html"
)

func urlsList() []string {
	return []string{"https://google.com", "https://golang.org", "https://youtube.com"}
}

func TestDirectMatch(t *testing.T) {
	for _, url := range urlsList() {
		data, err := Outline(url)
		if err != nil {
			t.Error(err)
		}

		reader := bytes.NewReader(data)
		_, err = html.Parse(reader)
		if err != nil {
			t.Error(err)
		}
	}
}

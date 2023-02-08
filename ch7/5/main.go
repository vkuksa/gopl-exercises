// Exercise 7.5: The LimitReader function in the io package accepts an io.Reader r and a
// number of bytes n, and returns another Reader that reads from r but reports an end-of-file
// condition after n bytes. Implement it.
// func LimitReader(r io.Reader, n int64) io.Reader

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	sample = "Lorem ipsum dolor"
)

type LimitedReader struct {
	r io.Reader
	n int64
}

func (l *LimitedReader) Read(b []byte) (n int, err error) {
	if l.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(b)) > l.n {
		b = b[0:l.n]
	}
	n, err = l.r.Read(b)
	l.n -= int64(n)

	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitedReader{r, n}
}

func main() {
	reader := strings.NewReader(sample)
	limitReader := LimitReader(reader, 5)
	bytes, err := io.Copy(os.Stdout, limitReader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nThe number of bytes: %d\n", bytes)
}

// Exercise 7.2: Write a function CountingWriter with the signature below that, given an
// io.Writer, returns a new Writer that wraps the original, and a pointer to an int64 variable
// that at any moment contains the number of bytes written to the new Writer.
//
//	func CountingWriter(w io.Writer) (io.Writer, *int64)
package main

import (
	"fmt"
	"io"
	"os"
)

type ByteCounter struct {
	wrapped io.Writer
	bytes   *int64
}

func (c *ByteCounter) Bytes() int64 {
	return *c.bytes
}

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c.bytes += int64(len(p))
	return c.wrapped.Write(p)
}

func CountingWriter(wrapped io.Writer) (io.Writer, *int64) {
	writer := &ByteCounter{wrapped: wrapped, bytes: new(int64)}
	*writer.bytes = 0
	return writer, writer.bytes
}

func main() {
	f, err := os.OpenFile("input.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var writer, bytes = CountingWriter(f)

	var name = "Dolly"
	fmt.Fprintf(writer, "hello, %s", name)
	fmt.Println(*bytes)
}

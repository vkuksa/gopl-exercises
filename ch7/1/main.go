// Exercise 7.1: Using the ideas from ByteCounter, implement counters for words and for lines.
// You will find bufio.ScanWords useful.
package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type WordCounter int
type LineCounter int
type ByteCounter int

func tokenCount(p []byte, splitter bufio.SplitFunc) int {
	reader := bytes.NewReader(p)
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitter)

	token := 0
	for scanner.Scan() {
		token++
	}

	return token
}

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

func (c *WordCounter) Write(p []byte) (int, error) {
	words := tokenCount(p, bufio.ScanWords)
	*c += WordCounter(words)
	return words, nil
}

func (c *LineCounter) Write(p []byte) (int, error) {
	lines := tokenCount(p, bufio.ScanLines)
	*c += LineCounter(lines)
	return lines, nil
}

const smallInput = `first second third fourth`
const biggerInput = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Viverra justo nec ultrices dui sapien. Magna etiam tempor orci eu lobortis elementum. 
Arcu felis bibendum ut tristique et egestas quis. Donec ac odio tempor orci dapibus ultrices in. 
Est ullamcorper eget nulla facilisi etiam dignissim. 
Cras tincidunt lobortis feugiat vivamus at augue eget arcu dictum. 
Porttitor massa id neque aliquam. Sed augue lacus viverra vitae congue eu consequat. 
Enim tortor at auctor urna nunc. Vitae justo eget magna fermentum iaculis eu non. 
Proin nibh nisl condimentum id. Et molestie ac feugiat sed lectus vestibulum mattis. 
At varius vel pharetra vel turpis nunc eget lorem dolor. 
Nisl suscipit adipiscing bibendum est ultricies integer quis. 
Volutpat blandit aliquam etiam erat.`

func main() {
	var bc ByteCounter
	var wc WordCounter
	var lc LineCounter

	fmt.Fprintf(&bc, smallInput)
	fmt.Printf("Bytes: %d\n", bc)

	fmt.Fprintf(&wc, smallInput)
	fmt.Printf("Words: %d\n", wc)

	fmt.Fprintf(&lc, smallInput)
	fmt.Printf("Lines: %d\n", lc)

	bc, wc, lc = 0, 0, 0

	fmt.Fprintf(&bc, biggerInput)
	fmt.Printf("Bytes: %d\n", bc)

	fmt.Fprintf(&wc, biggerInput)
	fmt.Printf("Words: %d\n", wc)

	fmt.Fprintf(&lc, biggerInput)
	fmt.Printf("Lines: %d\n", lc)
}

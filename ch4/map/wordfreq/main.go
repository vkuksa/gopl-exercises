// Exercises 4.9 from gopl.io

// Write a program wordfreq to report the frequency of each word in an input text
// file. Call input.Split(bufio.ScanWords) before the first call to Scan to break the input into
// words instead of lines.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	words := make(map[string]int) // counts word occurences

	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	input := bufio.NewScanner(file)
	input.Split(bufio.ScanWords)

	for input.Scan() {
		words[input.Text()]++
	}

	if err := input.Err(); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("word\tcount\n")
	for word, count := range words {
		fmt.Printf("%s\t\t\t\t%d\n", word, count)
	}
}

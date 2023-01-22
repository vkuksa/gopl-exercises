// Exercises 3.10, 3.11, 3.12 from gopl.io

package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	delimSize = 3
)

func main() {
	for i := 1; i < len(os.Args); i += 2 {
		fmt.Printf("  %t\n", anagram(os.Args[i], os.Args[i+1]))
	}
}

func comma(s string) string {
	var buf bytes.Buffer

	if len(s) == 0 {
		return buf.String()
	}

	// inputSize - size of chunk, which we apply comma delimiters to (everything before the dot in case of float)
	start, integerSize := 0, len(s)

	// We are interested in working with digits of integer, discarding fractional part and appending it later
	dotIndex := strings.IndexByte(s, '.')
	if dotIndex != -1 {
		integerSize -= (len(s) - dotIndex)
	}

	if s[0] == '-' || s[0] == '+' {
		buf.WriteByte(s[0])
		integerSize--
		start = 1
	}

	// Check whether we need to insert comma delimiter
	if integerSize > delimSize {
		end := start
		// Calculate the position of first comma
		if integerSize%delimSize != 0 {
			end += integerSize % delimSize
		} else {
			// If provided number has 6, 9, 12 etc digits - it's mod 3 equals to 0 and
			// We should start from the next comma
			end += delimSize
		}

		for ; end < integerSize; start, end = end, end+delimSize {
			buf.WriteString(s[start:end])
			buf.WriteByte(',')
		}
	}

	buf.WriteString(s[start:])

	return buf.String()
}

func Sort(input string) []rune {
	sorted := []rune(input)
	sort.Slice(sorted, func(i int, j int) bool { return sorted[i] < sorted[j] })
	return sorted
}

func Equal(first, second []rune) bool {
	if len(first) != len(second) {
		return false
	}
	for i, v := range first {
		if v != second[i] {
			return false
		}
	}
	return true
}

// No performance race here :)
// Counts space as distinct letter: "New York Times" != "monkeys write" (has 1 more space)
func anagram(first, second string) bool {
	firstSorted, secondSorted := Sort(first), Sort(second)
	return Equal(firstSorted, secondSorted)
}

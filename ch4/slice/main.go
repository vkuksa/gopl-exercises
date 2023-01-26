// Exercises 4.3, 4.4, 4.5, 4.6, 4.7 from gopl.io

package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const (
	arraySize = 4
)

func main() {
	data := [arraySize]int{1, 2, 3, 4}
	fmt.Printf("%d\n", data)
	reverseArray(&data)
	fmt.Printf("%d\n", data)

	sliced := []int{0, 1, 2, 3, 4, 5}
	fmt.Printf("%d\n", sliced)
	rotate(sliced, 8)
	fmt.Printf("%d\n", sliced) //"[2 3 4 5 0 1]"

	adjust := []string{"a", "b", "b", "a", "c", "a"}
	fmt.Printf("%s\n", adjust)
	adjust = eliminateAdjacentDuplicates(adjust)
	fmt.Printf("%s\n", adjust) //"c a"

	squash := []byte(string("Ḽơᶉëᶆ   ȋṕšᶙṁ ḍỡḽǭᵳ ʂǐť ӓṁệẗ.    |"))
	fmt.Printf("%s\n", squash)
	squash = squashSpaces(squash) // "Ḽơᶉëᶆ ȋṕšᶙṁ ḍỡḽǭᵳ ʂǐť ӓṁệẗ. |"
	fmt.Printf("%s\n", squash)

	reverse := []byte(string("Ḽơᶉëᶆ ȋṕšᶙṁ ḍỡḽǭᵳ ʂǐť ӓṁệẗ."))
	fmt.Printf("%s\n", reverse)
	reverse = reverseInPlace(reverse) // ".ẗệṁӓ ťǐʂ ᵳǭḽỡḍ ṁᶙšṕȋ ᶆëᶉơḼ"
	fmt.Printf("%s\n", reverse)
}

func reverseArray(s *[arraySize]int) {
	for i, j := 0, arraySize-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

func rotate(s []int, positions int) {
	positions = positions % len(s)
	buffer := make([]int, positions)

	for i, j := 0, 0; i < len(s); i++ {
		if i < positions {
			buffer[i] = s[i]
		} else {
			s[i-positions] = s[i]
		}

		if i >= len(s)-positions {
			s[i] = buffer[j]
			j++
		}
	}
}

func eliminateAdjacentDuplicates(input []string) []string {

	removals := 0
	for cleaned := false; !cleaned; {
		cleaned = true

		for i := 1; i < len(input)-1; i++ {
			if input[i-1] == input[i] {
				copy(input[i-1:], input[i+1:])
				removals += 2
				cleaned = false
			}
		}
	}

	return input[:len(input)-removals]
}

func squashSpaces(input []byte) []byte {
	removals := 0
	for i := 0; i < len(input)-removals; {
		r, size := utf8.DecodeRune(input[i:])
		if unicode.IsSpace(r) {
			shouldSquash, squashSize := false, 1
			for nextr, nextsize := utf8.DecodeRune(input[i+size : len(input)-removals]); unicode.IsSpace(nextr); {
				shouldSquash = true
				squashSize += nextsize
				nextr, nextsize = utf8.DecodeRune(input[i+squashSize : len(input)-removals])
			}

			if shouldSquash {
				copy(input[i+1:], input[i+squashSize:])
				input[i] = ' '
				size = 1
				removals += (squashSize - 1)
				shouldSquash, squashSize = false, 1
			}
		}

		i += size
	}

	return input[:len(input)-removals]
}

func reverseInPlace(input []byte) []byte {
	for i, j := 0, len(input); i < j; {
		frune, fsize := utf8.DecodeRune(input[i:j])
		lrune, lsize := utf8.DecodeLastRune(input[i:j])

		if fsize != lsize {
			copy(input[i+lsize:j-fsize], input[i+fsize:j-lsize])
		}

		utf8.EncodeRune(input[i:], lrune)
		utf8.EncodeRune(input[j-fsize:], frune)

		i += lsize
		j -= fsize
	}

	return input
}

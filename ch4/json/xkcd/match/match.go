package match

import (
	"math"
	"strings"
)

func levenshteinDistance(a, b string) int {
	if a == b {
		return 0
	}

	if len(a) > len(b) {
		a, b = b, a
	}
	lenA := len(a)
	lenB := len(b)

	var x []byte = make([]byte, lenA+1)

	for i := 1; i < len(x); i++ {
		x[i] = byte(i)
	}

	for i := 1; i <= lenB; i++ {
		prev := byte(i)
		for j := 1; j <= lenA; j++ {
			current := x[j-1]
			if b[i-1] != a[j-1] {
				current = min(min(x[j-1]+1, prev+1), x[j]+1)
			}
			x[j-1] = prev
			prev = current
		}
		x[lenA] = prev
	}
	return int(x[lenA])
}

func smallestEditDistance(a, b string) int {
	wordsA := strings.Fields(a)
	wordsB := strings.Fields(b)

	min, distance := math.MaxInt32, 0
	for _, wordA := range wordsA {
		for _, wordB := range wordsB {
			distance = levenshteinDistance(wordA, wordB)
			if distance < min {
				min = distance
			}
		}
	}

	return min
}

func min(a, b byte) byte {
	if a < b {
		return a
	}
	return b
}

func DirectMatch(string, match_against string) bool {
	return strings.Contains(string, match_against)
}

func FuzzyMatch(string, match_against string, threshold int) bool {
	return smallestEditDistance(string, match_against) < threshold
}

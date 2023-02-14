package palindrome

import (
	"sort"
	"testing"
)

func TestPalindrome(t *testing.T) {
	strings1 := sort.StringSlice{"a", "b", "cd", "x", "cd", "b", "a"}
	strings2 := sort.StringSlice{"a", "b", "cd", "x", "dc", "b", "a"}
	ints1 := sort.IntSlice{1, 2, 3, 4, 4, 3, 2, 1}
	ints2 := sort.IntSlice{1, 2, 3, 5, 4, 3, 2, 1}

	if !IsPalindrome(strings1) {
		t.Error("Failed for palindrome string slice")
	}

	if IsPalindrome(strings2) {
		t.Error("Failed for non-palindrome string slice")
	}

	if !IsPalindrome(ints1) {
		t.Error("Failed for palindrome int slice")
	}

	if IsPalindrome(ints2) {
		t.Error("Failed for non-palindrome int slice")
	}
}

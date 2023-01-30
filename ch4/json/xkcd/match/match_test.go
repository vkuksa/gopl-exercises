package match_test

import (
	"testing"

	. "gopl-exercises/ch4/json/xkcd/match"
)

func TestDirectMatch(t *testing.T) {
	input := "Black coat"
	test_against := "coat"

	if !DirectMatch(input, test_against) {
		t.Errorf("No DirectMatch for %s, %s", input, test_against)
	}

	test_against = "coal"
	if DirectMatch(input, test_against) {
		t.Errorf("No DirectMatch for %s, %s", input, test_against)
	}
}

func TestFuzzyMatch(t *testing.T) {
	input := "Black coat"
	test_against := "coat"

	if !FuzzyMatch(input, test_against, 1) {
		t.Errorf("No FuzzyMatch for %s, %s", input, test_against)
	}

	test_against = "coal"
	if FuzzyMatch(input, test_against, 1) {
		t.Errorf("No FuzzyMatch for %s, %s", input, test_against)
	}
	if !FuzzyMatch(input, test_against, 2) {
		t.Errorf("No FuzzyMatch for %s, %s", input, test_against)
	}

	test_against = "soal"
	if FuzzyMatch(input, test_against, 2) {
		t.Errorf("No FuzzyMatch for %s, %s", input, test_against)
	}
	if !FuzzyMatch(input, test_against, 3) {
		t.Errorf("No FuzzyMatch for %s, %s", input, test_against)
	}
}

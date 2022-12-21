package main

import "testing"

var test_case = []string{"A man, a plan, a canal: Panama", "123", "321"}

const n = 10000000

func BenchmarkEcho1(b *testing.B) {
	for i := 0; i < n; i++ {
		Echo1(test_case)
	}
}

func BenchmarkEcho2(b *testing.B) {
	for i := 0; i < n; i++ {
		Echo2(test_case)
	}
}

func BenchmarkEcho3(b *testing.B) {
	for i := 0; i < n; i++ {
		Echo3(test_case)
	}
}

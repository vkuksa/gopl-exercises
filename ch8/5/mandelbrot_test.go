package mandelbrot

import (
	"fmt"
	"testing"
)

func BenchmarkDrawSequential(b *testing.B) {
	for n := 0; n < b.N; n++ {
		DrawSequential()
	}
}

func BenchmarkDrawParallel(b *testing.B) {
	for i := 0; i <= 8; i++ {
		n := 1 << i
		b.Run(fmt.Sprintf("workers_%d", n), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				DrawParallel(n)
			}
		})
	}
}

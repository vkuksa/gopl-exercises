package pipeline

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestPipeline(t *testing.T) {
	in := make(chan int)

	expected := rand.Int()
	go func() {
		in <- expected
		close(in)
	}()

	if got := <-pipeline(in, 10); expected != got {
		t.Errorf("Invalid pipeline: expected %d, got %d", expected, got)
	}
}

func BenchmarkPipeline(b *testing.B) {
	for i := 15; i <= 25; i++ {
		stages := 1 << i
		b.Run(fmt.Sprintf("stage_%d", stages), func(b *testing.B) {
			in := make(chan int)
			go func() {
				for j := 0; j < i; j++ {
					in <- j
				}
				close(in)
			}()

			for range pipeline(in, stages) {
			}
		})
	}
}

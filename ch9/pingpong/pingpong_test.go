package pingpong

import (
	"fmt"
	"testing"
	"time"
)

const (
	n = 5
)

func BenchmarkPipeline(b *testing.B) {
	done := make(chan struct{})
	var comms int64

	pingpong(done, &comms)

	time.Sleep(n * time.Second)
	close(done)

	fmt.Printf("Result: %d/s\n", comms/n)
}

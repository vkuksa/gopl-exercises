// Exercise 9.5: Write a program with two goroutines that send messages back and forth over
// two unbuffered channels in ping-pong fashion. How many communications per second can
// the program sustain?

package pingpong

import (
	"sync/atomic"
)

func pingpong(done chan struct{}, comms *int64) {
	first, second := make(chan int), make(chan int)

	sender := func(in chan int, out chan int) {
		for {
			select {
			case <-done:
				return
			case v := <-in:
				atomic.AddInt64(comms, 1)
				out <- v
			}
		}
	}

	go sender(first, second)
	go sender(second, first)

	first <- 0
	return
}

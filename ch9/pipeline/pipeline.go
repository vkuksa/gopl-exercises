// Exercise 9.4: Construct a pipeline that connects an arbitrary number of goroutines with chan-
// nels. What is the maximum number of pipeline stages you can create without running out of
// memory? How long does a value take to transit the entire pipeline?

package pipeline

func pipeline(in chan int, stages int) (out chan int) {
	var localIn, localOut chan int

	localOut = in
	for i := 0; i < stages; i++ {
		localIn = localOut
		localOut = make(chan int)

		go func(in chan int, out chan int) {
			for value := range in {
				out <- value
			}
			close(out)
		}(localIn, localOut)

	}
	out = localOut

	return
}

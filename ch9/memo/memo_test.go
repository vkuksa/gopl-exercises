package memo_test

import (
	"testing"

	"gopl-exercises/ch9/memo"
	"gopl-exercises/ch9/memo/memotest"
)

var httpGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := memo.New(httpGetBody)
	memotest.Sequential(t, m)
	m.Close()
}

// NOTE: not concurrency-safe!  Test fails.
func TestConcurrent(t *testing.T) {
	m := memo.New(httpGetBody)
	memotest.Concurrent(t, m)
	m.Close()
}

func TestCancel(t *testing.T) {
	m := memo.New(httpGetBody)
	memotest.Cancelling(t, m)
	m.Close()
}

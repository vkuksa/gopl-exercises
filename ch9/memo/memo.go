// Exercise 9.3: Extend the Func type and the (*Memo).Get method so that callers may provide
// an optional done channel through which they can cancel the operation (§8.9). The results of a
// cancelled Func call should not be cached.

package memo

//!+Func

// Func is the type of the function to memoize.
type Func func(key string, done chan struct{}) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

//!-Func

//!+get

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	done     chan struct{}
	response chan<- result // the client wants a single result
}

type Memo struct{ requests chan request }

// New returns a memoization of f.  Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string, done chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, done, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

//!-get

//!+monitor

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key, req.done) // call f(key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string, done chan struct{}) {
	// Evaluate the function.
	e.res.value, e.res.err = f(key, done)
	// Broadcast the ready condition.
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}

//!-monitor

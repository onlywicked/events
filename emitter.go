package events

import (
	"sync"
)

// emitter is an in-memory EventEmitter
type emitter struct {
	sync.RWMutex
	once            sync.Once
	capacity        uint
	globalListeners []chan Data
	listeners       map[string][]chan Data
	closed          bool
}

func listener(ch <-chan Data, callback Listener) {
	for data := range ch {
		callback(data)
	}
}

func (e *emitter) OnAll(l Listener) {
	e.Lock()
	defer e.Unlock()

	ch := make(chan Data)
	e.globalListeners = append(e.globalListeners, ch)

	go listener(ch, l)
}

// On attaches a listener to a particular event
func (e *emitter) On(event string, l Listener) {
	e.Lock()
	defer e.Unlock()

	if e.closed {
		return
	}

	ch := make(chan Data, e.capacity)
	e.listeners[event] = append(e.listeners[event], ch)

	go listener(ch, l)
}

func (e *emitter) Emit(event string, data Data) {
	e.Lock()
	defer e.Unlock()

	if e.closed {
		return
	}

	// setting data event to emitted event for easier logging and debugging
	data.event = event

	var wg sync.WaitGroup

	// global listeners
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, ch := range e.globalListeners {
			ch <- data
		}
	}()
	// event listeners
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, ch := range e.listeners[event] {
			ch <- data
		}
	}()
	wg.Wait()
}

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

// NewEventEmitter returns an in-memory implementation event emitter
func NewEventEmitter(capacity ...uint) EventEmitter {
	var c uint = 1
	if len(capacity) > 0 {
		c = capacity[0]
	}

	e := new(emitter)
	e.capacity = c
	e.globalListeners = []chan Data{}
	e.listeners = make(map[string][]chan Data)
	e.closed = false

	return e
}

func (e *emitter) OnAll(l Listener) {
	e.Lock()
	defer e.Unlock()

	if e.closed {
		return
	}

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

// List returns all the events listeners are listening to
func (e *emitter) List() []string {
	var events []string

	for event := range e.listeners {
		events = append(events, event)
	}

	return events
}

// Close stops the event emitter from emitting events and
// clean up all the resources.
// Note: Once Close has been called it cannot be restarted again.
func (e *emitter) Close() {
	e.once.Do(func() {
		e.Lock()
		defer e.Unlock()

		for _, ch := range e.globalListeners {
			close(ch)
		}

		for _, ls := range e.listeners {
			for _, ch := range ls {
				close(ch)
			}
		}
		e.globalListeners = nil
		e.listeners = nil
		e.closed = true
	})
}

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

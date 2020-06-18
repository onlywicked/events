package events

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-ch:
		return false
	case <-time.After(timeout):
		return true
	}
}

func Test_emitter_OnAll(t *testing.T) {
	tests := []struct {
		name        string
		globalCount int
	}{
		{
			globalCount: 2,
		},
		{
			globalCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &emitter{}

			for i := 0; i < tt.globalCount; i++ {
				ee.OnAll(func(Data) {})
			}

			got := len(ee.globalListeners)
			if got != tt.globalCount {
				t.Errorf("emitter.OnAll(): global listeners count expected = %v, got = %v", tt.globalCount, got)
			}
		})
	}

}

func Test_emitter_On(t *testing.T) {
	tests := []struct {
		name           string
		listenersCount int
	}{
		{
			name:           "should attach 10 listeners",
			listenersCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &emitter{
				listeners: make(map[string][]chan Data),
			}

			for i := 0; i < tt.listenersCount; i++ {
				ee.On("event", func(Data) {})
			}

			got := len(ee.listeners["event"])
			if got != tt.listenersCount {
				t.Errorf("emitter.On(): listeners count expected = %v, got = %v", tt.listenersCount, got)
			}

		})
	}
}

func Test_emitter_Emit(t *testing.T) {
	tests := []struct {
		emitCount int
	}{
		{
			emitCount: 3,
		},
		{
			emitCount: 10,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("should emit %d events", tt.emitCount), func(t *testing.T) {
			ee := &emitter{
				listeners: make(map[string][]chan Data),
			}

			var wg sync.WaitGroup
			var lock sync.RWMutex

			listened := 0
			wg.Add(tt.emitCount)
			ee.On("event", func(Data) {
				lock.Lock()
				defer lock.Unlock()
				listened++
				wg.Done()
			})

			globalListened := 0
			wg.Add(tt.emitCount)
			ee.OnAll(func(Data) {
				lock.Lock()
				defer lock.Unlock()
				globalListened++
				wg.Done()
			})

			for i := 0; i < tt.emitCount; i++ {
				ee.Emit("event", Data{})
			}

			if waitTimeout(&wg, 2*time.Second) {
				t.Fatalf("emitter.Emit(): attached listener ran too long")
			}

			if listened != tt.emitCount {
				t.Errorf("emitter.Emit(): to emitted event count expected = %v, listened = %v", tt.emitCount, listened)
			}
			if globalListened != tt.emitCount {
				t.Errorf("emitter.Emit(): to emitted event count expected = %v, global listened = %v", tt.emitCount, globalListened)
			}

		})
	}
}

func Test_emitter_List(t *testing.T) {
	tests := []struct {
		name  string
		in    []string
		count int
	}{
		{
			name:  "should return all events",
			in:    []string{"event1", "event2", "event3"},
			count: 3,
		},
		{
			name:  "should return only distinct events",
			in:    []string{"event1", "event2", "event1", "event2"},
			count: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ee := &emitter{
				listeners: make(map[string][]chan Data),
				capacity:  1,
			}

			for _, event := range tt.in {
				ee.On(event, func(data Data) {})
			}

			gotCount := len(ee.List())

			if gotCount != tt.count {
				t.Errorf("emitter.List(): total event count expected = %v, got = %v", tt.count, gotCount)
			}
		})
	}
}

func Test_emitter_Close(t *testing.T) {
	tests := []struct {
		name           string
		globalCount    int
		listenersCount int
	}{
		{
			name:           "should clean up allocated resources",
			listenersCount: 10,
			globalCount:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &emitter{
				listeners:       make(map[string][]chan Data),
				globalListeners: []chan Data{},
			}

			for i := 0; i < tt.globalCount; i++ {
				ee.OnAll(func(Data) {})
			}

			for i := 0; i < tt.listenersCount; i++ {
				ee.On("event", func(Data) {})
			}

			ee.Close()

			if ee.closed != true {
				t.Errorf("emitter.Close(): unable to clean up resources")
			}

			if len(ee.globalListeners) != 0 {
				t.Errorf("emitter.Close(): unable to clean up global listeners")
			}

			if len(ee.listeners["event"]) != 0 {
				t.Errorf("emitter.Close(): unable to clean up event listeners")
			}

			for i := 0; i < tt.globalCount; i++ {
				ee.OnAll(func(Data) {})
			}

			for i := 0; i < tt.listenersCount; i++ {
				ee.On("event", func(Data) {})
			}

			if len(ee.globalListeners) != 0 {
				t.Errorf("emitter.Close(): attached global listeners after closing the emitter")
			}

			if len(ee.listeners["event"]) != 0 {
				t.Errorf("emitter.Close(): attached event listeners after closing the emitter")
			}
		})
	}
}

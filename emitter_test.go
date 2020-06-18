package events

import (
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

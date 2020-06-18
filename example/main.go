package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/onlywicked/events"
)

func main() {
	for i := 1; i <= 2; i++ {
		events.On("user:created", func(data events.Data) {
			log.Printf("logging from listener %d: %v\n", i, data)
		})
	}
	for i := 1; i <= 3; i++ {
		events.On("user:updated", func(data events.Data) {
			log.Printf("logging from listener %d: %v\n", i, data)
		})
	}

	for i := 1; i <= 3; i++ {
		// attaching global event listener to all events.
		events.OnAll(func(data events.Data) {
			// send the data to external services like rollbar or slack.
			// here we are just logging it to stderr
			log.Printf("logging from global %d: %v\n", i, data)
		})
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// simulating events from different user service
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			events.Emit("user:created", events.Data{
				Payload: fmt.Sprintf("payload %d", i),
				Message: fmt.Sprintf("message %d", i),
			})
		}
	}()

	wg.Add(1)
	// simulating events from different user service
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			events.Emit("user:updated", events.Data{
				Payload: fmt.Sprintf("payload %d", i),
				Message: fmt.Sprintf("message %d", i),
			})
		}
	}()

	wg.Wait()
}

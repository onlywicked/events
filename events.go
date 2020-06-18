// Package events provide an in-memory implementation of event emitter pattern
package events

func init() {
	defaultEmitter = NewEventEmitter()
}

var defaultEmitter EventEmitter

// Listener is function that listens on an event
type Listener func(data Data)

// EventEmitter provides an interface to work with events and act on them.
// It can help to decouple third party service from business logic.
//
// For example
//
// func (s *UserService) Create(user User) {
//  your business logic
//	eventemitter.Emit("user:created", events.Data{ Payload: user })
// }
//
// eventemitter.On("user:created", func (data event.Data) {
//  user, _ := data.Payload.(User)
//	emailService.SendWelcomeEmail(user.Name, user.Email)
// })
type EventEmitter interface {
	List() []string
	Emit(event string, data Data)
	On(event string, l Listener)
	Close()
	OnAll(l Listener)
}

// List returns all the events listeners are listening to
func List() []string {
	return defaultEmitter.List()
}

// Emit emits an event to all the listeners listening
func Emit(event string, data Data) {
	defaultEmitter.Emit(event, data)
}

// On attaches a listener to an event
func On(event string, l Listener) {
	defaultEmitter.On(event, l)
}

// OnAll attaches a listener to all events
// In simpler terms, attaches a global listener
func OnAll(l Listener) {
	defaultEmitter.OnAll(l)
}

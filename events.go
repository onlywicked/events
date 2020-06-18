// Package events provide an in-memory implementation of event emitter pattern
package events

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

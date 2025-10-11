package events

import (
	"sync"
	"time"
)

// EventInterface represents a generic event in the system.
type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

// EventHandlerInterface represents a handler that processes events.
type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

// EventDispatcherInterface represents the event dispatcher that manages event handlers.
type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}

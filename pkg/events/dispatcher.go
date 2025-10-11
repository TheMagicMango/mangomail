package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

// EventDispatcher manages event handlers and dispatches events to them.
type EventDispatcher struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandlerInterface
}

// NewEventDispatcher creates a new EventDispatcher instance.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

// Dispatch sends an event to all registered handlers for that event type.
// Handlers are executed concurrently using goroutines and wait groups.
func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	ed.mu.RLock()
	defer ed.mu.RUnlock()

	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

// Register adds a new handler for a specific event type.
// Returns ErrHandlerAlreadyRegistered if the handler is already registered for this event.
func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

// Has checks if a specific handler is registered for an event type.
func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	ed.mu.RLock()
	defer ed.mu.RUnlock()

	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

// Remove removes a specific handler from an event type.
func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

// Clear removes all handlers from all event types.
func (ed *EventDispatcher) Clear() {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	ed.handlers = make(map[string][]EventHandlerInterface)
}

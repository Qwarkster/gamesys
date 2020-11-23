package gamesys

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
)

var (
	// LastMove should be the last time we cycled
	LastMove time.Time

	// Dt will be our system update timing
	Dt float64
)

// Controller manages the key handlers, running handler methods as appropriate.
// The current state of the system will depend on what handlers are present. If
// there are any system handlers present, they will override the application
// handlers until there are not any system handlers present.
type Controller struct {
	// Handlers are the connected keystrokes and functions.
	Handlers []*Handler

	// SystemHandlers are system handlers that overrule the application handlers
	// when present. Example is a messagebox that blocks and waits to be
	// closed by the user. Effectively the concept of system modal.
	SystemHandlers []*Handler

	// Engine is the engine the controller is running on.
	Engine *Engine
}

// Handler is our structure that we will create and add to the controller
type Handler struct {
	// ID will be a string we can use to identify a handler when we need to
	// remove it.
	ID string

	// The keypress we are checking
	Button pixelgl.Button

	// Sensitive will indicate if we JustPress...usually for menus
	Sensitive bool

	// The action to perform
	Action func()
}

// Initialize will setup any structure elements that require not being nil.
func (c *Controller) Initialize() {
	// Empty handler arrays
	c.Handlers = make([]*Handler, 0)
	c.SystemHandlers = make([]*Handler, 0)

	// Set the starting time.
	c.Engine.LastMove = time.Now()
}

// AddApplicationHandler will add a pixelgl button handler to our list. The
// sensitive argument indicates if we use the JustPress method, firing the
// handler once per keypress as opposed to a handler that would react as long
// as the key is held down, as a motion handler might behave.
func (c *Controller) AddApplicationHandler(id string, button pixelgl.Button, sensitive bool, action func()) {
	c.Handlers = append(c.Handlers, &Handler{ID: id, Button: button, Sensitive: sensitive, Action: action})
}

// AddSystemHandler does the same as the application method, only on our
// system collection. TODO: Refactor appropriately.
func (c *Controller) AddSystemHandler(id string, button pixelgl.Button, sensitive bool, action func()) {
	c.SystemHandlers = append(c.SystemHandlers, &Handler{ID: id, Button: button, Sensitive: sensitive, Action: action})
}

// RemoveSystemHandler will remove a handler from the collection. If we provide
// an invalid id, it shouldn't remove anything.
func (c *Controller) RemoveSystemHandler(id string) {
	currenthandlers := len(c.SystemHandlers)
	newHandlers := make([]*Handler, 0)

	if currenthandlers > 0 {
		for _, h := range c.SystemHandlers {
			if h.ID != id {
				newHandlers = append(newHandlers, h)
			}
		}
	}

	c.SystemHandlers = newHandlers
}

// Run will loop through our controllers running any handlers that are setup.
func (c *Controller) Run() {
	// Manage timing
	c.Engine.Dt = time.Since(c.Engine.LastMove).Seconds()
	c.Engine.LastMove = time.Now()

	// If we have system handlers, we overrule application handlers.
	if len(c.SystemHandlers) > 0 {
		c.processHandlers(c.SystemHandlers)
	} else {
		c.processHandlers(c.Handlers)
	}
}

// This should likely not be used externally yet, if at all.
func (c *Controller) processHandlers(handlers []*Handler) {
	for _, h := range handlers {
		if h.Sensitive {
			if c.Engine.win.JustPressed(h.Button) {
				h.Action()
			}
		} else {
			if c.Engine.win.Pressed(h.Button) {
				h.Action()
			}
		}
	}
}

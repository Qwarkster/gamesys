package gamesys

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
)

// Controller manages the key handlers, running handler methods as appropriate.
// The current state of the system will depend on what handlers are present. If
// there are any system handlers present, they will override the application
// handlers until there are not any system handlers present.
type Controller struct {
	// Handlers are the collections of Handler Sets
	Handlers map[string][]*Handler

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
	// Setup handler map
	c.Handlers = make(map[string][]*Handler)

	// Set the starting time.
	c.Engine.LastMove = time.Now()
}

// AddHandler will add the indicated type of handler to this control. `sensitive`
// will indicate if the JustPress method is used, which triggers the handler
// once. Otherwise the handler will act as a game button, allowing it to be
// held for repeated effect.
func (c *Controller) AddHandler(class string, id string, button pixelgl.Button, sensitive bool, action func()) {
	handler, ok := c.Handlers[class]
	if !ok {
		handler = make([]*Handler, 0)
	}
	handler = append(handler, &Handler{ID: id, Button: button, Sensitive: sensitive, Action: action})
	c.Handlers[class] = handler
}

// RemoveHandler will remove a handler from the provided handler list.
func (c *Controller) RemoveHandler(class string, id string) {
	currenthandlers := c.Handlers[class]
	newHandlers := make([]*Handler, 0)

	if len(currenthandlers) > 0 {
		for _, h := range currenthandlers {
			if h.ID != id {
				newHandlers = append(newHandlers, h)
			}
		}
	}

	c.Handlers[class] = newHandlers
}

// Run will loop through our controllers running any handlers that are setup.
func (c *Controller) Run() {
	// Manage timing
	c.Engine.Dt = time.Since(c.Engine.LastMove).Seconds()
	c.Engine.LastMove = time.Now()

	// If we have system handlers, we overrule application handlers.
	if len(c.Handlers["system"]) > 0 {
		c.processHandlers(c.Handlers["system"])
	} else {
		c.processHandlers(c.Handlers["app"])
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

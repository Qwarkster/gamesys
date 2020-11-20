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

// Controller will handle keystrokes and perform the expected
// functions. These functions will be setup in the main app.
type Controller struct {
	// Handlers are the connected keystrokes and functions.
	Handlers []*Handler

	// Engine is the engine the controller is running on.
	Engine *Engine
}

// Handler is our structure that we will create and add to the controller
type Handler struct {
	// The keypress we are checking
	Button pixelgl.Button

	// Sensitive will indicate if we JustPress...usually for menus
	Sensitive bool

	// The action to perform
	Action func()
}

// Initialize will prepare our empty array of Handlers
func (c *Controller) Initialize() {
	// Set the starting time.
	c.Engine.LastMove = time.Now()
}

// Add will add a pixelgl button handler to our list
func (c *Controller) Add(button pixelgl.Button, sensitive bool, action func()) {
	c.Handlers = append(c.Handlers, &Handler{Button: button, Sensitive: sensitive, Action: action})
}

// Run will loop through our controllers running any actions
func (c *Controller) Run(win *pixelgl.Window) {
	// Manage timing
	c.Engine.Dt = time.Since(c.Engine.LastMove).Seconds()
	c.Engine.LastMove = time.Now()

	for _, h := range c.Handlers {
		if h.Sensitive {
			if win.JustPressed(h.Button) {
				h.Action()
			}
		} else {
			if win.Pressed(h.Button) {
				h.Action()
			}
		}
	}
}

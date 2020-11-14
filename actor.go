package gamesys

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
)

// Actor is an element that will move around within
// the view.
type Actor struct {
	// Position of the actor, needs to be relative to map
	Position pixel.Vec

	// Destinations will be preset by running scripts.
	Destinations []pixel.Vec

	// Src is the source graphic
	Src pixel.Picture

	// Output should be the sprite
	Output *pixel.Sprite

	// Clip is the area of the actor relative to the map.
	Clip pixel.Rect

	// Speed will set a speed modifier for the actor.
	Speed float64

	// Visible determines if we see the sprite or not.
	Visible bool

	// Collision determines if it collides with anything or not
	Collision bool
}

// NewActor creates a new actor and returns it
func NewActor(filename string, position pixel.Vec) Actor {
	newActor := Actor{Visible: false, Speed: Config.Float("DefaultActorSpeed"), Collision: true, Position: position}
	//newActor.Destinations = make([]pixel.Vec, 0)
	newActor.Src, err = LoadImage("characters/" + filename)

	if err != nil {
		fmt.Printf("Error loading image: %s", err.Error())
		os.Exit(2)
	}

	// Create our sprite.
	newActor.Output = pixel.NewSprite(newActor.Src, newActor.Src.Bounds())
	newActor.Clip = newActor.Output.Frame().Moved(newActor.Position)
	return newActor
}

// SetClip will create a clipping box based on the current actor position.
func (a *Actor) SetClip() {
	adjust := pixel.V(a.Output.Frame().W()/2, a.Output.Frame().H()/2)
	a.Clip = a.Output.Frame().Moved(a.Position.Sub(adjust))
}

// MoveTo will move the actor to an absolute map position.
func (a *Actor) MoveTo(position pixel.Vec) {
	a.Position = position
	a.SetClip()
}

// Move will move the actor according to the provided vector.
func (a *Actor) Move(distance pixel.Vec) {
	a.Position = a.Position.Add(distance)
	a.SetClip()
}

// Show will show the actor
func (a *Actor) Show() {
	a.Visible = true
}

// Hide will hide the actor
func (a *Actor) Hide() {
	a.Visible = false
}

// Toggle will toggle the visibility
func (a *Actor) Toggle() {
	if a.Visible {
		a.Hide()
	} else {
		a.Show()
	}
}

// Collides will set if the actor reacts to the collision layer.
func (a *Actor) Collides(collision bool) {
	a.Collision = collision
}

// Draw will draw the respective actor to the provided destination.
func (a *Actor) Draw(v *View) {
	drawMatrix := pixel.IM.Moved(a.Position.Sub(v.Camera.Min))
	a.Output.Draw(v.Rendered, drawMatrix)
}

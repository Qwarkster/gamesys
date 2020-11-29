package gamesys

import (
	"errors"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// View is the object that is rendered to the screen. It manages
// the associated graphics and actors, as well as any motion directly
// related to the view.
type View struct {
	// Visible indicates if the view should be rendered.
	Visible bool

	// Background is the background color of the view.
	Background color.RGBA

	// Src indicates the source picture data to draw from
	Src *pixel.PictureData

	// Focus is the actor that the view will follow. This actor will be
	// restricted to the bounds of the view.
	Focus *Actor

	// VisibleActors are the actors that are actually visible on this
	// view.
	VisibleActors []string

	// Output is the rendered map which is drawn to the screen.
	Output []*pixel.Sprite

	// Rendered is our background canvas to draw onto which will be
	// flipped to the screen.
	Rendered *pixelgl.Canvas

	// DesignView is the function that will be called to draw our view.
	// With consideration for views that don't focus on a map.
	DesignView func()

	// Position of our view on the window.
	Position pixel.Vec

	// Camera is the region of the map that is currently in view.
	// The camera also controls what is in display, period. So if we set it
	// once and don't change it, that's fine.
	Camera pixel.Rect

	// Speed is the speed modifier of this camera, not used yet.
	Speed float64

	// Scene will be the scene this view is a part of.
	Scene *Scene

	// Engine is passed through for ease of access
	Engine *Engine
}

// SetBackground sets the background color of the view.
func (v *View) SetBackground(bgcolor string) {
	v.Background = colornames.Map[bgcolor]
}

// UseMap will setup the view to use the scene map. It will error if there is
// no mapdata to use.
func (v *View) UseMap() error {
	// We should only be processing map data on a map
	// view, so this code needs a better home.
	if v.Scene.MapData != nil {
		v.Src = v.Scene.MapData.Img[0]
		v.Output = append(v.Output, pixel.NewSprite(v.Src, v.Src.Bounds()))
		return nil
	}

	// Return with error if we weren't able to use the map
	return errors.New("usemap: there is no mapdata available to use")
}

// Show the view.
func (v *View) Show() {
	v.Visible = true
}

// Hide the view.
func (v *View) Hide() {
	v.Visible = false
}

// Toggle view
func (v *View) Toggle() {
	if v.Visible {
		v.Hide()
	} else {
		v.Show()
	}
}

// Render will setup the viewable portion of the view.
func (v *View) Render() {
	if v.Visible {
		// Clear the existing view
		v.Rendered.Clear(v.Background)

		// Output means we have a map view to process.
		if v.Output != nil {
			for _, o := range v.Output {
				// Center on our actor if we have one.
				if v.Focus != nil {
					actor := v.Focus
					actorPos := actor.Position.Sub(v.Camera.Min)
					movement := actorPos.Sub(v.Rendered.Bounds().Center())
					v.CenterOn(movement)
				}

				// Grab the relevent section of map and place onto our view.
				o.Set(v.Src, v.Camera)
				o.Draw(v.Rendered, pixel.IM.Moved(v.Rendered.Bounds().Center()))
			}
		}

		// Now we work on the actors on the screen here.
		for _, a := range v.VisibleActors {
			v.Scene.Actors[a].Draw(v)
		}

		// See if this breaks first.
		if v.DesignView != nil {
			// This should draw to our debugger view.
			v.DesignView()
		}
	}
}

// FocusOn will focus on a specific actor
func (v *View) FocusOn(actor *Actor) {
	v.Focus = actor
}

// CenterOn will center the map on the actor who has the focus on that view.
func (v *View) CenterOn(movement pixel.Vec) {
	// Make a new temp camera based on where we would travel.
	newCamera := v.Camera.Moved(movement)

	if newCamera.Min.X < v.Src.Rect.Min.X {
		newCamera.Min.X = v.Src.Rect.Min.X
		newCamera.Max.X = newCamera.Min.X + v.Camera.W()
	} else if newCamera.Max.X > v.Src.Rect.Max.X {
		newCamera.Max.X = v.Src.Rect.Max.X
		newCamera.Min.X = newCamera.Max.X - v.Camera.W()
	}

	if newCamera.Min.Y < v.Src.Rect.Min.Y {
		newCamera.Min.Y = v.Src.Rect.Min.Y
		newCamera.Max.Y = newCamera.Min.Y + v.Camera.H()
	} else if newCamera.Max.Y > v.Src.Rect.Max.Y {
		newCamera.Max.Y = v.Src.Rect.Max.Y
		newCamera.Min.Y = newCamera.Max.Y - v.Camera.H()
	}

	// Assign the adjusted camera.
	v.Camera = newCamera
}

// Draw will draw the view to the scene, if it's visible.
func (v *View) Draw() {
	if v.Visible {
		// Ensure our view is rendered and up to date.
		v.Render()

		// This should draw onto the scene.
		v.Rendered.Draw(v.Scene.Rendered, pixel.IM.Moved(v.Position))
	}
}

// CameraContains will ensure that the camera contains the given
// rectangle. Should be refactored soon too.
func (v *View) CameraContains(target pixel.Rect) bool {
	return Contains(v.Rendered.Bounds(), target)
}

// Move will move view position within the window.
func (v *View) Move(position pixel.Vec) {
	v.Position = position
}

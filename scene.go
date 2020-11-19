package gamesys

import (
	"encoding/xml"
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Scene holds the information for a combination of views
// and actors. A view should be able to have multiple actors
// and multiple outputs.
type Scene struct {
	// XMLName is how we reference when loading xml information.
	XMLName xml.Name `xml:"scene"`

	// basespeed is the speed that this scene will run at.
	Basespeed float64 `xml:"basespeed,attr"`

	// Background is the background colour to clear this screen to.
	Background color.RGBA

	// Rendered is the canvas we draw to before flipping to screen.
	Rendered *pixelgl.Canvas

	// Views is the collection of views of the scene.
	Views map[string]*View

	// ViewOrder is the order the views are drawn down.
	ViewOrder []string

	// Actors is the collection of actors of the scene.
	Actors map[string]*Actor

	// MapData is the tiled data object.
	MapData *Map

	// Engine will be the reference to the engine we belong to.
	Engine *Engine
}

// Viewable will be useful at some point.
type Viewable interface {
	Show()
	Hide()
	Toggle()
}

// NewScene will create a new scene.
func (e *Engine) NewScene() *Scene {
	// Initialize our scene
	newScene := &Scene{Basespeed: e.Config.Default.Scene.Basespeed}

	// Setup the rest of our scene collections.
	newScene.Views = make(map[string]*View)
	newScene.Actors = make(map[string]*Actor)

	// Setup a drawing canvas based on screen size
	newRect := pixel.R(0, 0, e.Config.System.Window.Width, e.Config.System.Window.Height)
	newScene.Rendered = pixelgl.NewCanvas(newRect)

	// Pass it back
	return newScene
}

// NewMapScene will start a new scene with a mapfile. In this way we allow
// a scene to be run that has no map attached.
func (e *Engine) NewMapScene(file string) *Scene {
	// Start with a basic scene.
	newScene := e.NewScene()

	// Load our map file
	newMap := NewMap(file)

	// Attach mapdata to scene.
	newScene.MapData = &newMap

	// Get our actors from the mapfile.
	newScene.ActorsFromMapFile()

	// Return the completed scene.
	return newScene
}

// GetScene should grab a scene for easy reference.
func (e *Engine) GetScene(id string) *Scene {
	returnScene := e.Scenes[id]
	return returnScene
}

// GetActiveScene will get the currently active scene.
func (e *Engine) GetActiveScene() *Scene {
	return e.ActiveScene
}

// SetBackground will set the background color of the scene.
func (s *Scene) SetBackground(bgcolor string) {
	s.Background = colornames.Map[bgcolor]
}

// StartMapView sets a view up to use the scene map
// data.
func (s *Scene) StartMapView(view string) {
	// We should only be processing map data on a map
	// view, so this code needs a better home.
	if s.MapData != nil {
		s.Views[view].Src = s.MapData.Img[0]
		s.Views[view].Output = append(s.Views[view].Output,
			pixel.NewSprite(s.Views[view].Src,
				s.Views[view].Src.Bounds()))
	}
}

// ActorsFromMapFile will load up the actors that are setup
// on the current mapfile.
func (s *Scene) ActorsFromMapFile() {
	// We can load up the actors from the mapfile data at any time.

	for _, obj := range s.MapData.Src.ObjectGroups[0].Objects {
		if obj.Type == "Spawn" {
			// We need to grab our properties
			actorID := obj.Properties.GetString("gameID")
			collision := obj.Properties.GetBool("collide")
			file := obj.Properties.GetString("imgfile")

			// We need to reverse the Y position.
			newY := s.MapData.Size.Y - obj.Y
			startPos := pixel.V(obj.X, newY)

			// Create actor and populate fields.
			newActor := s.Engine.NewActor(file, startPos)
			newActor.Visible = obj.Visible
			newActor.Collision = collision

			// Add to main list.
			s.Engine.AddActor(actorID, newActor)

			// Attach this to the scene.
			s.AttachActor(actorID)
		}
	}
}

// AttachActor will attach an actor to the scene.
func (s *Scene) AttachActor(actor string) {
	newActor := s.Engine.Actors[actor]
	s.Actors[actor] = newActor
}

// MoveActor will move an actor within the scene.
func (s *Scene) MoveActor(actor *Actor, direction int) {
	// Calculate our base movement speed.
	speed := s.Basespeed * s.Engine.Dt

	// Adjust to account for speed of the actor we are moving.
	speed *= actor.Speed

	// Our movement values and flags.
	movement := pixel.ZV
	move := false

	// Create movement vector based on direction
	switch direction {
	case NORTH:
		movement.Y = speed
	case SOUTH:
		movement.Y -= speed
	case EAST:
		movement.X = speed
	case WEST:
		movement.X -= speed
	}

	// Find our new position.
	newPos := actor.Position.Add(movement)
	newClip := actor.Clip.Moved(movement)

	// We should not be doing the following stuff the way
	// we are doing it.

	// We need to find the containing view
	v := &View{}
	for _, view := range s.Views {
		if view.Focus == actor {
			v = view
		}
	}

	// If we have the focus, keep the actor on the screen.
	if v != nil && v.Focus == actor {
		// We can move only if we are within our limits.
		if v.CameraContains(newClip.Moved(pixel.ZV.Sub(v.Camera.Min))) {
			// We will check for collision here
			if actor.Collision {
				if s.CollisionFree(newClip) {
					// All is good to move safely.
					move = true
				}
			} else {
				// In screen but not checking for collisions
				move = true
			}
		} else if s.Contains(newClip) {
			move = true
		}

	} else {

		// We can still enable collisions
		if actor.Collision {
			if s.CollisionFree(newClip) {
				move = true
			}
		} else {
			move = true
		}

	}

	// Now to do what we gotta do.
	if move == true {
		actor.MoveTo(newPos)
	}
}

// CollisionFree will indicate the space is free of collisions.
func (s *Scene) CollisionFree(clip pixel.Rect) bool {
	// We can skip here so we don't have to worry about it in other places.
	for _, c := range s.MapData.Collision {
		if c.Intersects(clip) {
			return false
		}
	}

	// So if we skip collision, we are always safe.
	return true
}

// Contains will indicate if the rectangle is contained within this
// view's source map. Used for bounds checking against actors and
// the camera.
func (s *Scene) Contains(target pixel.Rect) bool {

	// False if any points are outside of the rect.
	for _, p := range target.Vertices() {
		if !s.MapData.Img[0].Rect.Contains(p) {
			return false
		}
	}
	return true
}

// ProcessActorDestinations will process any automated motions based on the
// destination lists of actors.
func (s *Scene) ProcessActorDestinations() {
	for _, a := range s.Actors {
		if a.Destinations != nil {
			// Here we need to figure out how far along to move.
			// We can move by a distance vector.
			dest := a.Destinations[0]
			motion := a.Position.To(dest)
			distance := math.Hypot(motion.X, motion.Y)
			travel := s.Engine.Dt * (s.Basespeed * a.Speed)

			// Do we travel all the way or not?
			if travel >= distance {
				// Here we reach the destination
				a.MoveTo(dest)
				if len(a.Destinations) > 1 {
					a.Destinations = a.Destinations[1:]
				} else {
					a.Destinations = nil
				}
			} else {
				// Here we calculate our finished motion position.
				diff := distance - travel
				ratio := diff / distance
				newDistance := pixel.V(ratio*motion.X, ratio*motion.Y)
				newPos := dest.Sub(newDistance)

				// Move to our hopefully new position
				a.MoveTo(newPos)
			}
		}
	}
}

// Render will draw our views onto our scene canvas.
func (s *Scene) Render() {
	for _, view := range s.ViewOrder {
		v := s.Views[view]
		// Render our view to canvas.
		v.Render()

		// Draw onto our scene
		v.Draw(s.Rendered)
		v.Rendered.Clear(v.Background)
	}
}

// Draw will draw the scene out to the provided destination.
func (s *Scene) Draw(win *pixelgl.Window) {
	// Adjust for the position on the screen
	win.SetMatrix(pixel.IM.Moved(s.Rendered.Bounds().Center()))

	// Now to put the canvas to the screen
	s.Rendered.Draw(win, pixel.IM)
}

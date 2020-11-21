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

	// Engine is the engine this scene belongs to.
	Engine *Engine
}

// NewView will create and return a new view. It's not tied to a map
// at this point, which is how it should be.
func (s *Scene) NewView(position pixel.Vec, camera pixel.Rect) *View {
	// A new view with some of our fields.
	newView := &View{Visible: false, Position: position, Camera: camera, Scene: s, Engine: s.Engine}

	// The canvas we prepare and flip to screen.
	newView.Rendered = pixelgl.NewCanvas(newView.Camera)

	return newView
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

// AttachView will add a view to the scene. We also need to add it to the view
// order so that it renders correctly.
func (s *Scene) AttachView(id string, view *View) {
	// Add to our view order on the scene.
	s.ViewOrder = append(s.ViewOrder, id)

	// Add to our scene here
	s.Views[id] = view
}

// RemoveView will destroy the view from the scene, also maintaining the vieworder.
func (s *Scene) RemoveView(id string) {
	// Loop through current view order, omitting the one matching id
	newViewOrder := make([]string, len(s.ViewOrder)-1)
	for _, v := range s.ViewOrder {
		if v != id {
			newViewOrder = append(newViewOrder, v)
		}
	}
	s.ViewOrder = newViewOrder

	// Remove our view
	delete(s.Views, id)
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

			// Add to our engine.
			s.Engine.AddActor(actorID, newActor)

			// Attach this to the scene.
			s.AttachActor(actorID)
		}
	}
}

// AttachActor will attach an actor to the scene.
func (s *Scene) AttachActor(actor string) {
	s.Actors[actor] = s.Engine.Actors[actor]
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

	if actor.Collision {
		if s.CollisionFree(newClip) {
			// All is good to move safely.
			move = true
		}
	} else if s.Contains(newClip) {
		move = true
	}

	// Now to do what we gotta do.
	if move == true {
		actor.MoveTo(newPos)
	}
}

// ProcessActorDestinations will move the relevent actors towards
// their respective destinations
func (s *Scene) ProcessActorDestinations() {
	for _, a := range s.Actors {
		if a.Destinations != nil {
			// Here we need to figure out how far along to move.
			// We can move by a distance vector.
			dest := a.Destinations[0]
			motion := a.Position.To(dest)
			distance := math.Hypot(motion.X, motion.Y)
			travel := Dt * (s.Basespeed * a.Speed)

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
	return Contains(s.MapData.Img[0].Rect, target)
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

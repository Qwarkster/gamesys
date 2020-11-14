package gamesys

import (
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// Our system constants we will use
const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

var (
	err error

	// Config holds our game configuration
	Config *Configuration

	// PixelWindow holds the configuration for the pixel window
	PixelWindow pixelgl.WindowConfig

	// win should be internal I think
	win *pixelgl.Window

	// ScriptActions holds the defined scripting actions for the game
	ScriptActions map[string]ScriptAction

	// Font is the basic text atlas for now
	Font *text.Atlas

	// Scenes holds a map of scenes
	Scenes map[string]Scene

	// ActiveScene will be the currently running scene.
	ActiveScene string

	// Actors will hold the actors for the game.
	Actors map[string]Actor

	// Control will hold a list of handlers
	Control *Controller

	// Logic will be the function that runs the game logic.
	Logic func()

	// LastMove should be the last time we cycled
	LastMove time.Time

	// Dt will be our system update timing
	Dt float64
)

// ConfigurePixel will build up the pixel configuration from our game
// configuration. In this way if any window options are changed, we simply
// use ConfigurePixel to update the pixel window.
func ConfigurePixel() {
	PixelWindow = pixelgl.WindowConfig{
		Title:  Config.String("title"),
		Bounds: pixel.R(0, 0, Config.Float("screenWidth"), Config.Float("screenHeight")),
		VSync:  true,
	}
}

// Initialize starts up the RPG engine
func Initialize() {
	// Set our pixel configuration
	ConfigurePixel()

	// Initialize window
	win, err = pixelgl.NewWindow(PixelWindow)
	if err != nil {
		panic(err)
	}

	// Setup the initial controller
	Control = &Controller{}
	Control.Initialize()

	// Setup our basic font
	Font = text.Atlas7x13

	// Initialize empty system maps.
	Scenes = make(map[string]Scene)
	Actors = make(map[string]Actor)
	ScriptActions = make(map[string]ScriptAction)

	// Now we can setup our core action library.
	CreateCoreActions()

}

// RunScriptAction will run the specified script action.
func RunScriptAction(action *Action) interface{} {
	if a, ok := ScriptActions[action.Action]; ok {
		return a.Runner(action.Args)
	}
	return nil
}

// RunScriptFile will load and run a script, presuming script directory
// and extension.
func RunScriptFile(file string) {
	if Config != nil {
		script := &Script{}
		script.Load(Config.String("ScriptFolder")+"/"+file+"."+Config.String("ScriptExtension"), false)
		_ = RunScript(script)
	}
}

// RunScript will run a game script, by default using our game script
// collection.
func RunScript(script *Script) interface{} {
	if script.Actions != nil {
		for _, a := range script.Actions {
			//ScriptActions
			_ = RunScriptAction(a)
		}
	}
	return nil
}

// AddScene will add a scene to the system.
func AddScene(id string, scene Scene) {
	Scenes[id] = scene
}

// ActivateScene will set the currently running scene.
func ActivateScene(scene string) {
	ActiveScene = scene
}

// AddActor will add an actor to the system.
func AddActor(id string, actor Actor) {
	Actors[id] = actor
}

// Run will run our main game processes
func Run() {

	for !win.Closed() {
		// Start main game loop, grab active scene.
		scene := Scenes[ActiveScene]

		// Clear to a color
		win.Clear(scene.Background)
		scene.Rendered.Clear(scene.Background)

		// Run our key handler
		Control.Run(win)

		// This will process custom game logic, technically optional.
		if Logic != nil {
			Logic()
		}

		// Process automatic movements via destinations.
		// Actors have destinations, the view is irrelevent.
		for _, a := range scene.Actors {
			// TODO: Convert this motion code into a method on actor.

			if a.Destinations != nil {
				// Here we need to figure out how far along to move.
				// We can move by a distance vector.
				dest := a.Destinations[0]
				motion := a.Position.To(dest)
				distance := math.Hypot(motion.X, motion.Y)
				travel := Dt * (scene.Speed * a.Speed)

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

		// We have multiple views, we should act accordingly
		for _, view := range scene.ViewOrder {
			v := scene.Views[view]
			// Render our view to canvas.
			v.Render()

			// Draw onto our scene
			v.Draw(scene.Rendered)
			v.Rendered.Clear(v.Background)
		}

		// Once we draw all our views, we should draw out to the window
		scene.Draw(win)

		win.Update()
	}
}

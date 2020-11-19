package gamesys

import (
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
)

// Engine is the core system that holds all running functionality.
type Engine struct {

	// Config holds our main configuration values to run the system.
	Config *Configuration

	// PixelWindow is our graphics window configuration.
	PixelWindow pixelgl.WindowConfig
	win         *pixelgl.Window

	// ScriptActions holds defined scripting actions.
	ScriptActions map[string]*ScriptAction

	// Font is our basic text atlas for system purposes.
	Font *text.Atlas

	// Scenes holds the various game scene contents.
	Scenes map[string]*Scene

	// ActiveScene is the currently running scene.
	ActiveScene *Scene

	// Actors holds the loaded actors for the game. They can be used across
	// scenes so it's not a good idea to tie them tightly to them. Scenes will
	// hold a list of actors that are visible or running on them.
	Actors map[string]*Actor

	// Control is the handlers that are loaded into the engine. Soon we should have global
	// control handlers, and scene independent handlers.
	Control *Controller

	// Logic is not used yet, but will be put in place to handle game logic functions that
	// should be run periodically, which are not tied to events.
	Logic func()

	// LastMove is the time of the last game cycle, used for managing game timing and motion.
	LastMove time.Time

	// Dt is used to calculate change in game cycle time, used for managing
	// game timing and motion.
	Dt float64
}

// ConfigurePixel will build up the pixel configuration from our game
// configuration. In this way if any window options are changed, we simply
// use ConfigurePixel to update the pixel window.
func (e *Engine) ConfigurePixel() {
	e.PixelWindow = pixelgl.WindowConfig{
		Title:  e.Config.System.Window.Title,
		Bounds: pixel.R(0, 0, e.Config.System.Window.Width, e.Config.System.Window.Height),
		VSync:  true,
	}
}

// Initialize starts up the RPG engine
func (e *Engine) Initialize() {
	// Set our pixel configuration
	e.ConfigurePixel()

	// Initialize window
	e.win, err = pixelgl.NewWindow(e.PixelWindow)
	if err != nil {
		panic(err)
	}

	// Setup the initial controller
	e.Control = &Controller{}
	e.Control.Initialize()

	// Setup our basic font
	e.Font = text.Atlas7x13

	// Initialize empty system maps.
	e.Scenes = make(map[string]*Scene)
	e.Actors = make(map[string]*Actor)
	e.ScriptActions = make(map[string]*ScriptAction)

	// Now we can setup our core action library.
	// TODO: This is too specific, should break it out of basic initialization.
	e.CreateCoreActions()
}

// RunScriptAction will run the specified script action.
func (e *Engine) RunScriptAction(action *Action) interface{} {
	if a, ok := e.ScriptActions[action.Action]; ok {
		return a.Runner(action.Args)
	}
	return nil
}

// RunScript will run a game script, by default using our game script
// collection.
func (e *Engine) RunScript(script *Script) interface{} {
	if script.Actions != nil {
		for _, a := range script.Actions {
			//ScriptActions
			_ = e.RunScriptAction(a)
		}
	}
	return nil
}

// RunScriptFile will load and run a script, presuming script directory
// and extension.
func (e *Engine) RunScriptFile(file string) {
	if e.Config != nil {
		script := &Script{}
		script.Load(e.Config.System.Scripting.Dir+"/"+file+"."+e.Config.System.Scripting.Extension, false)
		_ = e.RunScript(script)
	}
}

// AddScene will add a scene to the system.
func (e *Engine) AddScene(id string, scene *Scene) {
	e.Scenes[id] = scene
}

// ActivateScene will set the currently running scene.
func (e *Engine) ActivateScene(scene string) {
	e.ActiveScene = e.Scenes[scene]
}

// AddActor will add an actor to the system.
func (e *Engine) AddActor(id string, actor *Actor) {
	e.Actors[id] = actor
}

// Run will run our main game processes
func (e *Engine) Run() {

	for !e.win.Closed() {
		// Start main game loop, grab active scene.
		scene := e.ActiveScene

		// Clear to a color
		e.win.Clear(scene.Background)
		scene.Rendered.Clear(scene.Background)

		// Run our key handler
		e.Control.Run(e.win)

		// This will process custom game logic, technically optional.
		if e.Logic != nil {
			e.Logic()
		}

		// Process automatic movements via destinations.
		e.ActiveScene.ProcessActorDestinations()

		// We should render our scene appropriately.
		e.ActiveScene.Render()

		// Time to spit out the scene.
		scene.Draw(e.win)

		e.win.Update()
	}
}

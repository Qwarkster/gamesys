package gamesys

import (
	"os"
	"testing"

	"github.com/faiface/pixel/pixelgl"
	"github.com/stretchr/testify/assert"
)

var (
	// testEngine will be our test engine for this process.
	testEngine *Engine

	// mainLoop will test if we run our main game loop or not.
	mainLoop bool
)

func TestMain(m *testing.M) {
	// Wrap our testing system.
	pixelgl.Run(func() {
		testEngine = &Engine{}
		testEngine.Initialize("test_assets/config.xml")

		// test1.script contains 3 new scenes
		testEngine.RunScriptFile("test1")

		// Activate our first scene
		testEngine.ActivateScene("test1")

		// Decide if we are running main loop or not.
		mainLoop = true

		results := m.Run()
		os.Exit(results)
	})

}

func TestNewEngine(t *testing.T) {
	// We hate repetition
	e := testEngine

	// Check all systems are setup after initialization.
	assert.NotNil(t, e, "We should have an engine here that's not nil.")
	assert.NotNil(t, e.PixelWindow, "Our pixel configuration should be properly loaded.")
	assert.NotNil(t, e.win, "Our window should be created properly.")
	assert.NotNil(t, e.Control, "Our controller should be created properly.")
	assert.NotNil(t, e.Font, "Our system font should be created properly.")
	assert.NotNil(t, e.Scenes, "Our scenes map should be initialized properly.")
	assert.NotNil(t, e.Actors, "Our actors map should be initialized properly.")
	assert.NotNil(t, e.ScriptActions, "Our script actions should be created properly.")

	// We know we should have at least 1 core action for now.
	assert.Greater(t, len(e.ScriptActions), 0, "We should have at least 1 ScriptAction configured.")
}

func TestRunScriptAction(t *testing.T) {
	// See about a bad script action
	badscript := testEngine.RunScriptAction(&Action{Action: "blah", Args: make([]interface{}, 3)})

	assert.Nil(t, badscript)
}

func TestRunScriptFile(t *testing.T) {

	assert.Equal(t, 3, len(testEngine.Scenes), "We should have 3 scenes loaded.")
}

func TestActivateScene(t *testing.T) {

	assert.NotNil(t, testEngine.ActiveScene, "We should have an active scene")
}

func TestMessageBox(t *testing.T) {
	testEngine.DisplayMessageBox("Hello there!\nWe have multiple lines here.\nWhat shall we do with them?")

	scene := testEngine.ActiveScene

	assert.NotNil(t, scene.Views["messagebox"], "We should have a messagebox view.")
	assert.Equal(t, 1, len(testEngine.Control.Handlers["system"]), "We should have a system handler set.")
}

func TestRun(t *testing.T) {
	if !mainLoop {
		t.Skip("We don't always want to run the main loop.")
	}

	testEngine.Control.AddHandler("app", "right", pixelgl.KeyRight, false, func() {
		testEngine.ActiveScene.MoveActor(testEngine.Actors["monster"], 0)
	})
	testEngine.Control.AddHandler("app", "left", pixelgl.KeyLeft, false, func() {
		testEngine.ActiveScene.MoveActor(testEngine.Actors["monster"], 180)
	})
	testEngine.Control.AddHandler("app", "up", pixelgl.KeyUp, false, func() {
		testEngine.ActiveScene.MoveActor(testEngine.Actors["monster"], 90)
	})
	testEngine.Control.AddHandler("app", "down", pixelgl.KeyDown, false, func() {
		testEngine.ActiveScene.MoveActor(testEngine.Actors["monster"], 270)
	})
	testEngine.Run()
}

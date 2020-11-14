package rpg

import (
	"github.com/faiface/pixel"
)

// CreateCoreActions sets up the basic scripting actions that will
// always be included in the system.
func CreateCoreActions() {
	// ***********************************
	// NewScene will create a basic scene.
	// ===================================
	// NewScene scene_id background
	// -----------------------------------
	newScript := NewScriptAction("NewScene", func(args []interface{}) interface{} {
		// Setup arguments.
		id := args[0].(string)
		bgcolor := args[1].(string)

		// Create the new scene.
		newScene := NewScene()
		newScene.SetBackground(bgcolor)

		// Attach our scene.
		AddScene(id, newScene)

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// *****************************************************
	// NewMapScene will create a scene with a preloaded map.
	// =====================================================
	// NewMapScene scene_id file background
	// -----------------------------------------------------
	newScript = NewScriptAction("NewMapScene", func(args []interface{}) interface{} {
		// Setup arguments.
		id := args[0].(string)
		file := args[1].(string)
		bgcolor := args[2].(string)

		// Create the new scene.
		newScene := NewMapScene(file)
		newScene.SetBackground(bgcolor)

		// Attach our scene.
		AddScene(id, newScene)

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// ************************************************
	// NewView will load a map and attach a view to it.
	// ================================================
	// NewView scene_id view_id x y width height
	// ------------------------------------------------
	newScript = NewScriptAction("NewView", func(args []interface{}) interface{} {
		// Setup arguments.
		sceneID := args[0].(string)
		id := args[1].(string)
		x := StrFloat(args[2])
		y := StrFloat(args[3])
		width := StrFloat(args[4])
		height := StrFloat(args[5])
		bgcolor := args[6].(string)

		// Grab the scene we will be working with.
		scene := Scenes[sceneID]

		// Setup the position and camera rectangle.
		newPos := pixel.V(x, y)
		newCam := pixel.R(0, 0, width, height)

		// Create and add to our system.
		newView := NewView(id, newPos, newCam)

		// Set our background color on the view.
		newView.SetBackground(bgcolor)

		scene.Views[id] = &newView
		scene.ViewOrder = append(scene.ViewOrder, id)

		// Back into the collection
		Scenes[sceneID] = scene

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// *********************************************************
	// StartMapView will setup a view to use the scene map data.
	// =========================================================
	// StartMapView scene_id view_id
	// ---------------------------------------------------------
	newScript = NewScriptAction("StartMapView", func(args []interface{}) interface{} {
		// Setup arguments.
		scene := args[0].(string)
		view := args[1].(string)

		// This is an easy one we hope.
		Scenes[scene].StartMapView(view)

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// **************************************
	// ShowView will show the specified view.
	// ======================================
	// ShowView scene_id view_id
	// -------------------------
	newScript = NewScriptAction("ShowView", func(args []interface{}) interface{} {
		// Setup arguments.
		scene := args[0].(string)
		view := args[1].(string)

		// Show our view.
		Scenes[scene].Views[view].Show()

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// **********************************************************
	// NewActor will load a new actor. All visibility and collision options
	// should be set. We will not automatically focus on a view.
	// ==========================================================
	// NewActor scene_id actor_id imgfile x y visible collision
	// ----------------------------------------------------------------
	newScript = NewScriptAction("NewActor", func(args []interface{}) interface{} {
		// Setup arguments.
		sceneID := args[0].(string)
		id := args[1].(string)
		file := args[2].(string)
		x := StrFloat(args[3])
		y := StrFloat(args[4])
		visible := StrBool(args[5])
		collision := StrBool(args[6])

		// This area is copied right now, need to find the right home for it.
		// Create actor and populate fields.
		newActor := NewActor(file, pixel.Vec{X: x, Y: y})
		newActor.Visible = visible
		newActor.Collision = collision

		// Add to main list.
		AddActor(id, newActor)

		// Attach this to the scene.
		scene := Scenes[sceneID]
		scene.AttachActor(id)
		Scenes[sceneID] = scene

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// ****************************************************
	// ViewFocus will set which actor a view is focused on.
	// ====================================================
	// ViewFocus scene_id view_id actor_id
	// ----------------------------------------------------
	newScript = NewScriptAction("ViewFocus", func(args []interface{}) interface{} {
		// Setup arguments.
		scene := args[0].(string)
		viewID := args[1].(string)
		actor := args[2].(string)

		// Acting on a view so this is something for sanity.
		view := Scenes[scene].Views[viewID]

		view.FocusOn(Scenes[scene].Actors[actor])

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// ********************************************************************
	// ActorVisible will setup or add to the list on the view of the actors
	// that are rendered on that view.
	// ====================================================================
	// ActorVisible scene_id actor_id views_id...
	// ------------------------------------------
	newScript = NewScriptAction("ActorVisible", func(args []interface{}) interface{} {
		// Setup arguments.
		scene := args[0].(string)
		actor := args[1].(string)
		views := args[2:]

		// Attach our actor to the given views.
		for _, v := range views {
			view := v.(string)
			Scenes[scene].Views[view].VisibleActors = append(Scenes[scene].Views[view].VisibleActors, actor)
		}

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// *********************************************
	// ActorSpeed will set the actor speed modifier.
	// =============================================
	// ActorSpeed scene_id actor_id speed
	// ---------------------------------------------
	newScript = NewScriptAction("ActorSpeed", func(args []interface{}) interface{} {
		// Setup arguments.
		scene := args[0].(string)
		actor := args[1].(string)
		speed := StrFloat(args[2])

		// Set the speed on our actor
		Scenes[scene].Actors[actor].Speed = speed

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// ***********************************************************************
	// MoveActor will move an actor to the specified destination. The instant
	// flag will determine if we simply relocate, or animate the actor towards
	// that destination.
	// =======================================================================
	// MoveActor scene_id actor_id x y instant
	// -----------------------------------------------------------------------
	newScript = NewScriptAction("MoveActor", func(args []interface{}) interface{} {
		// Setup arguments.
		scene := args[0].(string)
		actorID := args[1].(string)
		x := StrFloat(args[2])
		y := StrFloat(args[3])
		instant := StrBool(args[4])

		actor := Scenes[scene].Actors[actorID]

		// If instant, then we just move it.
		if instant {
			actor.MoveTo(pixel.V(x, y))
		} else {
			// We give our actor a destination.
			actor.Destinations = append(actor.Destinations, pixel.V(x, y))
		}

		return nil
	})
	ScriptActions[newScript.Action] = newScript

	// **********************************************************************
	// MoveView will move the view accordingly. Instantaneous motion for now,
	// but eventually we could add some cool effects.
	// ======================================================================
	// MoveView scene_id view_id x y
	// ----------------------------------------------------------------------
	newScript = NewScriptAction("MoveView", func(args []interface{}) interface{} {
		// Setup arguments.
		scene := args[0].(string)
		view := args[1].(string)
		x := StrFloat(args[2])
		y := StrFloat(args[3])

		Scenes[scene].Views[view].Move(pixel.V(x, y))

		return nil
	})
	ScriptActions[newScript.Action] = newScript
}

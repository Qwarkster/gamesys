package gamesys

import (
	"bufio"
	"os"
	"strings"
)

// Script will hold a sequence or collection of commands. In theory this is
// is what a script file might get loaded into.
type Script struct {
	Actions []*Action
}

// Add will add the given event onto the end of our collection. It will
// return the index of the added command. This is useful when maintaining a
// collection of custom actions as opposed to a sequence.
func (s *Script) Add(action string, args ...interface{}) int {
	s.Actions = append(s.Actions, &Action{Action: action, Args: args})
	return len(s.Actions) - 1
}

// Load will open up the requested script file and parse the actions into
// the script. We can either overwrite or append.
func (s *Script) Load(file string, appendScript bool) error {
	// Open our file, return on error
	scriptfile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer scriptfile.Close()

	// If append is false, we need to clear the script.
	if !appendScript {
		s.Actions = []*Action{}
	}

	// Create our reader
	scriptReader := bufio.NewReader(scriptfile)

	// And now our scanner
	scriptScanner := bufio.NewScanner(scriptReader)

	scriptLine := ""
	// Loop while we are scanning
	for scriptScanner.Scan() {
		// Get our next line
		scriptLine = scriptScanner.Text()

		// Split our line up
		results := strings.Fields(scriptLine)

		// Our first item should be the command
		action := results[0]

		// Arguments should be the rest
		args := make([]interface{}, len(results)-1)
		for i := range results[1:] {
			args[i] = results[i+1]
		}

		// Add our action to our script
		s.Actions = append(s.Actions, &Action{Action: action, Args: args})
	}

	return nil
}

// ScriptAction will hold the actual content of the actions.
type ScriptAction struct {
	Action string
	Runner func([]interface{}) interface{}
}

// NewScriptAction will create and return a new ScriptAction.
func NewScriptAction(action string, runner func([]interface{}) interface{}) ScriptAction {
	newScriptAction := ScriptAction{Action: action, Runner: runner}

	return newScriptAction
}

package gamesys

// Action will hold the instructions on an action to perform
type Action struct {
	// The action keyword
	Action string

	// The arguments for this command
	Args []interface{}
}

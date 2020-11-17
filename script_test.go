package gamesys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScript(t *testing.T) {
	newScript := NewScript()
	assert.IsType(t, &Script{}, newScript, "We should be returning a proper script")
}

func TestScriptLoad(t *testing.T) {
	newScript := NewScript()

	// Bad Load
	err := newScript.Load("blahblah", false)
	assert.Error(t, err, "We should throw an error on bad file load")

	// Good Load
	err = newScript.Load("test_assets/scripts/testing.script", false)
	assert.NoError(t, err, "We should not have an error when loading a valid file")

	// Test script appending
	assert.Equal(t, 3, len(newScript.Actions), "We should start with 3 actions")
	err = newScript.Load("test_assets/scripts/testing.script", true)
	assert.Equal(t, 6, len(newScript.Actions), "We should now have 6 actions")

}

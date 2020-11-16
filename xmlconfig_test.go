package gamesys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	newconfig := LoadConfiguration("test_assets/config.xml")

	// It works, not empty
	assert.NotNil(t, newconfig, "Configuration should not be nil")

	// We expect some things
	assert.NotNil(t, newconfig.System, "Our system should have configuration information")

	// Check our window settings
	assert.Equal(t, float64(640), newconfig.System.Window.Width, "Width should be a float")
	assert.Equal(t, float64(480), newconfig.System.Window.Height, "Height should be a float")

	// Are our speeds failing?
	assert.Equal(t, 200.0, newconfig.Default.Scene.Basespeed, "Basespeed should not be 0")
	assert.Equal(t, 1.0, newconfig.Default.Actor.Speed, "Actor speed modifier should not be 0")
}

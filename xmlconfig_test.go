package gamesys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	// Bad config file?
	newconfig, newerr := LoadConfiguration("badfile")
	assert.Equal(t, newconfig, &Configuration{}, "Configuration should not be set")
	assert.Error(t, newerr, "Error should most definitely occur")

	newconfig, newerr = LoadConfiguration("test_assets/config.xml")

	// It works, not empty
	assert.NotNil(t, newconfig, "Configuration should not be nil")
	assert.NoError(t, newerr, "Error should not occur")

	// We expect some things
	assert.NotNil(t, newconfig.System, "Our system should have configuration information")

	// Check our window settings
	assert.Equal(t, float64(640), newconfig.System.Window.Width, "Width should be a float")
	assert.Equal(t, float64(480), newconfig.System.Window.Height, "Height should be a float")

	// Are our speeds failing?
	assert.Equal(t, 200.0, newconfig.Default.Scene.Basespeed, "Basespeed should not be 0")
	assert.Equal(t, 1.0, newconfig.Default.Actor.Speed, "Actor speed modifier should not be 0")

	// We need basic messagebox configuration
	msgConfig := newconfig.Default.MessageBox
	assert.NotNil(t, msgConfig.Color, "We should have a color set")
	assert.NotNil(t, msgConfig.BGColor, "We should have a background color set")
	assert.Equal(t, 320.0, msgConfig.X, "We should have an X position set")
	assert.Equal(t, 240.0, msgConfig.Y, "We should have a Y position set")
	assert.Equal(t, 100.0, msgConfig.Height, "We should have a height set")
	assert.Equal(t, 200.0, msgConfig.Width, "We should have a width set")
}

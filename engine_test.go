package gamesys

import (
	"testing"

	"github.com/faiface/pixel/pixelgl"
	"github.com/stretchr/testify/assert"
)

func TestNewEngine(t *testing.T) {
	// I think we need to start the pixelgl engine first, almost first
	testEngine := &Engine{}
	testEngine.Initialize("test_assets/config.xml")
	pixelgl.Run(testEngine.Run)

	assert.NotNil(t, testEngine, "Engine should create just fine")
}

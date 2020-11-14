package gamesys

import (
	"testing"

	"github.com/faiface/pixel"
	"github.com/stretchr/testify/assert"
)

func TestLoadImage(t *testing.T) {
	imgtype := &pixel.PictureData{}

	// Valid should be painless.
	newimg, err := LoadImage("test_assets/characters/demo.png")
	assert.IsType(t, imgtype, newimg)
	assert.Nil(t, err)

	// Invalid should be just as painless
	newimg, err = LoadImage("blahblah")
	assert.IsType(t, nil, newimg)
	assert.NotNil(t, err)

	// How about a corrupt file
	newimg, err = LoadImage("test_assets/characters/haha.png")
	assert.IsType(t, nil, newimg)
	assert.NotNil(t, err)

}

func TestContains(t *testing.T) {
	target := pixel.R(0, 0, 200, 200)
	inside := pixel.R(20, 20, 40, 40)
	outside := pixel.R(300, 200, 320, 220)
	edge := pixel.R(190, 100, 210, 120)

	assert.True(t, Contains(target, inside), "Target inside should be true")
	assert.False(t, Contains(target, outside), "Target outside should be false")
	assert.False(t, Contains(target, edge), "Target on edge is should be false")
}

func TestStrFloat(t *testing.T) {
	// Expected input
	f := StrFloat("5.0")

	// Unexpected input
	bad := StrFloat("what")

	// Numeric input
	what := StrFloat(123)

	assert.Equal(t, 5.0, f)
	assert.Equal(t, float64(0), bad)
	assert.Equal(t, float64(0), what)
}

func TestStrBool(t *testing.T) {
	// True input
	good := StrBool("true")

	// False input
	falsy := StrBool("false")

	// Random string input
	stringy := StrBool("whatstring")

	// Bad input
	weird := StrBool(293)

	assert.True(t, good, "true should return true")
	assert.False(t, falsy, "false should return false")
	assert.False(t, weird, "invalid input should return false")
	assert.False(t, stringy, "invalid input string should return false")
}

package gamesys

import (
	"image"
	"os"
	"strconv"

	"github.com/faiface/pixel"
)

// utility.go will hold functions that don't quite fit within the rpg
// concept, but are still useful. These functions could potentially be
// broken out into their own subpackage down the road.

// FailError will panic and fail on error. Use this on errors we can't
// really handle gracefully.
func FailError(e error) {
	if e != nil {
		panic(e)
	}
}

// LoadImage will give us some picture data.
func LoadImage(path string) (pixel.Picture, error) {
	file, err := os.Open(path)

	// Throw an error on file opening, likely file not found.
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	// Throw an error during file processing, possibly corrupt file.
	if err != nil {
		return nil, err
	}

	// Return our error in case we can handle it elsewhere.
	return pixel.PictureDataFromImage(img), err
}

// Contains will indicate if the rectangle is contained within another
// rectangle.
func Contains(container pixel.Rect, target pixel.Rect) bool {

	// False if any points are outside of the rect.
	for _, p := range target.Vertices() {
		if !container.Contains(p) {
			return false
		}
	}
	return true
}

// Some generic stuff we could break out later.

// StrFloat will return a string as a float64
func StrFloat(s interface{}) float64 {
	// We need to have a string.
	value, ok := s.(string)
	if ok {
		process, _ := strconv.ParseFloat(value, 64)
		return process
	}
	return float64(0)
}

// StrBool will return a string as a bool
func StrBool(s interface{}) bool {
	// If it's not true, it's false. With this logic, if it's not valid, it's false.
	value, ok := s.(string)
	if ok {
		process, _ := strconv.ParseBool(value)
		return process
	}
	return false

}

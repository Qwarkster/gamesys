package gamesys

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

// Map will contain our map with some easy to use stuff, like for rendering
type Map struct {
	// Basic map data as loaded from file
	Src *tiled.Map

	// Size will be the size of our map, pulled from our map data
	Size pixel.Vec

	// The rendered full map
	Img []*pixel.PictureData

	// Our collision information, a collection of map objects.
	Collision []*pixel.Rect
}

// NewMap will load and initialize a map from a mapfile. If we need
// to change maps during a game, it makes sense to reset everything about
// a map so that we don't have any lingering artifacts.
func NewMap(mapfile string) Map {
	newMap := Map{}

	newMap.Img = make([]*pixel.PictureData, 0)

	// Load up the source map file.
	newMap.Src, err = tiled.LoadFromFile(mapfile)

	// We will have to handle errors more gracefully.
	if err != nil {
		fmt.Printf("Error loading mapfile: %s", err.Error())
		os.Exit(2)
	}

	// Grab some of our map information
	newMap.Size = pixel.V(float64(newMap.Src.Width*newMap.Src.TileWidth), float64(newMap.Src.Height*newMap.Src.TileHeight))

	// This creates our map renderer.
	renderer, err := render.NewRenderer(newMap.Src)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}

	// Render all visible layers.
	err = renderer.RenderVisibleLayers()
	if err != nil {
		fmt.Printf("layer unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}

	// Convert into pixel's image/sprite format.
	newMap.Img = append(newMap.Img, pixel.PictureDataFromImage(renderer.Result))

	// Setup collision objects, the first object layer will hold that, eventually configurable
	for _, obj := range newMap.Src.ObjectGroups[0].Objects {
		if obj.Type == "Collision" {
			// We need to build the rects properly. The tiled
			// map starts from the top down, have to reverse y.
			newY := newMap.Size.Y - obj.Y - obj.Height
			newCollision := pixel.R(obj.X+2, newY+2, obj.X+obj.Width-2, newY+obj.Height-2)

			newMap.Collision = append(newMap.Collision, &newCollision)
		}
	}

	// Ready to go
	return newMap
}

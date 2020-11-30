package gamesys

import (
	"encoding/xml"
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

// MessageBox sets options for the system messagebox.
type MessageBox struct {
	XMLName xml.Name `xml:"messagebox"`
	Color   string   `xml:"color,attr"`
	BGColor string   `xml:"bgcolor,attr"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	Height  float64  `xml:"height,attr"`
	Width   float64  `xml:"width,attr"`
}

// DisplayMessageBox will display a message on screen and then wait for user
// input.
func (e *Engine) DisplayMessageBox(msg string) {
	// Grab our configuration options for simplicity.
	msgConfig := e.Config.Default.MessageBox

	// We have to work on the current scene and edit it's view order to
	// put the messagebox down last, and remove it when done.
	scene := e.ActiveScene

	// Create our new messagebox view.
	x := msgConfig.X
	y := msgConfig.Y
	height := msgConfig.Height
	width := msgConfig.Width

	scene.NewView("messagebox", pixel.V(x, y), pixel.R(0, 0, width, height), msgConfig.BGColor)

	// Grab our new message view
	msgView, err := scene.GetView("messagebox")
	if err != nil {
		// I do not expect to see this, unless something is broken elsewhere.
		panic(err)
	}

	// The messagebox should be visible
	msgView.Show()

	// Create a drawing method.
	msgView.DesignView = func() {
		// Get our configuration
		color := colornames.Map[msgConfig.Color]
		bgcolor := colornames.Map[msgConfig.BGColor]

		// Prepare colors
		msgView.Rendered.Clear(bgcolor)
		msgTxt := text.New(pixel.ZV, e.Font)
		msgTxt.Color = color

		// Calculate text size offset
		offsetTxt := text.New(pixel.ZV, e.Font)
		fmt.Fprint(offsetTxt, "Ay")
		offset := offsetTxt.Bounds().H()

		// TODO: We need to create word wrap within view.
		fmt.Fprintf(msgTxt, msg)

		// Render to our view, do this better soon.
		msgView.Rendered.SetMatrix(pixel.IM.Moved(pixel.V(2, height-offset+2)))
		msgTxt.Draw(msgView.Rendered, pixel.IM)
		msgView.Rendered.SetMatrix(pixel.IM)

	}

	// Creating our messagebox handler tells the system we have a messagebox
	// to process and wait for.
	e.Control.AddHandler("system", "messagebox", pixelgl.KeyEnter, true, func() {
		e.ActiveScene.RemoveView("messagebox")
		e.Control.RemoveHandler("system", "messagebox")
	})

}

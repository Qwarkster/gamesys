package gamesys

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// Configuration setting collection
type Configuration struct {
	XMLName xml.Name `xml:"configuration"`
	System  System   `xml:"system"`
	Default Default  `xml:"default"`
}

// System configuration setting.
type System struct {
	XMLName   xml.Name  `xml:"system"`
	Window    Window    `xml:"window"`
	Scripting Scripting `xml:"scripting"`
}

// Window is the options for starting pixel window.
type Window struct {
	XMLName xml.Name `xml:"window"`
	Width   float64  `xml:"width,attr"`
	Height  float64  `xml:"height,attr"`
	Title   string   `xml:"title,attr"`
}

// Scripting sets customizable script options.
type Scripting struct {
	XMLName   xml.Name `xml:"scripting"`
	Dir       string   `xml:"dir,attr"`
	Extension string   `xml:"extension,attr"`
}

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

// Default object values when not provided.
type Default struct {
	XMLName    xml.Name   `xml:"default"`
	Scene      Scene      `xml:"scene"`
	Actor      Actor      `xml:"actor"`
	MessageBox MessageBox `xml:"messagebox"`
}

// LoadConfiguration setups a configuration from an xml file.
func LoadConfiguration(file string) (*Configuration, error) {
	// New empty configuration
	newconfig := &Configuration{}

	// Open our xmlFile
	xmlFile, err := os.Open(file)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return newconfig, err
	}

	fmt.Println("Successfully Opened " + file)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Do the unmarshal thing

	xml.Unmarshal(byteValue, &newconfig)

	return newconfig, err
}

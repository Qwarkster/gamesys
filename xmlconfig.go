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

// Default object values when not provided.
type Default struct {
	XMLName xml.Name `xml:"default"`
	Scene   Scene    `xml:"scene"`
	Actor   Actor    `xml:"actor"`
}

// LoadConfiguration setups a configuration from an xml file.
func LoadConfiguration(file string) *Configuration {
	// Open our xmlFile
	xmlFile, err := os.Open(file)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened " + file)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Do the unmarshal thing
	newconfig := &Configuration{}
	xml.Unmarshal(byteValue, &newconfig)

	return newconfig
}
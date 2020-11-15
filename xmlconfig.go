/*package gamesys

// Configuration setting collection
type Configuration struct {
	System  System  `xml:"system"`
	Default Default `xml:"default"`
}

// System configuration setting
type System struct {
	Window Window `xml:"window"`
}

type Default struct {
	Scene	Scene `xml:scene`
	Actor	Actor `xml:actor`
}

/*<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <!--Basic system requirements-->
    <system>
        <window width="640" height="480" title="RPG Demo" />
        <scripting dir="scripts" extension="script" />
    </system>
    <!--Default structure values-->
    <default>
        <scene>
            <basespeed>200</basespeed>
        </scene>
        <actor>
            <speed>1</speed>
        </actor>
    </default>
</configuration>*/

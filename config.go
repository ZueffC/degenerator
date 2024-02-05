package main

import (
	"os"

	"gopkg.in/ini.v1"
)

var config = parseINI()

func parseINI() *ini.File {
	os.OpenFile("config.ini", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	config, err := ini.Load("config.ini")

	if err != nil {
		Error.Panicln("Can't parse config file: " + err.Error())
	}

	return config
}

func getPostsPath() string {
	if !config.HasSection("directory") {
		panic("Can't find directory section in config file")
	}

	section, _ := config.GetSection("directory")

	if !section.HasKey("path") {
		panic("Can't find path value in directory section")
	}

	iniPath, _ := section.GetKey("path")
	path := iniPath.String()

	return path
}

func getBlogName() string {
	if !config.HasSection("settings") {
		panic("Can't find settings section in config file")
	}

	section, _ := config.GetSection("settings")

	if !section.HasKey("name") {
		panic("Can't find path value in directory section")
	}

	iniBlogName, _ := section.GetKey("name")
	blogName := iniBlogName.String()

	return blogName
}

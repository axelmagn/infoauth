package main

import (
	"flag"
	"github.com/axelmagn/infoauth/infoauth"
)

var (
	configFile = flag.String("config", "config/settings.ecfg", "Config File")
)

const usageMsg = `
InfoAuth web server
`

func Init() {
	flag.Parse()

	// Set up Logging
	// TODO

	// Read config file
	err := infoauth.AddSettingsFromFile(*configFile)
	if err != nil { panic(err.Error()) }

	// set up models and db
	err = infoauth.InitModels()
	if err != nil { panic(err.Error()) }
}

func Serve() {}

func Close() {}

func main() {
	Init()
	Serve()
	Close()
}

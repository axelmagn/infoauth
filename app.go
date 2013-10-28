package main

import (
	"flag"
	"github.com/axelmagn/infoauth"
	"github.com/axelmagn/envcfg"
)

var (
	configFile	= flag.String("config", "config/settings.ecfg", "Config File")
)

const usageMsg = `
InfoAuth web server
`

func Init() {
	flag.Parse()

	// Set up Logging
	// TODO

	// Read config file
	settings, err := envcfg.ReadSettings(*configFile)
	infoauth.AddSettings(settings)
}

func Serve() {}

func Close() {}

func main() {
	Init()
	Serve()
	Close()
}
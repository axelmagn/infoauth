package main

import (
	"flag"
	"github.com/axelmagn/infoauth/infoauth"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	if err != nil {
		panic(err.Error())
	}

	// set up models and db
	err = infoauth.InitModels()
	if err != nil {
		panic(err.Error())
	}

	// set up a dummy user
	u := &infoauth.User{ID: 1}
	u.PlusProfile = "Google: Axel Magnuson"
	u.LinkedInProfile = "LinkedIn: Axel Magnuson"
	u.Save()
	log.Printf("Created Dummy user\n")

}

func Serve() {
	port := infoauth.GetSetting(infoauth.S_PORT)
	http.HandleFunc("/user/", infoauth.UserHandler)
	log.Printf("Starting Server on port %s...", port)
	http.ListenAndServe(":"+port, nil)
}

func Close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {

		log.Printf(sig.String())
		log.Printf("Recieved Interrupt Signal.")
		log.Printf("Closing...")

		log.Printf("Flushing Database...")
		infoauth.GetStore().Flush()

		os.Exit(0)
	}
}

func main() {
	Init()
	go Close()
	Serve()
}

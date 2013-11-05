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

	// set up config
	infoauth.InitOauthConfig()

	// set up models and db
	err = infoauth.InitModels()
	if err != nil {
		panic(err.Error())
	}
}

func Serve() {
	port := infoauth.GetSetting(infoauth.S_PORT)
	http.HandleFunc("/google/url/", infoauth.GetGoogleAuthURLHandler)
	http.HandleFunc("/linkedin/url/", infoauth.GetLinkedInAuthURLHandler)
	http.HandleFunc("/", infoauth.ExchangeCodeHandler)
	http.HandleFunc("/google/", infoauth.GoogleAuthRedirectHandler)
	http.HandleFunc("/linkedin/", infoauth.LinkedInAuthRedirectHandler)
	log.Printf("Starting Server on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
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

func CloseListen() {

}

func main() {
	Init()
	go Close()
	Serve()
}

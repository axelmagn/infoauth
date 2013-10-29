package main

import (
	"flag"
	"net/http"
	"log"
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

	// set up a dummy user
	u, err := infoauth.NewUser()
	if err != nil { panic(err.Error()) }
	u.PlusProfile = "Google: Axel Magnuson"
	u.LinkedInProfile = "LinkedIn: Axel Magnuson"
	u.Save()
	log.Printf("Created Dummy user\n")


}

func Serve() {
	port := infoauth.GetSetting(infoauth.S_PORT)
	http.HandleFunc("/user/", infoauth.UserHandler)
	http.ListenAndServe(":" + port, nil)
}

func Close() {

}

func main() {
	Init()
	Serve()
	Close()
}

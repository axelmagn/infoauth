package infoauth

import (
	"github.com/axelmagn/envcfg"
	"os"
)

var appSettings map[string]string = make(map[string]string)

var S_DB_PATH = "DB_PATH"
var S_DEBUG = "DEBUG"
var S_PORT = "PORT"

var S_LINKEDIN_CLIENT_ID = "LINKEDIN_CLIENT_ID"
var S_LINKEDIN_CLIENT_SECRET = "LINKEDIN_CLIENT_SECRET"
var S_LINKEDIN_SCOPE = "LINKEDIN_SCOPE"
var S_LINKEDIN_AUTH_URL = "LINKEDIN_AUTH_URL"
var S_LINKEDIN_TOKEN_URL = "LINKEDIN_TOKEN_URL"
var S_LINKEDIN_REDIRECT_URL = "LINKEDIN_REDIRECT_URL"

var S_GOOGLE_CLIENT_ID = "GOOGLE_CLIENT_ID"
var S_GOOGLE_CLIENT_SECRET = "GOOGLE_CLIENT_SECRET"
var S_GOOGLE_SCOPE = "GOOGLE_SCOPE"
var S_GOOGLE_AUTH_URL = "GOOGLE_AUTH_URL"
var S_GOOGLE_TOKEN_URL = "GOOGLE_TOKEN_URL"
var S_GOOGLE_REDIRECT_URL = "GOOGLE_REDIRECT_URL"
var S_GOOGLE_USERINFO_URL = "GOOGLE_USERINFO_URL"
var S_GOOGLE_PERSON_URL = "GOOGLE_PERSON_URL"

func GetSettings() map[string]string {
	return appSettings
}

func GetSetting(key string) string {
	return appSettings[key]
}

func AddSettingsFromFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	settings, err := envcfg.ReadSettings(f)
	if err != nil {
		return err
	}

	AddSettings(settings)
	return nil
}

func AddSettings(add map[string]string) {
	for k, v := range add {
		AddSetting(k, v)
	}
}
func AddSetting(k, v string) {
	appSettings[k] = v
}

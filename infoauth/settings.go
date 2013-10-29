package infoauth

import (
	"github.com/axelmagn/envcfg"
)

var appSettings map[string]string = make(map[string]string)

var S_DB_PATH = "DB_PATH"

func GetSettings() map[string]string {
	return appSettings
}

func GetSetting(key string) string {
	return appSettings[key]
}

func AddSettingsFromFile(fname string) error { 
	settings, err := envcfg.ReadSettings(fname)
	if err != nil {
		return err
	}
	AddSettings(settings)
	return nil
}

func AddSettings(add map[string]string) {
	for k, v := range add {
		appSettings[k] = v
	}
}

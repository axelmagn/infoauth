package infoauth

import (
	"github.com/axelmagn/envcfg"
	"os"
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

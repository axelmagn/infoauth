package infoauth

import (
	"github.com/axelmagn/envcfg"
)

var settings map[string]string = make(map[string]string)

func GetSettings() map[string]string {
	return settings
}

func AddSettings(add map[string]string) {
	for k, v := range add {
		settings[k] = v
	}
}

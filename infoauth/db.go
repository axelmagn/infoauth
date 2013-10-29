package infoauth

import (
	"github.com/steveyen/gkvlite"
	"os"
)

var defaultStore *gkvlite.Store

// init a new store from a settings path
func InitStore() (*gkvlite.Store, error) {
	db_path := GetSetting(S_DB_PATH)

	file, err := os.OpenFile(db_path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	s, err := gkvlite.NewStore(file)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func GetStore() *gkvlite.Store {
	return defaultStore
}

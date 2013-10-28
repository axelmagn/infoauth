package infoauth

import (
	"github.com/steveyen/gkvlite"
)


type Store interface {
	
}

func NewStore() (*Store, error) {
	return gkvlite.NewStore(f)
}

type Collection interface {

}

type Saveable interface {
	Collection() (Collection, error)
	Key()		([]byte, error)
	Value() 	([]byte, error)
}
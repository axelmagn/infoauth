package infoauth

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"errors"
	"github.com/steveyen/gkvlite"
)

var userCollectionKey = "users"
var userCollection *gkvlite.Collection

// set up model data
func InitModels() error {
	s, err := InitStore()
	if err != nil {
		return err
	}
	if s == nil {
		return errors.New("Store creation failed with no error.")
	}

	InitUserCollection()
	if userCollection == nil {
		return errors.New("Failed to initialize user collection")
	}

	return nil
}

// User Model
// implements Saveable
type User struct {
	ID              uint
	GoogleToken     oauth.Token
	LinkedInToken   oauth.Token
	PlusProfile     []byte
	LinkedInProfile []byte
}

// get the user's database key
func (u *User) Key() ([]byte, error) {
	return UintToHex(u.ID)
}

// get the user's database encoding (a byte array of json)
func (u *User) Value() ([]byte, error) {
	return json.Marshal(u)
}

// save the user to db
func (u *User) Save() error {
	k, err := u.Key()
	if err != nil {
		return err
	}
	v, err := u.Value()
	if err != nil {
		return err
	}

	return UserCollection().Set(k, v)
}

// retrieve a user from db by id
func GetUser(id uint) (*User, error) {
	out := &User{}

	idHex, err := UintToHex(id)
	if err != nil {
		return nil, err
	}
	raw, err := UserCollection().Get(idHex)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(raw, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// create a new user using pretty naive key assignment
func NewUser() (*User, error) {
	var userId uint

	user, err := UserCollection().MaxItem(false)
	if err != nil {
		return nil, err
	}

	if user != nil {
		userId, err = HexToUint(user.Key)
		userId++
		if err != nil {
			return nil, err
		}
	} else {
		userId = 0
	}

	out := &User{ID: userId}
	out.Save()
	return out, nil
}

// initialize the user collection
func InitUserCollection() *gkvlite.Collection {
	userCollection = GetStore().SetCollection(userCollectionKey, nil)
	return userCollection
}

// get the user collection
func UserCollection() *gkvlite.Collection {
	return userCollection
}

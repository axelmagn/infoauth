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
	_, err := InitStore()
	if err != nil {
		return err
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
	id              uint
	googleToken     oauth.Token
	linkedInToken   oauth.Token
	plusProfile     json.RawMessage
	linkedInProfile json.RawMessage
}

// get the user's database key
func (u *User) Key() ([]byte, error) {
	return UintToHex(u.id)
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
	user, err := UserCollection().MaxItem(false)
	if err != nil {
		return nil, err
	}

	userId, err := HexToUint(user.Key)
	userId++
	if err != nil {
		return nil, err
	}

	out := &User{id: userId}
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

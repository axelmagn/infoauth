package infoauth

import (
	"errors"
    "encoding/json"
    "code.google.com/p/goauth2/oauth"
    "github.com/steveyen/gkvlite"
)

var userCollectionKey = "users"
var userCollection *gkvlite.Collection

// set up model data
func InitModels() error {
	_, err := InitStore()
	if err != nil { return err }

	InitUserCollection()
	if userCollection == nil {
		return errors.New("Failed to initialize user collection")	
	}

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
	UserCollection().Set(u.Key(), u.Value())
}

// retrieve a user from db by id
func GetUser(id uint) (*User, error) {
	out := &User{}

	raw, err := UserCollection().Get(UintToHex(id))
	if err != nil { return _, err }

	json.Unmarshal(raw, out)
	if err != nil { return _, err }

	return out
}

// create a new user using pretty naive key assignment
func NewUser() (*User, error) {
	userId := HexToUint(UserCollection().MaxItem(false).Key) + 1
	out := &User{id: userId}
	out.Save()
	return out
}

// initialize the user collection
func InitUserCollection() *gkvlite.Collection {
	userCollection = GetStore().SetCollection(userCollectionKey, nil)
	return userCollection
}

// get the user collection
func UserCollection() (*gkvlite.Collection, error) {
	return userCollection
}

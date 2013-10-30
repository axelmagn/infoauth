package infoauth

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"errors"
	"github.com/steveyen/gkvlite"
)

var userCollectionKey = "users"
var userCollection *gkvlite.Collection

var googleUserIndexKey = "users.google_id"
var googleUserIndex *gkvlite.Collection

var linkedInUserIndexKey = "users.linkedin_id"
var linkedInUserIndex *gkvlite.Collection

var ErrorStoreCreationFailedSilently = errors.New("Store creation failed with no error.")
var ErrorInitUserCollectionFailed = errors.New("Failed to initialize user collection")

// initialize the user collection
func InitUserCollection() *gkvlite.Collection {
	userCollection = GetStore().SetCollection(userCollectionKey, nil)
	return userCollection
}

// initialize index of users by their google id
func InitGoogleUserIndex() *gkvlite.Collection {
	googleUserIndex = GetStore().SetCollection(googleUserIndexKey, nil)
	return googleUserIndex
}

// initialize index of users by their linkedin id
func InitLinkedInUserIndex() *gkvlite.Collection {
	linkedInUserIndex = GetStore().SetCollection(linkedInUserIndexKey, nil)
	return linkedInUserIndex
}

// get the user collection
func UserCollection() *gkvlite.Collection {
	return userCollection
}

func GoogleUserIndex() *gkvlite.Collection {
	return googleUserIndex
}

func LinkedInUserIndex() *gkvlite.Collection {
	return linkedInUserIndex
}

// User Model
// implements Saveable
type User struct {
	ID              uint
	GoogleToken     oauth.Token
	LinkedInToken   oauth.Token
	PlusProfile     string
	LinkedInProfile string
}

// create a new user using pretty naive key assignment
func NewUser() (*User, error) {
	userId, err := NewUserID()
	if err != nil {
		return nil, err
	}

	out := &User{ID: userId}
	err = out.Save()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewUserID() (uint, error) {
	var userId uint

	user, err := UserCollection().MaxItem(false)
	if err != nil {
		return 0, err
	}

	if user != nil {
		userId, err = HexToUint(user.Key)
		userId++
		if err != nil {
			return 0, err
		}
	} else {
		userId = 1
	}

	return userId, nil

}

// extract user object from raw json
func DecodeUser(raw []byte) (*User, error) {
	out := &User{}
	err := json.Unmarshal(raw, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// retrieve a user from db by id
// returns nil, nil if it can't find the user
// TODO: replace this corner case with proper error typing
func GetUser(id uint) (*User, error) {
	idHex, err := UintToHex(id)
	if err != nil {
		return nil, err
	}
	raw, err := UserCollection().Get(idHex)
	if err != nil {
		return nil, err
	} else if raw == nil {
		// we cheat a little here.
		// TODO: replace with properly typed errors
		return nil, nil
	}

	out, err := DecodeUser(raw)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// TODO
func GetUserByGoogleID(id string) (*User, error) {
	return nil, nil
}

// TODO
func GetUserByLinkedInID(id string) (*User, error) {
	return nil, nil
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

	// extract plus id for index
	// extract linkedin id for index

	return UserCollection().Set(k, v)
}

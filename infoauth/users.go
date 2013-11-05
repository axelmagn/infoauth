// This is a slightly broken library.  It was originally intended to
// store user information retrieved at the same time as an oauth
// exchange, but that requirement was later relaxed.  Users contains
// methods that were originally interspersed between infoauth package
// files. They have not been adapted to work as a standalone package.
package infoauth

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"errors"
	"github.com/steveyen/gkvlite"
	"log"
	"net/http"
	"strconv"
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue(UserIDKey)
	if idStr == "" {
		http.Error(w, "No User ID Specified.", http.StatusBadRequest)
		return
	}

	id64, err := strconv.ParseUint(idStr, 10, 0)
	id := uint(id64)
	if err != nil {
		http.Error(w, "Error parsing User id.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	u, err := GetUser(id)
	if err != nil {
		http.Error(w, "Error retrieving User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	if u == nil {
		http.Error(w, "User does not exist.", http.StatusBadRequest)
		return
	}

	raw, err := u.Value()
	if err != nil {
		http.Error(w, "Error encoding User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	w.Write(raw)
}

func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	raw := r.FormValue(UserContentKey)
	if raw == "" {
		http.Error(w, "No user submitted", http.StatusBadRequest)
		return
	}

	u, err := DecodeUser([]byte(raw))
	if err != nil {
		http.Error(w, "Could not parse User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	if u.ID == 0 {
		u.ID, err = NewUserID()
		if err != nil {
			http.Error(w, "Could not assign User ID.\n"+Debug(err), http.StatusInternalServerError)
			return
		}
	}

	err = u.Save()
	if err != nil {
		http.Error(w, "Could not save User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	v, err := u.Value()
	if err != nil {
		http.Error(w, "Could not decode user after saving.\n"+Debug(err), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

// takes a regex specifying path groups
// assumes that 2nd group is the user id
// returns a handler function for users
func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Printf("GET:\t%s", r.URL.Path)
		GetUserHandler(w, r)
	case "POST":
		log.Printf("POST:\t%s", r.URL.Path)
		SaveUserHandler(w, r)
	case "PUT":
		log.Printf("PUT:\t%s", r.URL.Path)
		SaveUserHandler(w, r)
	}
}

func GetGoogleUserInfo(token *oauth.Token) (*http.Response, error) {
	transport := &oauth.Transport{
		Token:     token,
		Config:    GoogleOauthConfig,
		Transport: http.DefaultTransport,
	}

	client := transport.Client()

	url := GetSetting(S_GOOGLE_USERINFO_URL)

	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func GetGooglePlusProfile(token *oauth.Token) (*http.Response, error) {
	transport := &oauth.Transport{
		Token:     token,
		Config:    GoogleOauthConfig,
		Transport: http.DefaultTransport,
	}

	client := transport.Client()

	url := GetSetting(S_GOOGLE_PERSON_URL)

	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return r, nil
}

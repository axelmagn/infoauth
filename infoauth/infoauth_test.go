package infoauth_test

import (
	"testing"
	"bytes"
	"github.com/axelmagn/infoauth/infoauth"
	"code.google.com/p/goauth2/oauth"
)

// test in rough order of init dependency

// settings

func TestAddSettings(t *testing.T) {
	settings := map[string]string {
		"DB_PATH": "/tmp/infoauth.gkvlite"
	}
	infoauth.AddSettings(settings)
}

func ExampleGetSettings() {
	fmt.Println(infoauth.GetSetting("DB_PATH"))
	// Output: /tmp/infoauth.gkvlite
}

// models

func TestInitModels(t *testing.T) {
	err := infoauth.InitModels()
	if err != nil { t.Error(err.Error()) }
}

func TestUser() {
	u, err := infoauth.NewUser()
	if err != nil { t.Error("NewUser error: " + err.Error) }
	if u == nil { t.Error("NewUser returns nil user.") }
	if u.id == nil { t.Error("NewUser returns user with nil id.") }

	u.googleToken = oauth.Token{"abc", "def", time.Time{}, nil}
	u.linkedInToken = oauth.Token{"hij", "klm", time.Time{}, nil}
	u.plusProfile = []byte("userPlusProfile")
	u.linkedInProfile = []byte("userLinkedInProfile")

	err = u.Save()
	if err != nil { t.Error("u.Save error: " + err.Error)}

	u2, err := GetUser(u.id)
	if err != nil {t.Error("GetUser error: " + err.Error)}

	if !bytes.Equal(u.Value(), u2.Value()) {
		t.Error("User value changed after storage.")
	}

}

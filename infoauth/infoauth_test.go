package infoauth_test

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"github.com/axelmagn/infoauth/infoauth"
	"math/rand"
	"testing"
	"time"
)

// test in rough order of init dependency

// settings

func TestAddSettings(t *testing.T) {
	settings := map[string]string{
		"DB_PATH": "/tmp/infoauth.gkvlite",
	}
	infoauth.AddSettings(settings)
}

// models

func TestInitModels(t *testing.T) {
	randString, err := infoauth.UintToHex(uint(rand.Uint32()))
	if err != nil {
		t.Error("String generation: " + err.Error())
	}
	dbPath := infoauth.GetSetting("DB_PATH") + string(randString)
	infoauth.AddSetting("DB_PATH", dbPath)
	err = infoauth.InitModels()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestUser(t *testing.T) {
	u, err := infoauth.NewUser()
	if err != nil {
		t.Error("NewUser error: " + err.Error())
	}
	if u == nil {
		t.Error("NewUser returns nil user.")
	}

	u.GoogleToken = oauth.Token{"abc", "def", time.Time{}, nil}
	u.LinkedInToken = oauth.Token{"hij", "klm", time.Time{}, nil}
	u.PlusProfile = "plusProfile"
	u.LinkedInProfile = "plusProfile"

	err = u.Save()
	if err != nil {
		t.Error("u.Save error: " + err.Error())
	}

	u2, err := infoauth.GetUser(u.ID)
	if err != nil {
		t.Error("GetUser error: " + err.Error())
	}

	uv, err := u.Value()
	if err != nil {
		t.Error("u.Value error: " + err.Error())
	}
	u2v, err := u2.Value()
	if err != nil {
		t.Error("u2.Value error: " + err.Error())
	}
	if !bytes.Equal(uv, u2v) {
		t.Error("User value changed after storage.")
	}

}

func TestEmptyUser(t *testing.T) {
	u, err := infoauth.NewUser()
	if err != nil {
		t.Error("NewUser error: " + err.Error())
	}
	if u == nil {
		t.Error("NewUser returns nil user.")
	}

	err = u.Save()
	if err != nil {
		t.Error("u.Save error: " + err.Error())
	}

	u2, err := infoauth.GetUser(u.ID)
	if err != nil {
		t.Error("GetUser error: " + err.Error())
	}

	uv, err := u.Value()
	if err != nil {
		t.Error("u.Value error: " + err.Error())
	}
	u2v, err := u2.Value()
	if err != nil {
		t.Error("u2.Value error: " + err.Error())
	}
	if !bytes.Equal(uv, u2v) {
		t.Error("User value changed after storage.")
	}

}

func TestGetUser_UserMissing(t *testing.T) {
	id := uint(rand.Uint32())
	u, err := infoauth.GetUser(id)
	if u != nil || err != nil {
		t.Error("Expected nil user and nil error for GetUser on missing id")
	}
}

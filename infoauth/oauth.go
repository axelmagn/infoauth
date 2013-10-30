//TODO: linkedin oauth handshakes
package infoauth

import (
	"time"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"code.google.com/p/goauth2/oauth"
	"github.com/steveyen/gkvlite"
)

var handshakeCollectionKey = "handshakes"
var handshakeCollection *gkvlite.Collection

var GoogleOauthConfig *oauth.Config
var GoogleOauthTransport *oauth.Transport
var LinkedInOauthConfig *oauth.Config

const StateLen = hex.DecodedLen(16)
const HandshakeExpireDuration = 5 * time.Minute

const (
	C_GOOGLE uint = iota,
	C_LINKEDIN
)


func InitOauthConfig() {
	GoogleOauthConfig = &oauth.Config{
		ClientId:		GetSetting(S_GOOGLE_CLIENT_ID),
		ClientSecret:	GetSetting(S_GOOGLE_CLIENT_SECRET),
		RedirectURL:	GetSetting(S_GOOGLE_REDIRECT_URL),
		Scope:			GetSetting(S_GOOGLE_SCOPE),
		AuthURL:		GetSetting(S_GOOGLE_AUTH_URL),
		TokenURL:		GetSetting(S_GOOGLE_TOKEN_URL),
	}
	GoogleOauthTransport = &oauth.Transport{Config: GoogleOauthConfig}

}

func InitHandshakeCollection() *gkvlite.Collection {
	handshakeCollection = GetStore().SetCollection(handshakeCollectionKey, nil)
	return handshakeCollection
}

func HandshakeCollection() *gkvlite.Collection {
	return handshakeCollection
}

type Handshake struct {
	State string
	Expires time.Time
	Config uint // we don't store a config pointer so that marshalling doesn't duplicate config

}

func NewGoogleAuthURL() (string, error) {
	// init & store new handshake struct
	randBytes := make([]bytes, StateLen)
	_, err := rand.Reader.Read(randBytes)
	if err != nil {return "", nil }
	randStr := hex.EncodeToString(randBytes)
	handshake := &Handshake{
		State: randStr,
		Expires: time.Now().Add(HandshakeExpireDuration),
		Config: C_GOOGLE,
	}

	// use handshake state token to get new url

}

func (h *Handshake) Key() ([]byte, error) {
	return []byte(h.State), nil
}

func (h *Handshake) Value() ([]byte, error) {
	return json.Marshal(h)
}

func (h *Handshake) Save() error {
	k := h.Key()

	v, err := h.Value()
	if err != nil { return err }

	return HandshakeCollection.Set(k, v)
}
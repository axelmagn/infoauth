//TODO: linkedin oauth handshakes
package infoauth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
	"github.com/steveyen/gkvlite"
	"code.google.com/p/goauth2/oauth"
)

var handshakeCollectionKey = "handshakes"
var handshakeCollection *gkvlite.Collection

var GoogleOauthConfig *oauth.Config
var GoogleOauthTransport *oauth.Transport
var LinkedInOauthConfig *oauth.Config
var LinkedInOauthTransport *oauth.Transport

var StateLen int

const StateByteLen = 16
const HandshakeExpireDuration = 5 * time.Minute

const (
	C_GOOGLE uint = iota
	C_LINKEDIN
)

func InitOauthConfig() {
	StateLen = hex.DecodedLen(StateByteLen)

	GoogleOauthConfig = &oauth.Config{
		ClientId:     GetSetting(S_GOOGLE_CLIENT_ID),
		ClientSecret: GetSetting(S_GOOGLE_CLIENT_SECRET),
		RedirectURL:  GetSetting(S_GOOGLE_REDIRECT_URL),
		Scope:        GetSetting(S_GOOGLE_SCOPE),
		AuthURL:      GetSetting(S_GOOGLE_AUTH_URL),
		TokenURL:     GetSetting(S_GOOGLE_TOKEN_URL),
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
	State   string
	Expires time.Time
	Config  uint // we don't store a config pointer so that marshalling doesn't duplicate config
	Exchanged bool
}

func DecodeHandshake(raw []byte) (*Handshake, error) {
	out := &Handshake{}
	err := json.Unmarshal(raw, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func NewGoogleAuthURL() (string, error) {
	// init & store new handshake struct
	randBytes := make([]byte, StateLen)
	_, err := rand.Reader.Read(randBytes)
	if err != nil {
		return "", nil
	}
	randStr := hex.EncodeToString(randBytes)
	h := &Handshake{
		State:   randStr,
		Expires: time.Now().Add(HandshakeExpireDuration),
		Config:  C_GOOGLE,
		Exchanged: false,
	}

	// use handshake state token to get new url
	err = h.Save()
	if err != nil {
		return "", err
	}
	return h.State, nil
}

func ExchangeCode(code, state string) (*oauth.Token, error) {
	// get handshake collection
	c := HandshakeCollection()
	if c == nil {
		return nil, errors.New("Could not get Handshake Collection")
	}

	// retrieve handshake by state and make sure it exists
	hraw, err := c.Get([]byte(state))
	if err != nil {
		return nil, err
	}
	if hraw == nil {
		return nil, errors.New("State token not found")
	}

	// decode serialized handshake
	h, err := DecodeHandshake(hraw)
	if err != nil {
		return nil, err
	}

	// get the correct trasport
	var transport *oauth.Transport
	switch h.Config {
	case C_GOOGLE:
		transport = GoogleOauthTransport
	case C_LINKEDIN:
		transport = LinkedInOauthTransport
	default:
		return nil, errors.New("Unknown Oauth configuration")
	}

	// exchange code for token
	token, err := transport.Exchange(code)
	if err != nil {
		return nil, err
	}

	// mark handshake as exchanged
	h.Exchanged = true
	err = h.Save()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (h *Handshake) Key() ([]byte, error) {
	return []byte(h.State), nil
}

func (h *Handshake) Value() ([]byte, error) {
	return json.Marshal(h)
}

func (h *Handshake) Save() error {
	k, err := h.Key()
	if err != nil {
		return err
	}

	v, err := h.Value()
	if err != nil {
		return err
	}

	return HandshakeCollection().Set(k, v)
}

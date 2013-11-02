//TODO: linkedin oauth handshakes
package infoauth

import (
	"code.google.com/p/goauth2/oauth"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/steveyen/gkvlite"
	"strings"
	"time"
	// "code.google.com/p/google-api-go-client/plus/v1"
)

var handshakeCollectionKey = "handshakes"
var handshakeCollection *gkvlite.Collection

var GoogleOauthConfig *oauth.Config
var GoogleOauthTransport *oauth.Transport
var LinkedInOauthConfig *oauth.Config
var LinkedInOauthTransport *oauth.Transport

const GoogleServiceName = "GOOGLE"
const LinkedInServiceName = "LINKEDIN"

const StateLen = 2 // size of uint16

const HandshakeExpireDuration = 5 * time.Minute

const (
	C_GOOGLE uint = iota
	C_LINKEDIN
	C_UNKNOWN
)

func InitOauthConfig() {

	scopeReplacer := strings.NewReplacer(",", " ")
	googleScope := scopeReplacer.Replace(GetSetting(S_GOOGLE_SCOPE))

	GoogleOauthConfig = &oauth.Config{
		ClientId:     GetSetting(S_GOOGLE_CLIENT_ID),
		ClientSecret: GetSetting(S_GOOGLE_CLIENT_SECRET),
		RedirectURL:  GetSetting(S_GOOGLE_REDIRECT_URL),
		Scope:        googleScope,
		AuthURL:      GetSetting(S_GOOGLE_AUTH_URL),
		TokenURL:     GetSetting(S_GOOGLE_TOKEN_URL),
	}
	GoogleOauthTransport = &oauth.Transport{Config: GoogleOauthConfig}

	linkedInScope := scopeReplacer.Replace(GetSetting(S_LINKEDIN_SCOPE))
	LinkedInOauthConfig = &oauth.Config{
		ClientId:     GetSetting(S_LINKEDIN_CLIENT_ID),
		ClientSecret: GetSetting(S_LINKEDIN_CLIENT_SECRET),
		RedirectURL:  GetSetting(S_LINKEDIN_REDIRECT_URL),
		Scope:        linkedInScope,
		AuthURL:      GetSetting(S_LINKEDIN_AUTH_URL),
		TokenURL:     GetSetting(S_LINKEDIN_TOKEN_URL),
	}
	LinkedInOauthTransport = &oauth.Transport{Config: LinkedInOauthConfig}

}

func InitHandshakeCollection() *gkvlite.Collection {
	handshakeCollection = GetStore().SetCollection(handshakeCollectionKey, nil)
	return handshakeCollection
}

func HandshakeCollection() *gkvlite.Collection {
	return handshakeCollection
}

func ConfigToConst(c *oauth.Config) uint {
	switch c {
	case GoogleOauthConfig:
		return C_GOOGLE
	case LinkedInOauthConfig:
		return C_LINKEDIN
	default:
		return C_UNKNOWN
	}
}

func ConstToConfig(c uint) *oauth.Config {
	switch c {
	case C_GOOGLE:
		return GoogleOauthConfig
	case C_LINKEDIN:
		return LinkedInOauthConfig
	default:
		return nil
	}
}

type Handshake struct {
	State         []byte
	Expires       time.Time
	Config        uint // we don't store a config pointer so that marshalling doesn't duplicate config
	Exchanged     bool
	OuterRedirect string
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
	u, _, err := NewAuthUrl(GoogleOauthConfig)
	return u, err
}

func NewLinkedInAuthURL() (string, error) {
	u, _, err := NewAuthUrl(LinkedInOauthConfig)
	return u, err
}

func NewGoogleAuthURLWithRedirect(redirectURL string) (string, error) {
	u, h, err := NewAuthUrl(GoogleOauthConfig)
	if err != nil {
		return u, err
	}

	h.OuterRedirect = redirectURL

	err = h.Save()
	if err != nil {
		return "", err
	}

	return u, err
}

func NewLinkedInAuthURLWithRedirect(redirectURL string) (string, error) {
	u, h, err := NewAuthUrl(LinkedInOauthConfig)
	if err != nil {
		return u, err
	}

	h.OuterRedirect = redirectURL

	err = h.Save()
	if err != nil {
		return "", err
	}

	return u, err
}

func NewAuthUrl(c *oauth.Config) (string, *Handshake, error) {
	// init & store new handshake struct
	randBytes := make([]byte, StateLen)
	_, err := rand.Reader.Read(randBytes)
	if err != nil {
		return "", nil, nil
	}

	configConst := ConfigToConst(c)
	if configConst == C_UNKNOWN {
		return "", nil, fmt.Errorf("Unknown config %v", c)
	}
	h := &Handshake{
		State:     randBytes,
		Expires:   time.Now().Add(HandshakeExpireDuration),
		Config:    configConst,
		Exchanged: false,
	}

	// use handshake state token to get new url
	err = h.Save()
	if err != nil {
		return "", nil, err
	}

	// convert to hex for printing
	stateHex := hex.EncodeToString(h.State)

	// get url using state
	url := c.AuthCodeURL(stateHex)

	// encode to string
	return url, h, nil
}

func ExchangeCode(code, stateHex string) (*oauth.Token, *Handshake, error) {
	// get handshake collection
	c := HandshakeCollection()
	if c == nil {
		return nil, nil, errors.New("Could not get Handshake Collection")
	}

	// retrieve handshake by stateHex and make sure it exists
	state, err := hex.DecodeString(stateHex)
	if err != nil {
		return nil, nil, err
	}
	hraw, err := c.Get(state)
	if err != nil {
		return nil, nil, err
	}
	if hraw == nil {
		return nil, nil, errors.New("State token not found")
	}

	// decode serialized handshake
	h, err := DecodeHandshake(hraw)
	if err != nil {
		return nil, nil, err
	}

	// TODO check that state isn't expired, and that it hasn't already been redeemed

	// get the correct trasport
	var transport *oauth.Transport
	switch h.Config {
	case C_GOOGLE:
		transport = GoogleOauthTransport
	case C_LINKEDIN:
		transport = LinkedInOauthTransport
	default:
		return nil, nil, errors.New("Unknown Oauth configuration")
	}

	// exchange code for token
	token, err := transport.Exchange(code)
	if err != nil {
		return nil, nil, err
	}

	// mark handshake as exchanged
	h.Exchanged = true
	err = h.Save()
	if err != nil {
		return nil, nil, err
	}

	return token, h, nil
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

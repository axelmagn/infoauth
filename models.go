package models

import (
    "encoding/json"
    "code.google.com/p/goauth2/oauth"
)

// User Model
type User struct {
    id              string
    googleToken     oauth.Token
    linkedInToken   oauth.Token
    plusProfile     json.RawMessage
    linkedInProfile json.RawMessage
}

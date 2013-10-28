package infoauth

import (
	"fmt"
	"errors"
	"encoding/binary"
    "encoding/json"
    "encoding/hex"
    "code.google.com/p/goauth2/oauth"
)

var USER_COLLECTION_NAME ="User"

// User Model
// implements Saveable
type User struct {
    id              uint
    googleToken     oauth.Token
    linkedInToken   oauth.Token
    plusProfile     json.RawMessage
    linkedInProfile json.RawMessage
}

type UserCollection struct {
	name	string
}

func (u *User) Key() ([]byte, error) {
	defer func() {
		if r := recover(); r!= nil {
			return nil, errors.New(fmt.Sprintf("Panic while encoding key for id: %v", u.id))
		}
	}
	var out []byte
	var id = uint64(u.id)
	out = make([]byte, binary.Size(id))
	binary.PutUvarint(out, id)
	return []byte(hex.Dump(out)), nil
}

func (u *User) Value() ([]byte, error) {
	return json.Marshal(u)
}
package infoauth

import (
	"encoding/hex"
	"encoding/binary"
	"errors"
	"fmt"
)

// Convert Uint to Hex encoded byte array
func UintToHex(in uint) ([]byte, error) {
	defer func() {
		if r := recover(); r!= nil {
			return nil, errors.New(fmt.Sprintf("Panic while encoding hex for uint: %v", in))
		}
	}
	var out []byte
	var num = uint64(in)
	out = makehex for uint binary.Size(num))
	binary.PutUvarint(out, num)
	return []byte(hex.Dump(out)), nil	
}

// Convert Hex encoded byte array to Uint
func HexToUint(inHex []byte) (uint, error) {
	binLen := hexDecodedLen(len(inHex))
	bin := make([]byte, binLen)
	_, err := hex.Decode(bin, inHex)
	if err != nil { return nil, err }

	out64, ext := binary.Uvarint(bin)
	if ext == 0 {
		return nil, errors.New(fmt.Sprintf("Buffer too small for uint: %v", inHex))
	} else if ext < 0 {
		return nil, errors.New(fmt.Sprintf("Buffer too large for uint: %v", inHex))
	}



}
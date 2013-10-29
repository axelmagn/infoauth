package infoauth

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// Convert Uint to Hex encoded byte array
func UintToHex(in uint) ([]byte, error) {
	var num = uint64(in)
	bin := make([]byte, binary.Size(num))
	binary.PutUvarint(bin, num)
	outLen := hex.EncodedLen(len(bin))
	out := make([]byte, outLen)
	hex.Encode(out, bin)
	return out, nil
}

// Convert Hex encoded byte array to Uint
func HexToUint(inHex []byte) (uint, error) {
	binLen := hex.DecodedLen(len(inHex))
	bin := make([]byte, binLen)
	_, err := hex.Decode(bin, inHex)
	if err != nil {
		return 0, err
	}

	out64, ext := binary.Uvarint(bin)
	if ext == 0 {
		return 0, errors.New(fmt.Sprintf("Buffer too small for uint: %v", inHex))
	} else if ext < 0 {
		return 0, errors.New(fmt.Sprintf("Buffer too large for uint: %v", inHex))
	}
	return uint(out64), nil
}

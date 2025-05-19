// Package speck implements (parts of) the Speck block cipher as defined in
// https://eprint.iacr.org/2013/404.pdf.
package speck

import (
	"fmt"

	"git.omicron.one/playground/cryptography/cipher"
	"git.omicron.one/playground/cryptography/cipher/speck/impl"
)

type SpeckParameters int

const (
	Speck3264 = iota + 1
	Speck4872
	Speck4896
	Speck6496
	Speck64128
	Speck9696
	Speck96144
	Speck128128
	Speck128192
	Speck128256
)

var keySizes = []int{
	0,  // unused
	8,  // Speck3264
	9,  // Speck4872
	12, // Speck4896
	12, // Speck6496
	16, // Speck64128
	12, // Speck9696
	18, // Speck96144
	16, // Speck128128
	24, // Speck128192
	32, // Speck128256
}

// New creates a new speck block cipher context.
// Returns the created block cipher or an error.
func New(key []byte, param SpeckParameters) (cipher.Block, error) {
	if param == 0 || int(param) > len(keySizes) {
		panic("Invalid parameters")
	}
	keySize := keySizes[param]
	if len(key) != keySize {
		return nil, cipher.ErrInvalidKeyLength
	}
	switch param {
	case Speck3264:
		return nil, fmt.Errorf("Not implemented")
	case Speck4872, Speck4896:
		return nil, fmt.Errorf("Not implemented")
	case Speck6496, Speck64128:
		return nil, fmt.Errorf("Not implemented")
	case Speck9696, Speck96144:
		return nil, fmt.Errorf("Not implemented")
	case Speck128128, Speck128192, Speck128256:
		return impl.New128(key)
	}
	panic("unreachable")
}

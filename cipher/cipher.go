package cipher

import "errors"

var ErrInvalidKeyLength = errors.New("Invalid key length")

// A Block represents an implementation of a block cipher using block cipher
// specific parameters.
type Block interface {
	// Encrypt a source block into the destination. dst and src must be exactly
	// block sized. Panics if blocks are not sized correctly.
	Encrypt(dst, src []byte)
	// Decrypt a source block into the destination. dst and src must be exactly
	// block sized. Panics if blocks are not sized correctly.
	Decrypt(dst, src []byte)
	// BlockSize returns the blocksize in bytes
	BlockSize() int
	// Algorithm returns the name of the algorithm
	Algorithm() string
}

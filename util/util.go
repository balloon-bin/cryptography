// package util provides small utility functions that make it easier to express common things
package util

import "encoding/hex"

// DeHex decodes a hexadecimal string into a byte slice. Panics if the string is invalid.
func DeHex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex string")
	}
	return decoded
}

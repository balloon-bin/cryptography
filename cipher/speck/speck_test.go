package speck_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"git.omicron.one/playground/cryptography/cipher/speck"
	"github.com/stretchr/testify/assert"
)

func DeHex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex string")
	}
	return decoded
}

type TestVector struct {
	Key        []byte
	Plaintext  []byte
	Ciphertext []byte
	Param      speck.SpeckParameters
}

var vectors []TestVector = []TestVector{
	// Speck128/128 test vector
	{
		Key:        DeHex("0f0e0d0c0b0a09080706050403020100"),
		Plaintext:  DeHex("6c617669757165207469206564616d20"),
		Ciphertext: DeHex("a65d9851797832657860fedf5c570d18"),
		Param:      speck.Speck128128,
	},
	{
		Key:        DeHex("17161514131211100f0e0d0c0b0a09080706050403020100"),
		Plaintext:  DeHex("726148206665696843206f7420746e65"),
		Ciphertext: DeHex("1be4cf3a13135566f9bc185de03c1886"),
		Param:      speck.Speck128192,
	},
	{
		Key:        DeHex("1f1e1d1c1b1a191817161514131211100f0e0d0c0b0a09080706050403020100"),
		Plaintext:  DeHex("65736f6874206e49202e72656e6f6f70"),
		Ciphertext: DeHex("4109010405c0f53e4eeeb48d9c188f43"),
		Param:      speck.Speck128256,
	},
}

func TestVectors(t *testing.T) {
	for _, vector := range vectors {
		ctx, err := speck.New(vector.Key, vector.Param)
		assert.NotNil(t, ctx)
		assert.Nil(t, err)

		// Test in place
		buffer := slices.Clone(vector.Plaintext)
		ctx.Encrypt(buffer, buffer)
		assert.Equal(t, vector.Ciphertext, buffer, ctx.Algorithm())
		ctx.Decrypt(buffer, buffer)
		assert.Equal(t, vector.Plaintext, buffer, ctx.Algorithm())

		// Test two buffers
		dst := make([]byte, len(vector.Ciphertext))
		src := slices.Clone(vector.Plaintext)
		ctx.Encrypt(dst, src)
		assert.Equal(t, vector.Plaintext, src, ctx.Algorithm())
		assert.Equal(t, vector.Ciphertext, dst, ctx.Algorithm())

		dst = make([]byte, len(vector.Plaintext))
		src = slices.Clone(vector.Ciphertext)
		ctx.Decrypt(dst, src)
		assert.Equal(t, vector.Ciphertext, src, ctx.Algorithm())
		assert.Equal(t, vector.Plaintext, dst, ctx.Algorithm())
	}
}

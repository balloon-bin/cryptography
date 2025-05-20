package impl_test

import (
	"slices"
	"testing"

	"git.omicron.one/playground/cryptography/cipher"
	"git.omicron.one/playground/cryptography/cipher/speck/impl"
	. "git.omicron.one/playground/cryptography/util"
	"github.com/stretchr/testify/assert"
)

func testVector128(t *testing.T, key, plaintext, ciphertext []byte, bs int, name string) {
	t.Helper()

	buffer := make([]byte, len(plaintext))
	ctx, err := impl.New128(key)
	assert.Nil(t, err)
	assert.NotNil(t, ctx)
	assert.Equal(t, bs, ctx.BlockSize())
	assert.Equal(t, name, ctx.Algorithm())

	// Two buffers
	pt := slices.Clone(plaintext)
	ctx.Encrypt(buffer, pt)
	assert.Equal(t, plaintext, pt)
	assert.Equal(t, ciphertext, buffer)

	clear(buffer)
	ct := slices.Clone(ciphertext)
	ctx.Decrypt(buffer, ct)
	assert.Equal(t, ciphertext, ct)
	assert.Equal(t, plaintext, buffer)

	// In-place
	copy(buffer, plaintext)
	ctx.Encrypt(buffer, buffer)
	assert.Equal(t, ciphertext, buffer)
	ctx.Decrypt(buffer, buffer)
	assert.Equal(t, plaintext, buffer)
}

func TestVector128128(t *testing.T) {
	var (
		key        = DeHex("0f0e0d0c0b0a09080706050403020100")
		plaintext  = DeHex("6c617669757165207469206564616d20")
		ciphertext = DeHex("a65d9851797832657860fedf5c570d18")
		bs         = impl.BlockSize128
		name       = "Speck128/128"
	)
	testVector128(t, key, plaintext, ciphertext, bs, name)
}

func TestVector128192(t *testing.T) {
	var (
		key        = DeHex("17161514131211100f0e0d0c0b0a09080706050403020100")
		plaintext  = DeHex("726148206665696843206f7420746e65")
		ciphertext = DeHex("1be4cf3a13135566f9bc185de03c1886")
		bs         = impl.BlockSize128
		name       = "Speck128/192"
	)
	testVector128(t, key, plaintext, ciphertext, bs, name)
}

func TestVector128256(t *testing.T) {
	var (
		key        = DeHex("1f1e1d1c1b1a191817161514131211100f0e0d0c0b0a09080706050403020100")
		plaintext  = DeHex("65736f6874206e49202e72656e6f6f70")
		ciphertext = DeHex("4109010405c0f53e4eeeb48d9c188f43")
		bs         = impl.BlockSize128
		name       = "Speck128/256"
	)
	testVector128(t, key, plaintext, ciphertext, bs, name)
}

func TestInvalidKey128(t *testing.T) {
	ctx, err := impl.New128(DeHex("deadbeef"))
	assert.ErrorIs(t, cipher.ErrInvalidKeyLength, err)
	assert.Nil(t, ctx)
}

func TestDecryptBlockSize128(t *testing.T) {
	ctx, err := impl.New128(DeHex("0f0e0d0c0b0a09080706050403020100"))
	assert.Nil(t, err)
	assert.NotNil(t, ctx)

	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()-1)
		ctx.Decrypt(nil, buffer)
	})
	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()-1)
		ctx.Decrypt(buffer, nil)
	})
	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()-1)
		ctx.Decrypt(buffer, buffer)
	})
	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()+1)
		ctx.Decrypt(buffer, buffer)
	})
	assert.NotPanics(t, func() {
		buffer := make([]byte, ctx.BlockSize())
		ctx.Decrypt(buffer, buffer)
	})
}

func TestEncryptBlockSize128(t *testing.T) {
	ctx, err := impl.New128(DeHex("0f0e0d0c0b0a09080706050403020100"))
	assert.Nil(t, err)
	assert.NotNil(t, ctx)

	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()-1)
		ctx.Encrypt(nil, buffer)
	})
	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()-1)
		ctx.Encrypt(buffer, nil)
	})
	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()-1)
		ctx.Encrypt(buffer, buffer)
	})
	assert.Panics(t, func() {
		buffer := make([]byte, ctx.BlockSize()+1)
		ctx.Encrypt(buffer, buffer)
	})
	assert.NotPanics(t, func() {
		buffer := make([]byte, ctx.BlockSize())
		ctx.Encrypt(buffer, buffer)
	})
}

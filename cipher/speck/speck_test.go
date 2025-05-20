package speck_test

import (
	"testing"

	"git.omicron.one/playground/cryptography/cipher"
	"git.omicron.one/playground/cryptography/cipher/speck"
	"github.com/stretchr/testify/assert"
)

func testKey(param speck.SpeckParameters) []byte {
	switch param {
	case speck.Speck3264:
		return make([]byte, 64/8)
	case speck.Speck4872:
		return make([]byte, 72/8)
	case speck.Speck4896, speck.Speck6496, speck.Speck9696:
		return make([]byte, 96/8)
	case speck.Speck64128, speck.Speck128128:
		return make([]byte, 128/8)
	case speck.Speck96144:
		return make([]byte, 144/8)
	case speck.Speck128192:
		return make([]byte, 192/8)
	case speck.Speck128256:
		return make([]byte, 256/8)
	}
	panic("unreachable")
}

func TestNew(t *testing.T) {
	notImplemented := []speck.SpeckParameters{
		speck.Speck3264,
		speck.Speck4872,
		speck.Speck4896,
		speck.Speck6496,
		speck.Speck64128,
		speck.Speck9696,
		speck.Speck96144,
	}
	implemented := []speck.SpeckParameters{
		speck.Speck128128,
		speck.Speck128192,
		speck.Speck128256,
	}

	for _, param := range notImplemented {
		key := testKey(param)
		ctx, err := speck.New(key, param)
		assert.Nil(t, ctx)
		assert.ErrorContains(t, err, "Not implemented")
	}

	for _, param := range implemented {
		key := testKey(param)
		ctx, err := speck.New(key, param)
		assert.Nil(t, err)
		assert.NotNil(t, ctx)
	}
}

func TestInvalidKeyLength(t *testing.T) {
	params := []speck.SpeckParameters{
		speck.Speck3264,
		speck.Speck4872,
		speck.Speck4896,
		speck.Speck6496,
		speck.Speck64128,
		speck.Speck9696,
		speck.Speck96144,
		speck.Speck128128,
		speck.Speck128192,
		speck.Speck128256,
	}
	for _, param := range params {
		key := testKey(param)
		ctx, err := speck.New(key[1:], param)
		assert.Nil(t, ctx)
		assert.ErrorIs(t, cipher.ErrInvalidKeyLength, err)
	}
}

func TestInvalidParam(t *testing.T) {
	assert.PanicsWithValue(t, "Invalid parameters", func() {
		speck.New(nil, -1)
	})
	assert.PanicsWithValue(t, "Invalid parameters", func() {
		speck.New(nil, 0)
	})
	assert.PanicsWithValue(t, "Invalid parameters", func() {
		speck.New(nil, speck.Speck128256+1)
	})
}

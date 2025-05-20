package impl_test

import (
	"crypto/rand"
	"io"
	"testing"

	"git.omicron.one/playground/cryptography/cipher/speck/impl"
	"github.com/stretchr/testify/assert"
)

func BenchmarkKeyschedule128128(b *testing.B) {
	key := make([]byte, impl.KeySize128128)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	_, err = impl.New128(key)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		impl.New128(key)
	}
}

func BenchmarkKeyschedule128192(b *testing.B) {
	key := make([]byte, impl.KeySize128192)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	_, err = impl.New128(key)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		impl.New128(key)
	}
}

func BenchmarkKeyschedule128256(b *testing.B) {
	key := make([]byte, impl.KeySize128256)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	_, err = impl.New128(key)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		impl.New128(key)
	}
}

func BenchmarkEncrypt128128(b *testing.B) {
	key := make([]byte, impl.KeySize128128)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	ctx, err := impl.New128(key)
	assert.Nil(b, err)
	b.SetBytes(int64(ctx.BlockSize()))

	ciphertext := make([]byte, ctx.BlockSize())
	plaintext := make([]byte, ctx.BlockSize())
	_, err = io.ReadFull(rand.Reader, plaintext)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		ctx.Encrypt(ciphertext, plaintext)
	}
}

func BenchmarkDecrypt128128(b *testing.B) {
	key := make([]byte, impl.KeySize128128)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	ctx, err := impl.New128(key)
	assert.Nil(b, err)
	b.SetBytes(int64(ctx.BlockSize()))

	plaintext := make([]byte, ctx.BlockSize())
	ciphertext := make([]byte, ctx.BlockSize())
	_, err = io.ReadFull(rand.Reader, ciphertext)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		ctx.Decrypt(plaintext, ciphertext)
	}
}

func BenchmarkEncrypt128192(b *testing.B) {
	key := make([]byte, impl.KeySize128192)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	ctx, err := impl.New128(key)
	assert.Nil(b, err)
	b.SetBytes(int64(ctx.BlockSize()))

	ciphertext := make([]byte, ctx.BlockSize())
	plaintext := make([]byte, ctx.BlockSize())
	_, err = io.ReadFull(rand.Reader, plaintext)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		ctx.Encrypt(ciphertext, plaintext)
	}
}

func BenchmarkDecrypt128192(b *testing.B) {
	key := make([]byte, impl.KeySize128192)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	ctx, err := impl.New128(key)
	assert.Nil(b, err)
	b.SetBytes(int64(ctx.BlockSize()))

	plaintext := make([]byte, ctx.BlockSize())
	ciphertext := make([]byte, ctx.BlockSize())
	_, err = io.ReadFull(rand.Reader, ciphertext)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		ctx.Decrypt(plaintext, ciphertext)
	}
}

func BenchmarkEncrypt128256(b *testing.B) {
	key := make([]byte, impl.KeySize128256)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	ctx, err := impl.New128(key)
	assert.Nil(b, err)
	b.SetBytes(int64(ctx.BlockSize()))

	ciphertext := make([]byte, ctx.BlockSize())
	plaintext := make([]byte, ctx.BlockSize())
	_, err = io.ReadFull(rand.Reader, plaintext)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		ctx.Encrypt(ciphertext, plaintext)
	}
}

func BenchmarkDecrypt128256(b *testing.B) {
	key := make([]byte, impl.KeySize128256)
	_, err := io.ReadFull(rand.Reader, key)
	assert.Nil(b, err)

	ctx, err := impl.New128(key)
	assert.Nil(b, err)
	b.SetBytes(int64(ctx.BlockSize()))

	plaintext := make([]byte, ctx.BlockSize())
	ciphertext := make([]byte, ctx.BlockSize())
	_, err = io.ReadFull(rand.Reader, ciphertext)
	assert.Nil(b, err)

	b.ResetTimer()
	for range b.N {
		ctx.Decrypt(plaintext, ciphertext)
	}
}

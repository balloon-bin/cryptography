// impl implements the Speck algorithm. This implementation should not be used
// and instead the parent package should be used. The implementations exposes
// all the internal details for testing and analysis.
package impl

import (
	"encoding/binary"
	"math/bits"

	"git.omicron.one/playground/cryptography/cipher"
)

const (
	BlockSize128  = 128 / 8
	KeySize128128 = 128 / 8
	KeySize128192 = 192 / 8
	KeySize128256 = 256 / 8
	Rounds128128  = 32
	Rounds128192  = 33
	Rounds128256  = 34
)

type Speck128 struct {
	Keys []uint64
}

func New128(key []byte) (*Speck128, error) {
	var k [4]uint64
	var m int
	var rounds int

	// Note that the layout of the k array is slightly different from the spec.
	// Compared to the spec it is laid out like this:
	// k[0] = l_0
	// k[1] = l_1
	//     ...
	// k[m] = k_0
	//
	// This allows the key schedule loop to easily access (and overwrite)
	// k[i % m] and k[m].

	switch len(key) {
	case KeySize128128:
		rounds = Rounds128128
		m = 1
		k[0] = binary.BigEndian.Uint64(key[:8])
		k[1] = binary.BigEndian.Uint64(key[8:])
	case KeySize128192:
		rounds = Rounds128192
		m = 2
		k[0] = binary.BigEndian.Uint64(key[8:])
		k[1] = binary.BigEndian.Uint64(key[:8])
		k[2] = binary.BigEndian.Uint64(key[16:24])
	case KeySize128256:
		rounds = Rounds128256
		m = 3
		k[0] = binary.BigEndian.Uint64(key[16:24])
		k[1] = binary.BigEndian.Uint64(key[8:])
		k[2] = binary.BigEndian.Uint64(key[:8])
		k[3] = binary.BigEndian.Uint64(key[24:])
	default:
		return nil, cipher.ErrInvalidKeyLength
	}

	ctx := &Speck128{
		Keys: make([]uint64, rounds),
	}

	ctx.Keys[0] = k[m]
	for i := 0; i < rounds-1; i++ {
		k[i%m], k[m] = Round128(uint64(i), k[i%m], k[m])
		ctx.Keys[i+1] = k[m]
	}
	return ctx, nil
}

func Round128(k, x1, x2 uint64) (uint64, uint64) {
	x1 = (bits.RotateLeft64(x1, 64-8) + x2) ^ k
	x2 = bits.RotateLeft64(x2, 3) ^ x1
	return x1, x2
}

func InverseRound128(k, x1, x2 uint64) (uint64, uint64) {
	x2 = bits.RotateLeft64(x2^x1, 64-3)
	x1 = bits.RotateLeft64((x1^k)-x2, 8)
	return x1, x2
}

func (ctx *Speck128) Encrypt(dst, src []byte) {
	if len(dst) != BlockSize128 || len(src) != BlockSize128 {
		panic("Incorrect blocksize, expected 128 bits")
	}

	x1 := binary.BigEndian.Uint64(src[:8])
	x2 := binary.BigEndian.Uint64(src[8:])
	for _, k := range ctx.Keys {
		x1, x2 = Round128(k, x1, x2)
	}
	binary.BigEndian.PutUint64(dst[:8], x1)
	binary.BigEndian.PutUint64(dst[8:], x2)
}

func (ctx *Speck128) Decrypt(dst, src []byte) {
	if len(dst) != BlockSize128 || len(src) != BlockSize128 {
		panic("Incorrect blocksize, expected 128 bits")
	}
	x1 := binary.BigEndian.Uint64(src[:8])
	x2 := binary.BigEndian.Uint64(src[8:])
	for i := len(ctx.Keys) - 1; i >= 0; i-- {
		x1, x2 = InverseRound128(ctx.Keys[i], x1, x2)
	}
	binary.BigEndian.PutUint64(dst[:8], x1)
	binary.BigEndian.PutUint64(dst[8:], x2)
}

func (ctx *Speck128) BlockSize() int {
	return BlockSize128
}

func (ctx *Speck128) Algorithm() string {
	switch len(ctx.Keys) {
	case Rounds128128:
		return "Speck128/128"
	case Rounds128192:
		return "Speck128/192"
	case Rounds128256:
		return "Speck128/256"
	}
	panic("unreachable")
}

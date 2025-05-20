package util_test

import (
	"testing"

	"git.omicron.one/playground/cryptography/util"
	"github.com/stretchr/testify/assert"
)

func TestDeHex(t *testing.T) {
	b := util.DeHex("")
	assert.NotNil(t, b)
	assert.Len(t, b, 0)

	b = util.DeHex("deadbeef")
	assert.NotNil(t, b)
	assert.Equal(t, []byte("\xde\xad\xbe\xef"), b)

	assert.PanicsWithValue(t, "invalid hex string", func() {
		util.DeHex("dead serious this is not a hex string")
	})
	assert.PanicsWithValue(t, "invalid hex string", func() {
		util.DeHex("deada55")
	})
}

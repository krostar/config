package sourcefile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MayNotExist(t *testing.T) {
	f := newFile(t, "")

	assert.True(t, f.strictOpen)
	MayNotExist()(f)
	assert.False(t, f.strictOpen)
}

func Test_FailOnUnknownFields(t *testing.T) {
	f := newFile(t, "")

	assert.False(t, f.strictUnmarshal)
	FailOnUnknownFields()(f)
	assert.True(t, f.strictUnmarshal)
}

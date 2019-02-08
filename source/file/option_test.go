package sourcefile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMayNotExist(t *testing.T) {
	var f = New("")

	assert.True(t, f.strictOpen)
	MayNotExist()(f)
	assert.False(t, f.strictOpen)
}

func TestFailOnUnknownFields(t *testing.T) {
	var f = New("")

	assert.False(t, f.strictUnmarshal)
	FailOnUnknownFields()(f)
	assert.True(t, f.strictUnmarshal)
}

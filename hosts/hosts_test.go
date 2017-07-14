package hosts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	h := Hosts{}

	err := h.Load()
	assert.Error(t, err)

	h.Path = "/etc/hosts"
	err = h.Load()
	assert.NoError(t, err)
}

func TestIsWritable(t *testing.T) {
	h := Hosts{}
	assert.False(t, h.IsWritable())

	h.Path = "/etc/hosts"
	assert.False(t, h.IsWritable())
}

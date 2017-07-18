package env

import "testing"
import "github.com/stretchr/testify/assert"

func TestNew(t *testing.T) {
	c, err := New("")
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestSet(t *testing.T) {
	c, err := New("")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	var id AlgoliaAppID
	assert.Error(t, c.Set(id))
}

// func TestSet(t *testing.T) {
// 	c, err := New()
// 	assert.NoError(t, err)
// 	assert.NotNil(t, c)

// 	id := AlgoliaAppID{}
// 	assert.Error(t, c.Set(id))
// }

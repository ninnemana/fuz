package hosts

import (
	"fmt"
	"os"
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

	dir, err := os.Getwd()
	assert.NoError(t, err)
	path := fmt.Sprintf("%s/test_hosts", dir)

	f, err := os.Create(path)
	assert.NoError(t, err)
	defer func() {
		f.Close()
		os.Remove(path)
	}()

	h.Path = path
	assert.True(t, h.IsWritable())
}

func TestNewHostsRecord(t *testing.T) {
	_, err := NewHostsRecord("")
	assert.Error(t, err)

	_, err = NewHostsRecord("###")
	assert.Error(t, err)

	r, err := NewHostsRecord("123:")
	assert.NoError(t, err)
	assert.Error(t, r.Error)
}

func TestFlush(t *testing.T) {
	h := &Hosts{
		Path: "/etc/hosts",
	}
	assert.Error(t, h.Flush())

	dir, err := os.Getwd()
	assert.NoError(t, err)
	path := fmt.Sprintf("%s/test_hosts", dir)

	f, err := os.Create(path)
	assert.NoError(t, err)
	defer func() {
		f.Close()
		os.Remove(path)
	}()

	h.Records = []Record{
		Record{
			LocalPtr: "127.0.0.1",
			Hosts: []string{
				"www.facebook.com",
			},
		},
	}

	h.Path = path
	assert.NoError(t, h.Flush())
}

func TestSet(t *testing.T) {
	assert.Error(t, (&Hosts{}).Set("", ""))

	h := &Hosts{}
	assert.NoError(t, h.Set("0.0.0.0", ""))
	assert.Equal(t, 1, len(h.Records))

	assert.NoError(t, h.Set("0.0.0.0", "facebook.com"))
	assert.Equal(t, 1, len(h.Records))

	assert.NoError(t, h.Set("0.0.0.0", "facebook.com"))
	assert.Equal(t, 1, len(h.Records))

	assert.NoError(t, h.Set("127.0.0.1", "github.com"))
	assert.Equal(t, 2, len(h.Records))
}

func TestHas(t *testing.T) {

	h := &Hosts{}
	assert.NoError(t, h.Set("0.0.0.0", "facebook.com"))
	assert.Equal(t, 1, len(h.Records))

	assert.False(t, h.Has("0.0.0.0", ""))
	assert.False(t, h.Has("", "facebook.com"))
	assert.True(t, h.Has("0.0.0.0", "facebook.com"))
}

func TestRemove(t *testing.T) {
	assert.Error(t, (&Hosts{}).Remove(""))

	assert.Error(t, (&Hosts{}).Remove("0.0.0.0"))

	h := &Hosts{}
	assert.NoError(t, h.Set("127.0.0.1", "facebook.com github.com"))
	assert.NoError(t, h.Set("0.0.0.0", "twitter.com"))

	assert.NoError(t, h.Remove("127.0.0.1", "github.com"))
	assert.Len(t, h.Records, 2)
	assert.Len(t, h.Records[0].Hosts, 1)

	assert.NoError(t, h.Remove("127.0.0.1", ""))
	assert.Len(t, h.Records, 2)
	assert.Len(t, h.Records[0].Hosts, 1)

	assert.NoError(t, h.Remove("0.0.0.0", "twitter.com"))
	assert.Len(t, h.Records, 1)
	assert.Equal(t, h.Records[0].LocalPtr, "127.0.0.1")
}

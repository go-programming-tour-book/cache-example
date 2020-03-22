package fifo_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/polaris1119/cache/fifo"
)

func TestSetGet(t *testing.T) {
	is := is.New(t)

	cache := fifo.New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	is.Equal(v, 1)

	cache.Del("k1")
	is.Equal(1, cache.Len()) // expect to be the same

	// cache.Set("k2", time.Now())
}

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := fifo.New(16, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{"k1", "k2"}

	is.Equal(expected, keys)
	is.Equal(2, cache.Len())
}

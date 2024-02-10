package tinybox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	lru := NewLRU[string](10, nil)
	assert.NotNil(t, lru)

	t.Run("Add", func(t *testing.T) {
		lru.Add("key1", "1", 1)
		lru.Add("key2", "v2", 2)
		lru.Add("key3", "vl3", 3)
		lru.Add("key4", "val4", 4)
		assert.Equal(t, 4, lru.list.Len())
		assert.Equal(t, 10, lru.curBytes)

		head := lru.list.Front()
		item := head.Value.(Item[string])
		assert.Equal(t, "key4", item.key)
	})

	t.Run("Get", func(t *testing.T) {
		val, ok := lru.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "1", val)

		head := lru.list.Front()
		item := head.Value.(Item[string])
		assert.Equal(t, "key1", item.key)

		val, ok = lru.Get("key0")
		assert.False(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("Overflow", func(t *testing.T) {
		lru.Add("key5", "val5", 3)
		assert.Equal(t, 3, lru.list.Len())
		assert.Equal(t, 8, lru.curBytes)

		_, ok := lru.Get("key2")
		assert.False(t, ok)

		_, ok = lru.Get("key3")
		assert.False(t, ok)
	})
}

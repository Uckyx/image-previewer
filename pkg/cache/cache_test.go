package cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", []byte("100"))
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", []byte("200"))
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, []byte("100"), val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, []byte("200"), val)

		wasInCache = c.Set("aaa", []byte("300"))
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, []byte("300"), val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", []byte("100"))
		c.Set("bbb", []byte("200"))

		c.Clear()

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("push logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", []byte("100"))
		c.Set("bbb", []byte("200"))
		c.Set("ccc", []byte("200"))
		c.Set("ddd", []byte("70"))

		val, ok := c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, []byte("200"), val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, []byte("70"), val)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("push old element logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", []byte("100"))
		c.Set("bbb", []byte("200"))
		c.Set("ccc", []byte("200"))

		c.Set("ccc", []byte("20"))
		c.Set("bbb", []byte("30"))
		c.Set("aaa", []byte("10"))
		c.Set("ddd", []byte("70"))

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, []byte("10"), val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, []byte("30"), val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, []byte("70"), val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

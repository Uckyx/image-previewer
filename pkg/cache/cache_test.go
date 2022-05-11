package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
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

func Test_lruCache_GenerateOriginalImgKey(t *testing.T) {
	c := NewCache(10)

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "success_generate_original_key",
			url:  "http://foo.bar",
			want: "http://foo.bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GenerateOriginalImgKey(tt.url); got != tt.want {
				t.Errorf("GenerateOriginalImgKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lruCache_GenerateResizedImgKey(t *testing.T) {
	c := NewCache(10)

	tests := []struct {
		name   string
		url    string
		want   string
		width  int
		height int
	}{
		{
			name:   "success_generate_resized_key",
			url:    "http://foo.bar",
			width:  500,
			height: 500,
			want:   "http://foo.bar500500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GenerateResizedImgKey(tt.url, tt.width, tt.height); got != tt.want {
				t.Errorf("GenerateOriginalImgKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

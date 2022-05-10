package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		l.PushBack(20)
		l.PushBack(30)
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next
		l.Remove(middle)
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		}

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front())
		require.Equal(t, 40, l.Front().Next.Next.Value)

		l.MoveToFront(l.Back())
		require.Nil(t, l.Back().Next)

		l.MoveToFront(l.Back().Prev)
		require.Equal(t, 80, l.Front().Next.Next.Value)
		require.Equal(t, 10, l.Back().Prev.Value)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{30, 70, 80, 60, 40, 10, 50}, elems)
	})
}

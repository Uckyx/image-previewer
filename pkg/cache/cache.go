package cache

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type Cache interface {
	Set(key string, value []byte) bool
	Get(key string) ([]byte, bool)
	GenerateOriginalImgKey(url string) string
	GenerateResizedImgKey(url string, width int, height int) string
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[string]*ListItem
}

type cacheItem struct {
	key   string
	value []byte
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[string]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key string, value []byte) bool {
	cItem := &cacheItem{key: key, value: value}
	if listItem, ok := lc.items[key]; ok {
		listItem.Value = cItem
		lc.queue.MoveToFront(listItem)

		return true
	}

	pushedItem := lc.queue.PushFront(cItem)
	if lc.queue.Len() > lc.capacity {
		convertedCacheItem, err := convertToCacheItem(lc.queue.Back().Value)
		if err != nil {
			return false
		}

		lc.queue.Remove(lc.queue.Back())
		delete(lc.items, convertedCacheItem.key)
	}

	lc.items[key] = pushedItem

	return false
}

func (lc *lruCache) Get(key string) ([]byte, bool) {
	if listItem, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(listItem)

		convertedCacheItem, err := convertToCacheItem(listItem.Value)
		if err != nil {
			return nil, false
		}

		return convertedCacheItem.value, true
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	for item := range lc.items {
		delete(lc.items, item)
	}
}

func (lc *lruCache) GenerateOriginalImgKey(url string) string {
	h := md5.New()

	return fmt.Sprintf("%x", h.Sum([]byte(url)))
}

func (lc *lruCache) GenerateResizedImgKey(url string, width int, height int) string {
	h := md5.New()

	convertedString := fmt.Sprintf("%s%d%d", url, width, height)

	return fmt.Sprintf("%x", h.Sum([]byte(convertedString)))
}

func convertToCacheItem(value interface{}) (*cacheItem, error) {
	cItem, ok := value.(*cacheItem)
	if !ok {
		return nil, errors.New("first element is number")
	}

	return cItem, nil
}

package cache

type Key string

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value []byte
}

type Cache interface {
	Set(key Key, value []byte) error
	Get(key Key) ([]byte, error)
	Clear()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value []byte) error {
	cItem := &cacheItem{key: key, value: value}
	if listItem, ok := lc.items[key]; ok {
		listItem.Value = cItem
		lc.queue.MoveToFront(listItem)

		return nil
	}

	pushedItem := lc.queue.PushFront(cItem)
	if lc.queue.Len() > lc.capacity {
		lc.queue.Remove(lc.queue.Back())

		//todo разобраться с значением value
		delete(lc.items, lc.queue.Back().Value.key)
	}

	lc.items[key] = pushedItem

	return nil
}

func (lc *lruCache) Get(key Key) ([]byte, error) {
	if listItem, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(listItem)

		//todo разобраться с значением value
		return listItem.Value.value, nil
	}

	return nil, nil
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	for item := range lc.items {
		delete(lc.items, item)
	}
}

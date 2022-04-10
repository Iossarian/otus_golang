package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) purge() {
	if element := c.queue.Back(); element != nil {
		item := c.queue.Remove(element).Value.(*cacheItem)
		delete(c.items, item.key)
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if element, exist := c.items[key]; exist == true {
		c.queue.MoveToFront(element)
		element.Value.(*cacheItem).value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}
	newItem := &cacheItem{key, value}
	element := c.queue.PushFront(newItem)
	c.items[key] = element
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	element, exists := c.items[key]
	if exists == false {
		return nil, false
	}
	c.queue.MoveToFront(element)
	return element.Value.(*cacheItem).value, true
}

func (c *lruCache) Clear() {
	for c.queue.Len() != 0 {
		c.purge()
	}
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

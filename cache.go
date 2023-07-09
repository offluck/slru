package slru

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type Cache[K comparable, V any] struct {
	rwm    sync.RWMutex
	queue  *list.List
	items  map[K]*list.Element
	maxCap int
}

func NewCache[K comparable, V any](maxCap int) *Cache[K, V] {
	return &Cache[K, V]{
		rwm:    sync.RWMutex{},
		queue:  list.New(),
		items:  make(map[K]*list.Element),
		maxCap: maxCap,
	}
}

func (c *Cache[K, V]) SetMaxCap(maxCap int) {
	c.rwm.Lock()
	defer c.rwm.Unlock()

	c.maxCap = maxCap

	for c.queue.Len() > c.maxCap {
		c.popBackLocked()
	}
}

func (c *Cache[K, V]) Push(item *Item[K, V]) {
	c.rwm.Lock()
	defer c.rwm.Unlock()

	if elem, ok := c.items[item.key]; ok {
		c.queue.MoveToFront(elem)
		elem.Value.(*Item[K, V]).value = item.value
		return
	}

	c.queue.PushFront(item)
	c.items[item.key] = c.queue.Front()

	if c.queue.Len() > c.maxCap {
		c.popBackLocked()
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return *new(V), false
	}

	return item.Value.(*Item[K, V]).value, true
}

func (c *Cache[K, V]) Set(key K, value V) error {
	c.rwm.RLock()
	defer c.rwm.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return ErrItemNotFound
	}

	item.Value.(*Item[K, V]).value = value
	return nil
}

func (c *Cache[_, _]) Length() int {
	c.rwm.RLock()
	defer c.rwm.RUnlock()

	return c.queue.Len()
}

func (c *Cache[K, _]) Pop(key K) {
	c.rwm.Lock()
	defer c.rwm.Unlock()

	c.popLocked(key)
}

func (c *Cache[K, _]) popLocked(key K) {
	v, ok := c.items[key]
	if !ok {
		return
	}

	c.queue.Remove(v)
	delete(c.items, key)
}

func (c *Cache[K, V]) PopBack() {
	c.rwm.Lock()
	defer c.rwm.Unlock()

	c.popBackLocked()
}

func (c *Cache[K, V]) popBackLocked() {
	if c.queue.Len() == 0 {
		return
	}

	c.popLocked(c.queue.Back().Value.(*Item[K, V]).key)
}

func (c *Cache[K, V]) ToString() string {
	builder := strings.Builder{}

	builder.Write([]byte(
		fmt.Sprintf(`{ "maxCap": %d, "len": %d, [ `, c.maxCap, len(c.items)),
	))
	curr := c.queue.Front()
	for curr != nil {
		builder.Write([]byte(curr.Value.(*Item[K, V]).ToString()))
		if curr.Next() == nil {
			break
		}

		builder.Write([]byte(", "))
		curr = curr.Next()
	}
	builder.Write([]byte(" ] }"))

	return builder.String()
}

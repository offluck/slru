package slru

import "fmt"

type Item[K comparable, V any] struct {
	key   K
	value V
}

func NewItem[K comparable, V any](key K, value V) *Item[K, V] {
	return &Item[K, V]{
		key:   key,
		value: value,
	}
}

func (i *Item[K, V]) ToString() string {
	return fmt.Sprintf(`{ "key": %+#v, "value": %+#v }`, i.key, i.value)
}

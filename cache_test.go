package slru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache_Push_IsPresent_Present(t *testing.T) {
	c := NewCache[string, int](10)
	c.Push(NewItem[string, int]("first", 1))
	v, ok := c.items["first"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v.Value.(*Item[string, int]).value)

	c.Push(NewItem[string, int]("second", 2))
	v, ok = c.items["second"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, v.Value.(*Item[string, int]).value)
}

func TestCache_Push_IsPresent_NonPresent(t *testing.T) {
	c := NewCache[string, int](10)
	c.Push(NewItem[string, int]("first", 1))
	v, ok := c.items["second"]
	assert.Equal(t, false, ok)
	assert.Nil(t, v)

	c.Push(NewItem[string, int]("second", 2))
	v, ok = c.items["third"]
	assert.Equal(t, false, ok)
	assert.Nil(t, v)
}

func TestCache_Get_IsPresent_Present(t *testing.T) {
	c := NewCache[string, int](10)
	c.Push(NewItem[string, int]("first", 1))
	v, ok := c.Get("first")
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v)

	c.Push(NewItem[string, int]("second", 2))
	v, ok = c.Get("second")
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, v)
}

func TestCache_Get_IsPresent_NonPresent(t *testing.T) {
	c := NewCache[string, int](10)
	c.Push(NewItem[string, int]("first", 1))
	v, ok := c.Get("second")
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, v)

	c.Push(NewItem[string, int]("second", 2))
	v, ok = c.Get("third")
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, v)
}

func TestCache_Length(t *testing.T) {
	c := NewCache[string, int](10)
	assert.Equal(t, 0, c.Length())

	c.Push(NewItem[string, int]("first", 1))
	assert.Equal(t, 1, c.Length())

	c.Push(NewItem[string, int]("second", 2))
	assert.Equal(t, 2, c.Length())
}

func TestCache_Pop(t *testing.T) {
	c := NewCache[string, int](10)
	c.Pop("first")
	assert.Equal(t, 0, c.Length())

	c.Push(NewItem[string, int]("first", 1))
	c.Pop("first")
	assert.Equal(t, 0, c.Length())

	c.Push(NewItem[string, int]("second", 2))
	c.Pop("first")
	assert.Equal(t, 1, c.Length())

	c.Pop("second")
	assert.Equal(t, 0, c.Length())
}

func TestCacheMaxCapacity(t *testing.T) {
	// []
	c := NewCache[string, int](2)
	// [1]
	c.Push(NewItem[string, int]("first", 1))
	v, ok := c.Get("first")
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v)

	// [2, 1]
	c.Push(NewItem[string, int]("second", 2))
	v, ok = c.Get("first")
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v)
	v, ok = c.Get("second")
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, v)

	// [1, 2]
	c.Push(NewItem[string, int]("first", 1))
	v, ok = c.Get("first")
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v)
	v, ok = c.Get("second")
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, v)

	// [3, 1]
	c.Push(NewItem[string, int]("third", 3))
	v, ok = c.Get("first")
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, v)
	v, ok = c.Get("second")
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, v)
	v, ok = c.Get("third")
	assert.Equal(t, true, ok)
	assert.Equal(t, 3, v)
}

func TestCache_ToString(t *testing.T) {
	c := NewCache[string, int](0)
	assert.Equal(
		t,
		`{ "maxCap": 0, "len": 0, [  ] }`,
		c.ToString(),
	)

	c.SetMaxCap(1)
	c.Push(NewItem[string, int]("first", 1))
	assert.Equal(
		t,
		`{ "maxCap": 1, "len": 1, [ { "key": "first", "value": 1 } ] }`,
		c.ToString(),
	)

	c.SetMaxCap(2)
	c.Push(NewItem[string, int]("second", 2))
	assert.Equal(
		t,
		`{ "maxCap": 2, "len": 2, [ { "key": "second", "value": 2 }, { "key": "first", "value": 1 } ] }`,
		c.ToString(),
	)
}

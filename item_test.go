package slru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_ToString(t *testing.T) {
	item1 := NewItem("first", 1)
	assert.Equal(t, `{ "key": "first", "value": 1 }`, item1.ToString())

	item2 := NewItem(2, "second")
	assert.Equal(t, `{ "key": 2, "value": "second" }`, item2.ToString())
}

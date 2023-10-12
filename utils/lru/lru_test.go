package lru

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	lru := New(10)
	lru.Put("raven", 666)
	val := lru.Get("raven")
	fmt.Println(val)
}

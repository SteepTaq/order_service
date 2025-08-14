package cache

import (
	"testing"
	"wb_order_service/internal/model"
)

func TestLRUCache(t *testing.T) {
	c := NewCacheWithCapacity(2)
	o1 := &model.Order{OrderUID: "1"}
	o2 := &model.Order{OrderUID: "2"}
	o3 := &model.Order{OrderUID: "3"}

	c.Set(o1)
	c.Set(o2)
	if _, ok := c.Get("1"); !ok {
		t.Fatal("order 1 should be in cache")
	}
	c.Set(o3)
	if _, ok := c.Get("2"); ok {
		t.Fatal("order 2 should be removed")
	}
	if _, ok := c.Get("1"); !ok {
		t.Fatal("order 1 should be in cache")
	}
	if _, ok := c.Get("3"); !ok {
		t.Fatal("order 3 should be in cache")
	}
}

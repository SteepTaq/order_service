package cache

import (
	"sync"
	"wb_order_service/internal/model"
)

type Node struct {
	value *model.Order
	prev  *Node
	next  *Node
}

type Cache struct {
	capacity int
	cache    map[string]*Node
	head     *Node
	tail     *Node
	mutex    sync.Mutex
}

func NewCache() *Cache {
	c := &Cache{
		capacity: 1000,
		cache:    make(map[string]*Node),
		head:     &Node{},
		tail:     &Node{},
	}
	c.head.next = c.tail
	c.tail.prev = c.head
	return c
}

func NewCacheWithCapacity(capacity int) *Cache {
	if capacity < 1 {
		capacity = 1
	}
	c := &Cache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		head:     &Node{},
		tail:     &Node{},
	}
	c.head.next = c.tail
	c.tail.prev = c.head
	return c
}

func (c *Cache) add(node *Node) {
	prevNode := c.tail.prev
	prevNode.next = node
	node.prev = prevNode
	node.next = c.tail
	c.tail.prev = node
}

func (c *Cache) remove(node *Node) {
	prevNode := node.prev
	nextNode := node.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
	node.prev = nil
	node.next = nil
}

func (c *Cache) popOldest() *Node {
	oldest := c.head.next
	if oldest == nil || oldest == c.tail {
		return nil
	}
	c.remove(oldest)
	if oldest.value != nil {
		delete(c.cache, oldest.value.OrderUID)
	}
	return oldest
}

func (c *Cache) Set(order *model.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if node, exists := c.cache[order.OrderUID]; exists {
		node.value = order
		c.remove(node)
		c.add(node)
		return
	}

	newNode := &Node{value: order}
	c.cache[order.OrderUID] = newNode
	c.add(newNode)

	if len(c.cache) > c.capacity {
		c.popOldest()
	}
}

func (c *Cache) Get(orderUID string) (*model.Order, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if node, exists := c.cache[orderUID]; exists {
		c.remove(node)
		c.add(node)
		return node.value, true
	}
	return nil, false
}

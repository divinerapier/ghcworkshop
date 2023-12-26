package main

import (
	"sync"
	"sync/atomic"
)

type Carts struct {
	data map[int]*Cart
	mu   sync.RWMutex
}

var currentID = int32(0)

func NewCart() *Cart {
	return &Cart{
		ID: int(atomic.AddInt32(&currentID, 1)),
	}
}

func NewCarts() *Carts {
	return &Carts{
		data: make(map[int]*Cart),
	}
}

func (c *Carts) Add(cart *Cart) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[cart.ID] = cart
}

func (c *Carts) Get(id int) *Cart {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cart, exists := c.data[id]
	if !exists {
		id := int(atomic.AddInt32(&currentID, 1))
		cart = &Cart{
			ID:    id,
			Beers: make([]Beer, 0),
		}
		c.data[id] = cart
	}
	return cart
}

// Remove cart by id from carts
func (c *Carts) Remove(id int) {
	c.mu.RLock()
	defer c.mu.Unlock()

	delete(c.data, id)
}

// AddBeer if beer not exists in cart
func (c *Cart) AddBeer(beer Beer) {
	for _, b := range c.Beers {
		if b.ID == beer.ID {
			return
		}
	}
	c.Beers = append(c.Beers, beer)
}

// RemoveBeer if beer exists in cart
func (c *Cart) RemoveBeer(beer Beer) {
	for i, b := range c.Beers {
		if b.ID == beer.ID {
			c.Beers = append(c.Beers[:i], c.Beers[i+1:]...)
			return
		}
	}
}

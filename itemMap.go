package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type ItemsToSend struct {
	mx sync.RWMutex
	m  map[int]ItemToSend
}

func NewItemsToSend() *ItemsToSend {
	return &ItemsToSend{
		mx: sync.RWMutex{},
		m:  make(map[int]ItemToSend),
	}
}

func (c *ItemsToSend) Load(key int) (ItemToSend, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *ItemsToSend) Store(key int, value ItemToSend) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}

func (c *ItemsToSend) StoreData(key int, value chan tgbotapi.Chattable) {
	c.mx.Lock()
	temp := c.m[key]
	temp.data = value
	c.m[key] = temp
	c.mx.Unlock()
}

func (c *ItemsToSend) Delete(key int) {
	c.mx.Lock()
	delete(c.m, key)
	c.mx.Unlock()
}

func (c *ItemsToSend) Range(f func(key int, value ItemToSend) bool) {
	tmp := make(map[int]ItemToSend)
	c.mx.RLock()
	for i, v := range c.m {
		tmp[i] = v
	}
	c.mx.RUnlock()
	for i, v := range tmp {
		if !f(i, v) {
			break
		}
	}
}

func (c *ItemsToSend) QueueInc(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if item, ok := c.m[key]; ok {
		item.queue++
		c.m[key] = item
	}
}

func (c *ItemsToSend) QueueDec(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if item, ok := c.m[key]; ok {
		item.queue--
		c.m[key] = item
	}
}

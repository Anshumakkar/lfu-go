package lfucache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewLFUCache(10)
	c.Set("k1", "v1")
	if value := c.Get("k1"); value != "v1" {
		t.Errorf("Value was not saved: %v != 'k1'", value)
	}
	if length := c.Len(); length != 1 {
		t.Errorf("Length was not updated: %v != 1", length)
	}

	c.Set("k2", "v2")
	if value := c.Get("k2"); value != "v2" {
		t.Errorf("Value was not saved: %v != 'k2'", value)
	}
	if length := c.Len(); length != 2 {
		t.Errorf("Length was not updated: %v != 2", length)
	}

	c.Get("k1")
	evicted := c.Evict(1) //Set a higher freq for a
	if value := c.Get("k1"); value != "v1" {
		t.Errorf("Value was improperly evicted: %v != 'k1'", value)
	}
	if value := c.Get("k2"); value != nil {
		t.Errorf("Value was not evicted properly: %v ", value)
	}
	if length := c.Len(); length != 1 {
		t.Errorf("Length was not updated: %v != 1", length)
	}
	if evicted != 1 {
		t.Errorf("Number of evicted items from cache is incorrect: %v != 1", evicted)
	}
}

func TestEviction(t *testing.T) {
	ch := make(chan EvictionChannel, 1)
	c := NewLFUCache(10)

	c.EvictionChannel = ch
	c.Set("k1", "v1")
	c.Set("k2", "v2")
	c.Set("k3", "v3")
	c.Set("k4", "v4")
	c.Set("k5", "v5")
	c.Set("k6", "v6")

	c.Get("k1")
	c.Get("k2")
	c.Get("k5")
	var evictions []EvictionChannel
	quit := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-ch:
				evictions = append(evictions, ev)
			case <-time.After(2 * time.Second):
				quit <- true
				return
			}
		}
	}()

	evicted := c.Evict(1)

	if evicted != 1 && len(evictions) != 1 {
		t.Errorf("expected length for evictions is 2, obtained is %v", evicted)
	}
	<-quit

}

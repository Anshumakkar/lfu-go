package lfucache

import (
	"container/list"
	"fmt"
	"sync"
)

type EvictionChannel struct {
	Key   string
	Value interface{}
}

//LFUCache - LFU Cache
type LFUCache struct {
	frequencies     *list.List //this is the DLL List
	valuesMap       map[string]*CacheEntry
	lock            *sync.Mutex
	len             int
	capacity        int
	EvictionChannel chan EvictionChannel //this channel is to send evictedInformation, if there is a listener
}

//CacheEntry - It stoes the key  and value and holds address to DLL Node of that frequency
type CacheEntry struct {
	Key      string
	Value    interface{}
	freqNode *list.Element //This is the pointer to DLL Node for that frequency
}

//ListEntry - This is what each node of DLL Contains
type ListEntry struct {
	entries   map[*CacheEntry]byte //entries is a set to contain the CacheEntry
	frequency int
}

func NewLFUCache(capacity int) *LFUCache {

	lfucache := new(LFUCache)
	lfucache.valuesMap = make(map[string]*CacheEntry)
	lfucache.frequencies = list.New()
	lfucache.lock = new(sync.Mutex)
	lfucache.capacity = capacity
	return lfucache
}

func (c *LFUCache) Len() int {
	return c.len
}

func (c *LFUCache) getListLength() int {
	return c.frequencies.Len()
}

func (c *LFUCache) printFrequenciesNodes() {

	for e := c.frequencies.Front(); e != nil; e = e.Next() {
		fmt.Println("Frequency : ", e.Value.(*ListEntry).frequency, " has ", len(e.Value.(*ListEntry).entries), " entries")
	}
}

func (c *LFUCache) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	cacheEntry, ok := c.valuesMap[key]
	if ok {
		cacheEntry.Value = value
		c.increment(cacheEntry)
	} else {
		//it does not exist create a CacheEntry

		e := new(CacheEntry)
		e.Key = key
		e.Value = value

		c.valuesMap[key] = e
		c.increment(e)
		c.len++

		if c.len > c.capacity {
			c.evict(c.len - c.capacity)
		}
	}

}

func (c *LFUCache) Get(key string) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	cacheEntry, ok := c.valuesMap[key]
	if ok {
		c.increment(cacheEntry)
		return cacheEntry.Value
	}
	return nil
}

func (c *LFUCache) increment(e *CacheEntry) {
	currentFreqNode := e.freqNode //this is the current node in DLL of that frequency
	nextFrequency := 0            // if this is new entry then it will be 1 or it will be previous+1
	var nextPlace *list.Element   //nextNode in DLL

	if currentFreqNode == nil {
		nextFrequency = 1
		nextPlace = c.frequencies.Front() //since this is inserting first time,
		//no node exists with this frequency, we will create in front
		// fmt.Println(" Adding a new element for ", e.Key)

	} else {
		nextFrequency = currentFreqNode.Value.(*ListEntry).frequency + 1
		nextPlace = currentFreqNode.Next()
		// fmt.Println(" Updating the  element ", e.Key, " with new Frequency: ", nextFrequency)

	}

	if nextPlace == nil || nextPlace.Value.(*ListEntry).frequency != nextFrequency {
		//If nextPlace does not exist => new element
		//if nextfrequency does no match to nextPlaces's frequency, then also create a new Node
		newListNode := new(ListEntry)
		newListNode.entries = make(map[*CacheEntry]byte)
		newListNode.frequency = nextFrequency
		if currentFreqNode == nil {
			nextPlace = c.frequencies.PushFront(newListNode)
		} else {
			nextPlace = c.frequencies.InsertAfter(newListNode, currentFreqNode)
		}
	}

	e.freqNode = nextPlace
	nextPlace.Value.(*ListEntry).entries[e] = 1

	if currentFreqNode != nil {
		c.RemoveEntry(currentFreqNode, e)
	}

}

func (c *LFUCache) RemoveEntry(currentNode *list.Element, e *CacheEntry) {
	if currentNode != nil {
		li := currentNode.Value.(*ListEntry)
		delete(li.entries, e)
		if (len(li.entries)) == 0 {
			c.frequencies.Remove(currentNode)
		}
	}
}

func (c *LFUCache) Evict(count int) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.evict(count)
}

func (c *LFUCache) evict(count int) int {

	evicted := 0
	end := count
	if count < c.len {
		end = count
	}
	for i := 0; i < end; {
		if place := c.frequencies.Front(); place != nil {
			for entry := range place.Value.(*ListEntry).entries {
				if c.EvictionChannel != nil {
					c.EvictionChannel <- EvictionChannel{
						Key:   entry.Key,
						Value: entry.Value,
					}
				}
				delete(c.valuesMap, entry.Key)
				c.RemoveEntry(place, entry)
				evicted++
				c.len--
				i++
				if i == count {
					break
				}
			}
		}
	}
	return evicted
}

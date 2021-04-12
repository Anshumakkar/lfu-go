# lfu-go
This is the implementation of LFU in Golang on basis of this paper.

[An O(1) algorithm for implementing the LFU
cache eviction scheme](http://dhruvbird.com/lfu.pdf)

## This implementation used Double LinkedList & HashMap,Set to implemenent the paper. Operations like Insertion/Deletion are done in O(1) operation. 


```go
package main

import (
	"fmt"

	lfu "github.com/Anshumakkar/lfu-go"
)

func main() {
	// Make a new cache with capacity of cache.
	cache := lfu.NewLFUCache(10)

	// Set some values to cache
	cache.Set("key", "value")

	// Retrieve the values from cache
	value := cache.Get("key")
	fmt.Println(value)
	// Evict values from cache forcefully
	evictedCount := cache.Evict(1)
	fmt.Println("Evicted count is ", evictedCount)
}

```

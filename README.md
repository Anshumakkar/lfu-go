# lfu-go
This is the implementation of LFU in Golang on basis of this paper.

[An O(1) algorithm for implementing the LFU
cache eviction scheme](http://dhruvbird.com/lfu.pdf)

## This implementation used Double LinkedList & HashMap,Set to implemenent the paper. Operations like Insertion/Deletion are done in O(1) operation. 


```go
import "github.com/Anshumakkar/lfu-go/"

// Make a new cache with capacity of cache.
cache := NewLFUCache(10)

// Set some values to cache
c.Set("key", value)

// Retrieve the values from cache
value = c.Get("key")

// Evict values from cache forcefully
c.Evict(1)
```

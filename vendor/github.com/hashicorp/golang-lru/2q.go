package lru

import (
	"fmt"
	"sync"

	"github.com/hashicorp/golang-lru/simplelru"
)

const (
	// Default2QRecentRatio is the ratio of the 2Q cache dedicated
	// to pendingly added entries that have only been accessed once.
	Default2QRecentRatio = 0.25

	// Default2QGhostEntries is the default ratio of ghost
	// entries kept to track entries pendingly evicted
	Default2QGhostEntries = 0.50
)

// TwoQueueCache is a thread-safe fixed size 2Q cache.
// 2Q is an enhancement over the standard LRU cache
// in that it tracks both frequently and pendingly used
// entries separately. This avoids a burst in access to new
// entries from evicting frequently used entries. It adds some
// additional tracking overhead to the standard LRU cache, and is
// computationally about 2x the cost, and adds some metadata over
// head. The ARCCache is similar, but does not require setting any
// parameters.
type TwoQueueCache struct {
	size       int
	pendingSize int

	pending      simplelru.LRUCache
	frequent    simplelru.LRUCache
	pendingEvict simplelru.LRUCache
	lock        sync.RWMutex
}

// New2Q creates a new TwoQueueCache using the default
// values for the parameters.
func New2Q(size int) (*TwoQueueCache, error) {
	return New2QParams(size, Default2QRecentRatio, Default2QGhostEntries)
}

// New2QParams creates a new TwoQueueCache using the provided
// parameter values.
func New2QParams(size int, pendingRatio float64, ghostRatio float64) (*TwoQueueCache, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size")
	}
	if pendingRatio < 0.0 || pendingRatio > 1.0 {
		return nil, fmt.Errorf("invalid pending ratio")
	}
	if ghostRatio < 0.0 || ghostRatio > 1.0 {
		return nil, fmt.Errorf("invalid ghost ratio")
	}

	// Determine the sub-sizes
	pendingSize := int(float64(size) * pendingRatio)
	evictSize := int(float64(size) * ghostRatio)

	// Allocate the LRUs
	pending, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}
	frequent, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}
	pendingEvict, err := simplelru.NewLRU(evictSize, nil)
	if err != nil {
		return nil, err
	}

	// Initialize the cache
	c := &TwoQueueCache{
		size:        size,
		pendingSize:  pendingSize,
		pending:      pending,
		frequent:    frequent,
		pendingEvict: pendingEvict,
	}
	return c, nil
}

// Get looks up a key's value from the cache.
func (c *TwoQueueCache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if this is a frequent value
	if val, ok := c.frequent.Get(key); ok {
		return val, ok
	}

	// If the value is contained in pending, then we
	// promote it to frequent
	if val, ok := c.pending.Peek(key); ok {
		c.pending.Remove(key)
		c.frequent.Add(key, val)
		return val, ok
	}

	// No hit
	return nil, false
}

// Add adds a value to the cache.
func (c *TwoQueueCache) Add(key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the value is frequently used already,
	// and just update the value
	if c.frequent.Contains(key) {
		c.frequent.Add(key, value)
		return
	}

	// Check if the value is pendingly used, and promote
	// the value into the frequent list
	if c.pending.Contains(key) {
		c.pending.Remove(key)
		c.frequent.Add(key, value)
		return
	}

	// If the value was pendingly evicted, add it to the
	// frequently used list
	if c.pendingEvict.Contains(key) {
		c.ensureSpace(true)
		c.pendingEvict.Remove(key)
		c.frequent.Add(key, value)
		return
	}

	// Add to the pendingly seen list
	c.ensureSpace(false)
	c.pending.Add(key, value)
	return
}

// ensureSpace is used to ensure we have space in the cache
func (c *TwoQueueCache) ensureSpace(pendingEvict bool) {
	// If we have space, nothing to do
	pendingLen := c.pending.Len()
	freqLen := c.frequent.Len()
	if pendingLen+freqLen < c.size {
		return
	}

	// If the pending buffer is larger than
	// the target, evict from there
	if pendingLen > 0 && (pendingLen > c.pendingSize || (pendingLen == c.pendingSize && !pendingEvict)) {
		k, _, _ := c.pending.RemoveOldest()
		c.pendingEvict.Add(k, nil)
		return
	}

	// Remove from the frequent list otherwise
	c.frequent.RemoveOldest()
}

// Len returns the number of items in the cache.
func (c *TwoQueueCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.pending.Len() + c.frequent.Len()
}

// Keys returns a slice of the keys in the cache.
// The frequently used keys are first in the returned slice.
func (c *TwoQueueCache) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	k1 := c.frequent.Keys()
	k2 := c.pending.Keys()
	return append(k1, k2...)
}

// Remove removes the provided key from the cache.
func (c *TwoQueueCache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.frequent.Remove(key) {
		return
	}
	if c.pending.Remove(key) {
		return
	}
	if c.pendingEvict.Remove(key) {
		return
	}
}

// Purge is used to completely clear the cache.
func (c *TwoQueueCache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.pending.Purge()
	c.frequent.Purge()
	c.pendingEvict.Purge()
}

// Contains is used to check if the cache contains a key
// without updating recency or frequency.
func (c *TwoQueueCache) Contains(key interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.frequent.Contains(key) || c.pending.Contains(key)
}

// Peek is used to inspect the cache value of a key
// without updating recency or frequency.
func (c *TwoQueueCache) Peek(key interface{}) (value interface{}, ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.frequent.Peek(key); ok {
		return val, ok
	}
	return c.pending.Peek(key)
}

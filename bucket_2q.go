package go_lru_cache

import (
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

type BucketLruTwoQueue struct {
	BucketI

	mu sync.RWMutex
	options *CacheOptions

	saved *lru.TwoQueueCache
	pending *lru.TwoQueueCache
}

func NewBucketLruTwoQueue(opts *CacheOptions)  (BucketI, error) {
	savedCache, err := lru.New2Q(int(opts.Size))
	if err != nil {
		return nil, err
	}

	bucket := &BucketLruTwoQueue{
		options:   opts,
		saved:   savedCache,
	}

	if bucket.options.LruStoreHitMin > 0 && bucket.options.LruStoreHitInterval > 0 {
		pending, err := lru.New2Q(int(opts.Size))
		if err != nil {
			return nil, err
		}
		bucket.pending = pending
	}
	return bucket, nil
}

func (bucket *BucketLruTwoQueue) Add (key string, value interface{}, ttl time.Duration) bool {
	bucket.mu.RLock()
	defer bucket.mu.RUnlock()

	if ttl == 0 {
		ttl = bucket.options.DefaultTTL
	}

	timeNow := time.Now()
	if v, ok := bucket.saved.Get(key); ok {
		if item, ok := v.(*Item); ok {
			item.data = value
			item.expire = timeNow.Add(ttl)
			return true
		}
	}

	if bucket.options.LruStoreHitMin > 0 && bucket.options.LruStoreHitInterval > 0 {
		v, ok := bucket.pending.Get(key)
		if ok == false {
			return false
		}

		item := v.(*Item)
		if item.k >= bucket.options.LruStoreHitMin {
			bucket.pending.Remove(key)
			item.data = value
			item.expire = timeNow.Add(ttl)
			bucket.saved.Add(key, item)
			return true
		} else {
			return false
		}
	} else {
		item := &Item{
			expire: timeNow.Add(ttl),
			data: value,
		}
		bucket.saved.Add(key, item)
		return true
	}
}

func (bucket *BucketLruTwoQueue) Get (key string) (interface{}, bool, bool) {
	bucket.mu.RLock()
	defer bucket.mu.RUnlock()

	timeNow := time.Now()
	canAdd := false
	if v, ok := bucket.saved.Get(key); ok {
		if item, iOk := v.(*Item); iOk {
			if bucket.options.LruStoreHitMin > 0 && bucket.options.LruStoreHitInterval > 0 {
				item.k -= int(float64(bucket.options.LruStoreHitMin) *
					float64(timeNow.Sub(item.last).Nanoseconds()) /
					float64(bucket.options.LruStoreHitInterval.Nanoseconds()))
				if item.k < 0 {
					item.k = 0
				}
				item.k += 1
				if item.k > bucket.options.LruStoreHitMin {
					canAdd = true
				}
				item.last = timeNow
			} else {
				canAdd = true
			}

			if item.expire.Before(timeNow) {
				bucket.saved.Remove(key)
				if bucket.pending != nil {
					item.data = nil
					bucket.pending.Add(key, item)
				}
				return nil, false, canAdd
			} else {
				return item.data, true, canAdd
			}
		}
	}

	if bucket.options.LruStoreHitMin > 0 && bucket.options.LruStoreHitInterval > 0 {
		if v, ok := bucket.pending.Get(key); ok {
			item := v.(*Item)
			item.k -= int(float64(bucket.options.LruStoreHitMin) *
				float64(timeNow.Sub(item.last).Nanoseconds()) /
				float64(bucket.options.LruStoreHitInterval.Nanoseconds()))
			if item.k < 0 {
				item.k = 0
			}
			item.k += 1
			if item.k >= bucket.options.LruStoreHitMin {
				canAdd = true
			}
			item.last = timeNow
		} else {
			newItem := &Item{
				k:    1,
				last: time.Now(),
			}
			bucket.pending.Add(key, newItem)
		}
	} else {
		canAdd = true
	}

	return nil, false, canAdd
}



func (bucket *BucketLruTwoQueue) Remove(key string) {
	bucket.mu.RLock()
	defer bucket.mu.RUnlock()

	if bucket.pending != nil {
		bucket.pending.Remove(key)
	}

	if bucket.saved == nil {
		bucket.saved.Remove(key)
	}
}
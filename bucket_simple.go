package go_lru_cache

import (
	gcache "github.com/patrickmn/go-cache"
	"time"
)
type BucketSimple struct {
	BucketI
	items *gcache.Cache
}

func NewBucketSimple(opts *CacheOptions) BucketI {
	bucket := &BucketSimple{
		items:   gcache.New(opts.DefaultTTL, opts.DefaultCleanUpInterval),
	}
	return bucket
}

func (bucket *BucketSimple)Add(key string, value interface{}, ttl time.Duration) bool {
	bucket.items.Add(key, value, ttl)
	return false
}

func (bucket *BucketSimple)Get(key string) (interface{}, bool, bool) {
	v, ok := bucket.items.Get(key)
	return v, ok, true
}

func (bucket *BucketSimple)Remove(key string) {
	bucket.items.Delete(key)
}


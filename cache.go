package go_lru_cache

import (
	"time"
)

const DefaultBucketSize = 32
const DefaultSize = 10240

const CacheTypeSimple = "simple"
const CacheTypeLru = "lru"
const CacheTypeArc = "arc"
const CacheTypeTwoQueue = "2q"

type CacheOptions struct {
	CacheType string
	Size int
	BucketSize int

	DefaultTTL time.Duration
	DefaultCleanUpInterval time.Duration
	LruStoreHitMin int
	LruStoreHitInterval time.Duration

	QpsMax int64
	ReservedQpsMin int64
	MissedQpsMax int64
}

type CacheI interface {
    Add(key string, value interface{}, ttl time.Duration) bool
	Get(key string) (data interface{}, found bool, circuitBreaker bool, canAdd bool)
    Remove(key string)
}

type Cache struct {
	CacheI

	options *CacheOptions
	buckets []BucketI
	bucketMast uint32

	Limiter *Limiter
}

func NewCache(options *CacheOptions) (CacheI, error) {
	if options == nil {
		options = &CacheOptions{}
	}

	if options.BucketSize <= 0 {
		options.BucketSize = DefaultBucketSize
	}

	if options.Size <= 0 {
		options.Size = DefaultSize
	}

	if len(options.CacheType) == 0 {
		options.CacheType = CacheTypeSimple
	}

	if options.DefaultTTL == 0 {
		options.DefaultTTL = 10 * 365 * 24 * 3600 * time.Second
	}

	cache := &Cache{
		options:    options,
		bucketMast: uint32(options.BucketSize),
	}

	if cache.options.QpsMax > 0 || cache.options.MissedQpsMax > 0 {
		cache.Limiter =  &Limiter{
			options:       options,
		}
	}

	var buckets []BucketI
	for i:= 0; i < options.BucketSize; i++ {
		bucket, err := NewBucket(cache.options)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}
	cache.buckets = buckets

	return cache, nil
}

func (cache *Cache) Add(key string, value interface{}, ttl time.Duration) bool {
	if ttl == 0 {
		ttl = cache.options.DefaultTTL
	}

	bucketId := int(bucketHash(key)%uint64(cache.bucketMast))
	return cache.buckets[bucketId].Add(key, value, ttl)
}


func (cache *Cache) Get(key string) (interface{}, bool, bool, bool) {
	if cache.options.QpsMax > 0 || cache.options.MissedQpsMax > 0 {
		if cache.Limiter.Acquire() == false {
			return nil, false, true, false
		}
	}

	bucketId := int(bucketHash(key)%uint64(cache.bucketMast))
	value, ok, canAdd := cache.buckets[bucketId].Get(key)
	if ok {
		if cache.options.QpsMax > 0 || cache.options.MissedQpsMax > 0 {
			cache.Limiter.Hit()
		}
	}
	return value, ok, false, canAdd
}

func (cache *Cache) Remove(key string) {
	bucketId := int(bucketHash(key)%uint64(cache.bucketMast))
	cache.buckets[bucketId].Remove(key)
}

const offset64 = 14695981039346656037
const prime64 = 1099511628211

func bucketHash(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return hash
}
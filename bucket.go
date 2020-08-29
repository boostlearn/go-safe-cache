package go_lru_cache

import (
	"errors"
	"time"
)

type BucketI interface {
	Add(key string, value interface{}, ttl time.Duration) bool
	Get(key string) (data interface{}, found bool, canAdd bool)
	Remove(key string)
}

func NewBucket(opts *CacheOptions) (BucketI, error) {
	switch opts.CacheType {
	case CacheTypeSimple:
		return NewBucketSimple(opts), nil
	case CacheTypeLru:
		return NewBucketLru(opts)
	case CacheTypeTwoQueue:
		return NewBucketLruTwoQueue(opts)
	case CacheTypeArc:
		return NewBucketLruArc(opts)
	default:
		return nil, errors.New("not support cache")
	}
}


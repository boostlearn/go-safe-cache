package go_lru_cache

import (
	"log"
	"testing"
)

func TestNewCache(t *testing.T) {
	cache, err := NewCache(&CacheOptions{
		CacheType:              "",
		Size:                   0,
		BucketSize:             0,
		DefaultTTL:             0,
		DefaultCleanUpInterval: 0,
		LruStoreHitMin:         0,
		LruStoreHitInterval:    0,
		QpsMax:                 0,
		ReservedQpsMin:         0,
		MissedQpsMax:           0,
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = cache.Add("aaa", "bbb", -1)
	v, ok, _, _ := cache.Get("aaa")
	log.Println(v, ok)

}

package go_lru_cache

import "time"

type Item struct {
	k int
	last time.Time
	expire time.Time
	data interface{}
}


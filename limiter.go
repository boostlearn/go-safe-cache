package go_lru_cache

import (
	"sync"
	"sync/atomic"
	"time"
)

const HistorySize = 12
const DefaultDecreaseRatio = 0.9

type Limiter struct {
	mu sync.RWMutex
	options *CacheOptions

	qps float64
	missedQps float64
	hitRate float64

	recordIdx int
	recordTs int64
	requestRecorder [HistorySize]int64
	hitRecorder [HistorySize]int64
}

func (limiter *Limiter) Acquire() bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	timeNow := time.Now()
	if atomic.LoadInt64(&limiter.recordTs) != timeNow.Unix() {
		limiter.qps = limiter.qps * DefaultDecreaseRatio +
			float64(limiter.requestRecorder[limiter.recordIdx]) * ( 1 - DefaultDecreaseRatio)
		limiter.missedQps = limiter.missedQps * DefaultDecreaseRatio +
			float64(limiter.requestRecorder[limiter.recordIdx] - limiter.hitRecorder[limiter.recordIdx]) * ( 1 - DefaultDecreaseRatio)
		limiter.hitRate = limiter.hitRate * DefaultDecreaseRatio +
			float64(limiter.hitRecorder[limiter.recordIdx])/float64(1 + limiter.requestRecorder[limiter.recordIdx]) * ( 1 - DefaultDecreaseRatio)
		limiter.recordTs = timeNow.Unix()
		limiter.recordIdx = int(limiter.recordTs % int64(HistorySize))
		limiter.requestRecorder[limiter.recordIdx] = 0
		limiter.hitRecorder[limiter.recordIdx] = 0

		atomic.StoreInt64(&limiter.recordTs, timeNow.Unix())
	}

	limiter.requestRecorder[limiter.recordIdx] +=  1

	if limiter.options.ReservedQpsMin > 0 && limiter.requestRecorder[limiter.recordIdx] < limiter.options.ReservedQpsMin {
		return true
	}

	if limiter.options.QpsMax > 0 && limiter.requestRecorder[limiter.recordIdx] > limiter.options.QpsMax {
		return false
	}

	//if limiter.options.MissedQpsMax > 0 &&
	//	limiter.requestRecorder[limiter.recordIdx] - limiter.hitRecorder[limiter.recordIdx] > limiter.options.MissedQpsMax {
	//	return false
	//}

	return true
}

func (limiter *Limiter) Hit() {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	limiter.hitRecorder[limiter.recordIdx] += 1
}

func (limiter *Limiter) Qps() float64 {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	return limiter.qps
}

func (limiter *Limiter) MissedQps() float64 {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	return limiter.missedQps
}

func (limiter *Limiter) HitRate() float64 {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	return limiter.hitRate
}



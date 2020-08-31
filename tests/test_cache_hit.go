package main

import (
	"fmt"
	go_lru_cache "github.com/boostlearn/go-safe-cache"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

var paretoData []string
var dataFile = "pareto_data/pareto_%v.data"

type CacheMetric struct {
	Alpha int
	CacheType string
	CacheSize int64
	Get int64
	Hit int64
	Add int64
	K int
}

func main() {
	fmt.Println("alpha,cache_size,cache_type,k,get, add,hit")
	for _, alpha := range []int{4} {
		data, err := ioutil.ReadFile(fmt.Sprintf(dataFile, alpha))
		if err != nil {
			log.Fatal("openfile error:", dataFile)
		}

		paretoData = strings.Split(string(data), "\n")

		for _, CacheSize := range []int{100, 500, 1000, 2000} {
			for _, minHit := range []int{1, 2, 4, 8, 16} {
				for _, cacheType := range []string{
					go_lru_cache.CacheTypeLru,
					go_lru_cache.CacheTypeArc,
					go_lru_cache.CacheTypeTwoQueue} {

					cacheSingle, err := go_lru_cache.NewCache(&go_lru_cache.CacheOptions{
						CacheType:  cacheType,
						Size:       CacheSize,
						BucketSize: 5,
					})
					if err != nil {
						log.Fatal(err)
					}

					cacheK, err := go_lru_cache.NewCache(&go_lru_cache.CacheOptions{
						CacheType:           cacheType,
						Size:                CacheSize,
						BucketSize:          5,
						LruStoreHitMin:      minHit,
						LruStoreHitInterval: time.Minute,
					})
					if err != nil {
						log.Fatal(err)
					}

					var metricSingle = &CacheMetric{
						Alpha: alpha,
						CacheType: cacheType,
						CacheSize: int64(CacheSize),
					}
					var metricK = &CacheMetric{CacheType: cacheType,
						Alpha: alpha,
						K:         minHit,
						CacheSize: int64(CacheSize),
					}

					for i := 0; i < 10000000; i++ {
						d := paretoData[rand.Intn(len(paretoData))]
						_, found, canAdd := cacheSingle.Get(d)
						metricSingle.Get += 1
						if found {
							metricSingle.Hit += 1
						}
						if canAdd && found == false {
							cacheSingle.Add(d, "v", 0)
							metricSingle.Add += 1
						}

						_, foundK, canAddK := cacheK.Get(d)
						metricK.Get += 1
						if foundK {
							metricK.Hit += 1
						}
						if canAddK && foundK == false {
							cacheK.Add(d, "v", 0)
							metricK.Add += 1
						}
					}

					for _, m := range []*CacheMetric{metricSingle, metricK} {
						fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n",
							m.Alpha,
							m.CacheSize,
							m.CacheType,
							m.K, m.Get,
							m.Add,
							m.Hit)
					}
				}

			}
		}
	}
}



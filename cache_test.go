package go_lru_cache

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var paretoData []string
var dataFile = "tests/pareto_data/pareto_2.data"
var CacheSize = 1000
var minHit = 8

func init() {
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		log.Fatal("openfile error:", dataFile)
	}

	paretoData =  strings.Split(string(data), "\n")
}


func TestNewCache(t *testing.T) {
	cache, err := NewCache(&CacheOptions{})
	if err != nil {
		log.Fatal(err)
	}

	_ = cache.Add("aaa", "bbb", -1)
	v, ok, _ := cache.Get("aaa")
	log.Println(v, ok)

}

func BenchmarkBucketLru_Single(b *testing.B) {
	cache, err := NewCache(&CacheOptions{
		CacheType:              CacheTypeLru,
		Size:                   CacheSize,
		BucketSize:             5,
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

	b.ReportAllocs()
	b.ResetTimer()

	getCounter := 0
	hitCounter := 0
	addCounter := 0
	for i := 0; i < b.N; i++ {
		d := paretoData[rand.Intn(len(paretoData))]
		_, found,canAdd := cache.Get(d)
		getCounter += 1
		if found {
			hitCounter += 1
		}
		if canAdd && found == false {
			cache.Add(d, "v", 0)
			addCounter += 1
		}
	}
	//fmt.Printf("get:%v, hit:%v, add:%v\n", getCounter, hitCounter, addCounter)
}

func BenchmarkBucketLru_K(b *testing.B) {
	cache, err := NewCache(&CacheOptions{
		CacheType:              CacheTypeLru,
		Size:                   CacheSize,
		BucketSize:             5,
		DefaultTTL:             0,
		DefaultCleanUpInterval: 0,
		LruStoreHitMin:         minHit,
		LruStoreHitInterval:    time.Second,
		QpsMax:                 0,
		ReservedQpsMin:         0,
		MissedQpsMax:           0,
	})
	if err != nil {
		log.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	getCounter := 0
	hitCounter := 0
	addCounter := 0
	for i := 0; i < b.N; i++ {
		d := paretoData[rand.Intn(len(paretoData))]
		_, found, canAdd := cache.Get(d)
		getCounter += 1
		if found {
			hitCounter += 1
		}
		if canAdd && found == false {
			cache.Add(d, "v", 0)
			addCounter += 1
		}
	}
	//fmt.Printf("get:%v, hit:%v, add:%v\n", getCounter, hitCounter, addCounter)
}

func BenchmarkBucket2Q_Single(b *testing.B) {
	cache, err := NewCache(&CacheOptions{
		CacheType:              CacheTypeTwoQueue,
		Size:                   CacheSize,
		BucketSize:             5,
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

	b.ReportAllocs()
	b.ResetTimer()

	getCounter := 0
	hitCounter := 0
	addCounter := 0
	for i := 0; i < b.N; i++ {
		d := paretoData[rand.Intn(len(paretoData))]
		_, found, canAdd := cache.Get(d)
		getCounter += 1
		if found {
			hitCounter += 1
		}
		if canAdd && found == false {
			cache.Add(d, "v", 0)
			addCounter += 1
		}
	}
	//fmt.Printf("get:%v, hit:%v, add:%v\n", getCounter, hitCounter, addCounter)
}

func BenchmarkBucket2Q_K(b *testing.B) {
	cache, err := NewCache(&CacheOptions{
		CacheType:              CacheTypeTwoQueue,
		Size:                   CacheSize,
		BucketSize:             5,
		DefaultTTL:             0,
		DefaultCleanUpInterval: 0,
		LruStoreHitMin:         minHit,
		LruStoreHitInterval:    time.Second,
		QpsMax:                 0,
		ReservedQpsMin:         0,
		MissedQpsMax:           0,
	})
	if err != nil {
		log.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	getCounter := 0
	hitCounter := 0
	addCounter := 0
	for i := 0; i < b.N; i++ {
		d := paretoData[rand.Intn(len(paretoData))]
		_, found, canAdd := cache.Get(d)
		getCounter += 1
		if found {
			hitCounter += 1
		}
		if canAdd && found == false {
			cache.Add(d, "v", 0)
			addCounter += 1
		}
	}
	//fmt.Printf("get:%v, hit:%v, add:%v\n", getCounter, hitCounter, addCounter)
}

func BenchmarkBucketArc_Single(b *testing.B) {
	cache, err := NewCache(&CacheOptions{
		CacheType:              CacheTypeArc,
		Size:                   CacheSize,
		BucketSize:             5,
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

	b.ReportAllocs()
	b.ResetTimer()

	getCounter := 0
	hitCounter := 0
	addCounter := 0
	for i := 0; i < b.N; i++ {
		d := paretoData[rand.Intn(len(paretoData))]
		_, found, _, canAdd := cache.Get(d)
		getCounter += 1
		if found {
			hitCounter += 1
		}
		if canAdd && found == false {
			cache.Add(d, "v", 0)
			addCounter += 1
		}
	}
	//fmt.Printf("get:%v, hit:%v, add:%v\n", getCounter, hitCounter, addCounter)
}


func BenchmarkBucketArc_K(b *testing.B) {
	cache, err := NewCache(&CacheOptions{
		CacheType:              CacheTypeArc,
		Size:                   CacheSize,
		BucketSize:             5,
		DefaultTTL:             0,
		DefaultCleanUpInterval: 0,
		LruStoreHitMin:         minHit,
		LruStoreHitInterval:    time.Second,
		QpsMax:                 0,
		ReservedQpsMin:         0,
		MissedQpsMax:           0,
	})
	if err != nil {
		log.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	getCounter := 0
	hitCounter := 0
	addCounter := 0
	for i := 0; i < b.N; i++ {
		d := paretoData[rand.Intn(len(paretoData))]
		_, found, _, canAdd := cache.Get(d)
		getCounter += 1
		if found {
			hitCounter += 1
		}
		if canAdd && found == false {
			cache.Add(d, "v", 0)
			addCounter += 1
		}
	}
	//fmt.Printf("get:%v, hit:%v, add:%v\n", getCounter, hitCounter, addCounter)
}
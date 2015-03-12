package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/csimplestring/go-concurrent-map/bench/uniuri"
	"github.com/csimplestring/go-concurrent-map/ccmap"
)

var (
	byte1024 = make([]byte, 1024)
	byte2048 = make([]byte, 2048)
)

func init() {
	for i := 0; i < 1024; i++ {
		byte1024[i] = '0'
	}
	for i := 0; i < 2048; i++ {
		byte2048[i] = '0'
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	writerNum := 15
	readerNum := 1
	writeTotalOps := 100000
	readTotalOps := 100000

	cmap := ccmap.NewSimpleMap()

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < writerNum; i++ {
		wg.Add(1)
		go write(cmap, writeTotalOps, &wg)
	}
	for i := 0; i < readerNum; i++ {
		wg.Add(1)
		go read(cmap, readTotalOps, &wg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	wQps := math.Floor(float64(writeTotalOps*2*(writerNum+readerNum)) / elapsed.Seconds())
	fmt.Printf("qps %f\n", wQps)

	//cmap.Stats()
}

func write(c ccmap.ConcurrentMap, ops int, wg *sync.WaitGroup) {
	for i := 0; i < ops; i++ {
		c.Put(ccmap.NewStringKey(randKey(30)), randKey(20))
	}
	wg.Done()
}

func read(c ccmap.ConcurrentMap, ops int, wg *sync.WaitGroup) {
	for i := 0; i < ops; i++ {
		c.Get(ccmap.NewStringKey(randKey(30)))
	}
	wg.Done()
}

// randKey generates a key whose size 1-256
func randKey(size int) string {
	if size < 0 {
		l := rand.Intn(256)
		if l == 0 {
			l = 1
		}
		return uniuri.NewLen(l)
	}
	return uniuri.NewLen(size)
}

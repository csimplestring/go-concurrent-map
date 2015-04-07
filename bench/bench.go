package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/csimplestring/go-concurrent-map/algo/random"
	"github.com/csimplestring/go-concurrent-map/ccmap"
	"github.com/csimplestring/go-concurrent-map/ccmap/key"
	"github.com/csimplestring/go-concurrent-map/ccmap/v1"
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
	runtime.GOMAXPROCS(4)
	h1, _ := v1.NewHashMap(1024)
	h2, _ := v1.NewConcurrentMap(16)

	printStat("locked hash map", suit(15, 100000, 30, 100000, h1))
	printStat("concurrent hash map", suit(15, 100000, 30, 100000, h2))

}

func suit(wNum, wOps, rNum, rOps int, cmap ccmap.Map) []string {
	rand.Seed(time.Now().UTC().UnixNano())

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < wNum; i++ {
		wg.Add(1)
		go write(cmap, wOps, i, &wg)
	}

	for i := 0; i < rNum; i++ {
		wg.Add(1)
		go read(cmap, rOps, i, &wg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	wQps := math.Floor(float64(wNum*wOps) / elapsed.Seconds())
	rQps := math.Floor(float64(rNum*rOps) / elapsed.Seconds())
	qps := math.Floor(float64(wNum*wOps+rNum*rOps) / elapsed.Seconds())

	stat := make([]string, 3)
	stat[0] = fmt.Sprintf("writes qps %f\n", wQps)
	stat[1] = fmt.Sprintf("reads qps %f\n", rQps)
	stat[2] = fmt.Sprintf("average qps %f\n", qps)
	return stat
}

func printStat(suitName string, stat []string) {
	fmt.Printf("--------- %s -----------\n", suitName)
	for _, v := range stat {
		fmt.Println(v)
	}
}

func write(c ccmap.Map, ops, routineIdx int, wg *sync.WaitGroup) {
	for i := 1; i < ops; i++ {
		c.Put(key.NewStringKey(random.NewLen(15)), "string")
	}
	wg.Done()
}

func read(c ccmap.Map, ops, routineIdx int, wg *sync.WaitGroup) {
	for i := 1; i < ops; i++ {
		c.Get(key.NewStringKey(random.NewLen(15)))
	}
	wg.Done()
}

// randKey generates a key whose size 1-256
func randKey() string {
	return strconv.Itoa(rand.Int())
}

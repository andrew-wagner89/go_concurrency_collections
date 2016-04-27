package main

import (
	"./Lists"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var numthreads = 8
var itersperthread = 1024 * 1024

type List interface {
	Insert(key int) bool
	Remove(key int) bool
	Contains(key int) bool
	Init()
	Printlist()
}

func testList(list List, seed int, wg *sync.WaitGroup) {
	fmt.Printf("Testing with thread %d\n", seed)
	rand.Seed((int64)(seed))
	method := rand.Intn(3)
	key := rand.Intn(256)
	for i := 0; i < itersperthread; i++ {
		if method == 0 {
			list.Insert(key)
		} else if method == 1 {
			list.Remove(key)
		} else {
			list.Contains(key)
		}
	}
	wg.Done()
}

func main() {
	list := new(Lists.CGList)
	list.Init()
	var wg sync.WaitGroup
	wg.Add(numthreads)

	start := time.Now()
	for i := 0; i < numthreads; i++ {
		go testList(list, i, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Finished testing %d threads with %d iterations per thread in:\n%s\n", numthreads, itersperthread, elapsed)
}

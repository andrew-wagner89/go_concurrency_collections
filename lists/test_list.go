package main

import 
(	"fmt"
	"math/rand"
	"sync"
	"time"
	
)

var numthreads = 8
var itersperthread = 1024 * 1024

func main() {
	list := NewList()
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

func testOneList() {
	list := NewList()
	list.insert(1)
	list.insert(7)
	list.insert(2)
	list.insert(16)
	fmt.Println(list.insert(2))
	fmt.Println(list.contains(3))
	fmt.Println(list.remove(20))
	fmt.Println(list.contains(1))
	fmt.Println(list.remove(2))
	list.printlist()
}

func testList(list *List, seed int, wg *sync.WaitGroup) {
	fmt.Printf("Testing with thread %d\n", seed)
	rand.Seed((int64)(seed))
	method := rand.Intn(3)
	key := rand.Intn(256)
	for i := 0; i < itersperthread; i++ {
		if method == 0 {
			list.insert(key)
		} else if method == 1 {
			list.remove(key)
		} else {
			list.contains(key)
		}
	}
	wg.Done()
}
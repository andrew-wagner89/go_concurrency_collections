package main

import 
(	"fmt"
	"math/rand"
	"sync"
	"time"
	"bufio"
	"strconv"
	"os"
	"../concurrent_list"
	"../lock_free_list"
	"../coarse_grain_list"
	"../lazy_locking_list"
)

var numthreads = 8
var itersperthread = 1024 * 1024

func main() {

	//take in input to see which list to use, TODO: change to command line input
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter 1 for coarse grain, 2 for lock free and 3 for lazy locking: ")
    text, _ := reader.ReadString('\n')
    listTypeInt := strconv.Atoi(text)

	var list *concurrent_list.Concurrent_List //uses interface here
	var listTypeStr string


	switch listTypeInt {
	case 1:
		list = coarse_grain_list.NewList()
		listTypeStr = "coarse grain"
	case 2:
		list = lock_free_list.NewList()
		listTypeStr = "lock free"
	case 3:
		list = lazy_locking_list.NewList()
		listTypeStr = "lazy locking"
	default:
		fmt.Printf("improper input detected")
		os.Exit(1)

	}

	var wg sync.WaitGroup
	wg.Add(numthreads)

	start := time.Now()
	for i := 0; i < numthreads; i++ {
		go testList(list, i, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Finished testing %s list with %d threads and %d iterations per thread in:\n%s\n", listTypeStr, numthreads, itersperthread, elapsed)
}

func testList(list *concurrent_list.Concurrent_List, seed int, wg *sync.WaitGroup) { //interface used here too
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
package main

import (
	"./Lists"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var numthreads = 1
var itersperthread = 1024 * 64
var maxkeyval = 4096
var numBuckets = 12

//test function for the map, each thread will run this
func testMap(hMap *Lists.HashMap, seed int, wg *sync.WaitGroup) {
	fmt.Printf("Testing with thread %d\n", seed)
	rand.Seed((int64)(seed))
	var method int
	var key int
	var val int
	for i := 0; i < itersperthread; i++ {
		key = rand.Intn(maxkeyval)
		val = rand.Intn(maxkeyval)
		method = rand.Intn(3)

		if method == 0 {
			hMap.Insert(key, val)
		} else if method == 1 {
			hMap.Remove(key)
		} else {
			hMap.Get(key)
		}
	}
	wg.Done()
}

func main() {
	//take in input to see which list to use
	//TODO: change to command line input
	fmt.Print("Enter 1 for coarse grain, 2 for lock free and 3 for lazy locking: ")
	var inputstr string
	_, err := fmt.Scanf("%s", &inputstr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	input, e := strconv.Atoi(inputstr)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	hMap := new(Lists.HashMap)
	hMap.Init(numBuckets, input)

	//fmt.Println("Running tests...")
	//Lists.Runtests(list)
	//fmt.Println("Tests complete\n")

	var wg sync.WaitGroup
	wg.Add(numthreads)

	start := time.Now()
	for i := 0; i < numthreads; i++ {
		go testMap(hMap, i, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Finished testing %d threads with %d iterations per thread in:\n%s\n", numthreads, itersperthread, elapsed)
}

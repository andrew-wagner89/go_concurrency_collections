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

var numthreads = 4
var itersperthread = 1024

func testList(list Lists.List, seed int, wg *sync.WaitGroup) {
	fmt.Printf("Testing with thread %d\n", seed)
	rand.Seed((int64)(seed))
	method := rand.Intn(3)
	key := rand.Intn(256)
	val := rand.Intn(256)
	for i := 0; i < itersperthread; i++ {
		if method == 0 {
			list.Insert(key, val)
		} else if method == 1 {
			list.Remove(key)
		} else {
			list.Get(key)
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
	var list Lists.List

	switch input {
	case 1:
		list = new(Lists.CGList)
	case 2:
		list = new(Lists.LFList)
	//case 3:
	//list = new(Lists.LazyList)
	default:
		fmt.Printf("improper input detected")
		os.Exit(1)
	}
	list.Init()

	fmt.Println("Running tests...")
	Lists.Runtests(list)
	fmt.Println("Tests complete\n")

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

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
var itersperthread = 1024 * 64
var maxkeyval = 4096

func testList(list Lists.List, seed int, wg *sync.WaitGroup) {
	fmt.Printf("Testing with thread %d\n", seed)
	rand.Seed((int64)(seed))
	var key int
	var val int
	var method int
	for i := 0; i < itersperthread; i++ {
		key = rand.Intn(maxkeyval)
		val = rand.Intn(maxkeyval)
		method = rand.Intn(3)
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

func testHash() {
	rand.Seed((int64)(0))
	start := time.Now()
	var key int
	var hash uint64
	for i := 0; i < itersperthread; i++ {
		key = rand.Intn(maxkeyval)
		hash, _ = Lists.GetHash(key)
		_ = hash % numBuckets
		//fmt.Printf("Hash of %d is %d\n", key, hash)
	}
	elapsed := time.Since(start)
	fmt.Printf("Computing %d hashes took %s\n", itersperthread, elapsed)

}

func main() {
	//take in input to see which list to use
	//TODO: change to command line input

	testHash()

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
	case 3:
		list = new(Lists.LazyList)
	default:
		fmt.Printf("improper input detected\n")
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

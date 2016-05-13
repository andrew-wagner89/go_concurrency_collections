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

var numthreads = 8
var itersperthread = 1024 * 64
var maxkeyval = 4096
var numBuckets = 64

func testHash() {
	rand.Seed((int64)(0))
	start := time.Now()
	var key int
	var hash uint64
	for i := 0; i < itersperthread; i++ {
		key = rand.Intn(maxkeyval)
		hash, _ = Lists.GetHash(key)
		_ = hash % uint64(numBuckets)
		//fmt.Printf("Hash of %d is %d\n", key, hash)
	}
	elapsed := time.Since(start)
	fmt.Printf("Computing %d hashes took %s\n", itersperthread, elapsed)
}

//test function for the map, each thread will run this
func testHashMap(hMap *Lists.HashMap, seed int, wg *sync.WaitGroup) {

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
	testHash()
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

	startConc := time.Now()
	for i := 0; i < numthreads; i++ {
		go testHashMap(hMap, i, &wg)
	}
	wg.Wait()
	elapsedConc := time.Since(startConc)

	//test go's map
	goMap := make(map[int]int)

	startSeq := time.Now()
	rand.Seed((int64)(0))
	var method int
	var key int
	var val int
	var trash int
	for i := 0; i < itersperthread*numthreads; i++ {
		key = rand.Intn(maxkeyval)
		val = rand.Intn(maxkeyval)
		method = rand.Intn(3)

		if method == 0 {
			goMap[key] = val
		} else if method == 1 {
			delete(goMap, key)
		} else {
			trash = goMap[key]
		}
	}

	elapsedSeq := time.Since(startSeq)
	fmt.Printf("IGNORE: this output is to satisfy go%d\n\n\n", trash)
	fmt.Printf("Finished testing %d threads with %d iterations per thread:\n", numthreads, itersperthread)
	fmt.Printf("Concurrent Hash map tooK: %s\n", elapsedConc)
	fmt.Printf("Go's sequential hash map took: %s\n", elapsedSeq)

}

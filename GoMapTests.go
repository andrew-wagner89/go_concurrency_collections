package main

import (
	"./Lists"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var numthreads = 16
var itersperthread = 1024 * 64
var maxkeyval = 4096
var numBuckets = 4

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

	var listTypeStr = flag.String("type", "CG", "Type of list")
	flag.Parse()

	var listType = Lists.ParseType(*listTypeStr)

	hMap := new(Lists.HashMap)
	hMap.Init(numBuckets, listType)

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
	for i := 0; i < itersperthread*numthreads; i++ {
		key = rand.Intn(maxkeyval)
		val = rand.Intn(maxkeyval)
		method = rand.Intn(3)

		if method == 0 {
			goMap[key] = val
		} else if method == 1 {
			delete(goMap, key)
		} else {
			_ = goMap[key]
		}
	}

	elapsedSeq := time.Since(startSeq)
	fmt.Printf("Finished testing %d threads with %d iterations per thread:\n",
		numthreads, itersperthread)
	fmt.Printf("Concurrent Hash map (%s) took: %s\n", *listTypeStr, elapsedConc)
	fmt.Printf("Go's sequential hash map took: %s\n", elapsedSeq)

}

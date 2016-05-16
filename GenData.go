package main

import (
	"./Lists"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var numBuckets = 1
var maxkeyval = 4096 * 4
var threadsArr = []int{1, 2, 4, 6, 8, 12, 16, 20, 24, 28, 32}

func testHashMap(hMap *Lists.HashMap, seed int, wg *sync.WaitGroup, iters int) {
	rand.Seed((int64)(seed))
	var method int
	var key int
	var val int
	for i := 0; i < iters; i++ {
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

//Keep buckes constant, loop thru diff values of threads and iters
func genDataBucketsConst() {
	fmt.Println("#threads iterations seconds(CG) seconds(LF) seconds(LL)")
	for _, threads := range threadsArr {
		for iters := 1024; iters <= 1024*128; iters *= 2 {
			fmt.Fprintf(os.Stderr, "%d %d\n", threads, iters)
			cgS := oneTest(threads, iters, Lists.CGListType, numBuckets)
			fmt.Fprintf(os.Stderr, "CG done in %f\n", cgS)
			lfS := oneTest(threads, iters, Lists.LFListType, numBuckets)
			fmt.Fprintf(os.Stderr, "LF done in %f\n", lfS)
			llS := oneTest(threads, iters, Lists.LLListType, numBuckets)
			fmt.Fprintf(os.Stderr, "LL done in %f\n", llS)
			fmt.Printf("%d %d %f %f %f\n", threads, iters, cgS, lfS, llS)
		}
		fmt.Println()
	}
}

//Keep iters constant, loop thru diff values of buckets and threads
func genDataItersConst() {
	iters := 1024 * 64
	fmt.Println("#buckets threads seconds(CG) seconds(LF) seconds(LL)")
	for numBuckets := 1; numBuckets <= 1024; numBuckets *= 2 {
		for _, threads := range threadsArr {
			fmt.Fprintf(os.Stderr, "%d %d\n", numBuckets, threads)
			cgS := oneTest(threads, iters, Lists.CGListType, numBuckets)
			fmt.Fprintf(os.Stderr, "CG done in %f\n", cgS)
			lfS := oneTest(threads, iters, Lists.LFListType, numBuckets)
			fmt.Fprintf(os.Stderr, "LF done in %f\n", lfS)
			llS := oneTest(threads, iters, Lists.LLListType, numBuckets)
			fmt.Fprintf(os.Stderr, "LL done in %f\n", llS)
			fmt.Printf("%d %d %f %f %f\n", numBuckets, threads, cgS, lfS, llS)
		}
		fmt.Println()
	}
}

//Run a single test given the parameters
func oneTest(threads, iters int, listType Lists.ListType, numBuckets int) float64 {
	hMap := new(Lists.HashMap)
	hMap.Init(numBuckets, listType)
	var wg sync.WaitGroup
	wg.Add(threads)
	startConc := time.Now()
	for i := 0; i < threads; i++ {
		go testHashMap(hMap, i, &wg, iters)
	}
	wg.Wait()
	elapsedConc := time.Since(startConc)
	return elapsedConc.Seconds()
}

func main() {
	genDataItersConst()
}

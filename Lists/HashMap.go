package Lists

import (
	"fmt"
	"os"
)

type HashMap struct {
	buckets      []List
	numBuckets   int
	numPerBucket uint64
}

func (hm *HashMap) Init(numBuckets int, listType int) {
	hm.numBuckets = numBuckets
	hm.numPerBucket = MAX_UINT64 / uint64(numBuckets)

	hm.buckets = make([]List, numBuckets)

	for i := 0; i < numBuckets; i++ {
		switch listType {
		case 1:
			hm.buckets[i] = new(CGList)
		case 2:
			hm.buckets[i] = new(LFList)
		case 3:
			hm.buckets[i] = new(LazyList)
		default:
			fmt.Printf("improper hashmap type detected\n")
			os.Exit(1)
		}

		hm.buckets[i].Init()
	}
}

func (hm *HashMap) Get(key interface{}) (interface{}, bool) {
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	bucketId := keyHash / hm.numPerBucket

	return hm.buckets[bucketId].Get(key)
}

func (hm *HashMap) Remove(key interface{}) bool {
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	bucketId := keyHash / hm.numPerBucket

	return hm.buckets[bucketId].Remove(key)
}

func (hm *HashMap) Insert(key interface{}, val interface{}) bool {
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	bucketId := keyHash / hm.numPerBucket

	return hm.buckets[bucketId].Insert(key, val)
}

func (hm *HashMap) PrintMap() {
	for i := 0; i < hm.numBuckets; i++ {
		hm.buckets[i].Printlist()
	}
}

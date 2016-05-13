package Lists

import (
	"fmt"
	"os"
)

type HashMap struct {
	buckets    []List
	numBuckets uint64
}

type ListType int

const (
	CGListType = iota
	LFListType
	LLListType
)

func ParseType(str string) ListType {
	switch str {
	case "CG":
		return CGListType
	case "LF":
		return LFListType
	case "LL":
		return LLListType
	default:
		fmt.Printf("Must supply list type: either CG, LF, or LL\n")
		os.Exit(1)
		return CGListType
	}
}

func (hm *HashMap) Init(numBuckets int, listType ListType) {
	hm.numBuckets = uint64(numBuckets)

	hm.buckets = make([]List, numBuckets)

	for i := 0; i < numBuckets; i++ {
		switch listType {
		case CGListType:
			hm.buckets[i] = new(CGList)
		case LFListType:
			hm.buckets[i] = new(LFList)
		case LLListType:
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

	bucketId := keyHash % hm.numBuckets

	return hm.buckets[bucketId].Get(key)
}

func (hm *HashMap) Remove(key interface{}) bool {
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	bucketId := keyHash % hm.numBuckets

	return hm.buckets[bucketId].Remove(key)
}

func (hm *HashMap) Insert(key interface{}, val interface{}) bool {
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	bucketId := keyHash % hm.numBuckets

	return hm.buckets[bucketId].Insert(key, val)
}

func (hm *HashMap) PrintMap() {
	for i := 0; i < int(hm.numBuckets); i++ {
		hm.buckets[i].Printlist()
	}
}

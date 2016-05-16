package Lists

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"strconv"
)

const MAX_UINT64 uint64 = 1<<64 - 1 //max 64 bit number possible
const MIN_UINT64 uint64 = 0         //min 64 bit number possible

//returns the hash value of any interface
//calculates hash based on the interface's bytes
func getHash(key interface{}) (uint64, error) {
	return GetHash(key)
}
func getHashOld(key interface{}) (uint64, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return 0, err
	}

	h := fnv.New64a()
	h.Write([]byte(buf.Bytes()))

	// Map to between [1,MAX_UINT64-1], so we can have sentinel nodes
	hash := h.Sum64()%(MAX_UINT64-2) + 1
	return hash, nil

}

//Not working, binary.Write does nothing?
func getHash2(key interface{}) (uint64, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, key)

	fmt.Printf("Bytes of %d is %x\n", key, buf.Bytes())

	// Map to between [1,MAX_UINT64-1], so we can have sentinel nodes
	h := fnv.New64a()
	h.Write([]byte(buf.Bytes()))
	hash := h.Sum64()%(MAX_UINT64-2) + 1
	return hash, nil
}

func GetHash(key interface{}) (uint64, error) {
	h := fnv.New64a()

	switch v := key.(type) {
	case *int:
		h.Write([]byte(strconv.Itoa(*v)))
	case int:
		h.Write([]byte(strconv.Itoa(v)))
	case *string:
		h.Write([]byte(*v))
	case string:
		h.Write([]byte(v))
	default:
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(key)
		if err != nil {
			return 0, err
		}
		h.Write([]byte(buf.Bytes()))
	}

	// Map to between [1,MAX_UINT64-1], so we can have sentinel nodes
	hash := h.Sum64()%(MAX_UINT64-2) + 1
	return hash, nil
}

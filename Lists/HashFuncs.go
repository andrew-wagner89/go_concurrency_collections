package Lists

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
)

const MAX_UINT64 uint64 = 1<<64 - 1
const MIN_UINT64 uint64 = 0

func getHash(key interface{}) (uint64, error) {
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

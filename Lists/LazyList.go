package Lists

import (
	"fmt"
	"sync"
)

//Class LazyList
//Constructor: NewLazyList()

type NodeLL struct {
	key    interface{}
	val    interface{}
	next   *NodeLL
	marked bool
	hash uint64
	lock   *sync.Mutex
}

type LazyList struct {
	head *NodeLL
	tail *NodeLL
}

func make_nodeLL(key interface{},val interface{}, next *NodeLL) *NodeLL {
	n := new(NodeLL)
	n.key = key
	n.val = val
	n.next = next
	n.marked = false
	hash, _ := getHash(key)
	n.hash = uint64(hash)
	n.lock = &sync.Mutex{}
	return n
}

func (l *LazyList) Init() {
	l.tail = make_nodeLL(0, nil, nil)
	l.head = make_nodeLL(0, nil, l.tail)
	l.head.hash = MIN_UINT64
	l.tail.hash = MAX_UINT64

}

func (l *LazyList) Printlist() {

	t := l.head
	for t != nil {
		fmt.Printf("%+v: %+v", t.key, t.val)
		t = t.next
	}
}

func validate(pred *NodeLL, curr *NodeLL) bool {
	return !pred.marked && !curr.marked && pred.next == curr
}

//Member funcs for LazyList

func (l *LazyList) Insert(key interface{}, val interface{}) bool {
	var returnval bool

	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	for {
		pred := l.head
		curr := pred.next

		for curr.hash < keyHash {
			pred = curr
			curr = curr.next
		}

		pred.lock.Lock()

		if validate(pred, curr) {
			if curr.hash == keyHash && curr.key == key {
				curr.val = val
				returnval = true
			} else {
				new_node := make_nodeLL(key, val, curr)
				pred.next = new_node
				returnval = true
			}
			pred.lock.Unlock()
			break
		}

		pred.lock.Unlock()
	}
	return returnval
}

func (l *LazyList) Get(key interface{}) (interface{}, bool) {

	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	var curr *NodeLL = l.head

	for curr.hash < keyHash {
		curr = curr.next
	}

	for curr.hash == keyHash {
		if curr.hash == keyHash && curr.key == key {
			if !curr.marked {
				return curr.val, true
			}
		}
		curr = curr.next
	}

	return nil, false
}

func (l *LazyList) Remove(key interface{}) bool {
	returnval := false
	breakInfinite := false

	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	for {
		pred := l.head
		curr := l.head.next

		for curr.hash < keyHash {
			pred = curr
			curr = curr.next
		}

		pred.lock.Lock()
		curr.lock.Lock()


		if curr.hash == keyHash{
			for curr.hash == keyHash {
				if validate(pred, curr) {
					if curr.key == key{
						curr.marked = true
						pred.next = curr.next
						returnval = true
						breakInfinite = true
						break
					} else {
						pred = curr
						curr = curr.next
						continue
					}
				} else {
					continue
				}

			}
		} else {
			returnval = false
			pred.lock.Unlock()
			curr.lock.Unlock()
			break
		}

		pred.lock.Unlock()
		curr.lock.Unlock()

		if breakInfinite {
			break
		}
	}

	return returnval
}

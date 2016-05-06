package Lists

import (
	"fmt"
	"sync"
)

//Class LazyList
//Constructor: NewLazyList()

type Node struct {
	next *Node
	key  interface{}
	val  interface{}
	hash uint64
}

type CGList struct {
	head      *Node
	tail      *Node
	list_lock *sync.Mutex
}

func make_node(key interface{}, val interface{}, next *Node) *Node {
	n := new(Node)
	n.key = key
	n.val = val
	n.next = next
	hash, _ := getHash(key)
	n.hash = uint64(hash)
	return n
}

func (l *CGList) Init() {
	l.tail = make_node(0, nil, nil)
	l.head = make_node(0, nil, l.tail)
	l.head.hash = MIN_UINT64
	l.tail.hash = MAX_UINT64
	l.list_lock = &sync.Mutex{}
}

func (l *CGList) Printlist() {
	l.list_lock.Lock()

	t := l.head
	for t != nil {
		fmt.Printf("%+v: %+v", t.key, t.val)
		t = t.next
	}

	l.list_lock.Unlock()
}

//Member funcs for List

func (l *CGList) Insert(key interface{}, val interface{}) bool {
	var returnval bool

	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	pred := l.head
	curr := pred.next

	for curr.hash < keyHash {
		pred = curr
		curr = curr.next
	}

	l.list_lock.Lock()

	if curr.hash == keyHash && curr.key == key {

		returnval = false
	} else {
		new_node := make_node(key, val, curr)
		pred.next = new_node
		returnval = true
	}

	l.list_lock.Unlock()

	return returnval

}

func (l *CGList) Get(key interface{}) (interface{}, bool) {
	l.list_lock.Lock()

	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	var curr *Node = l.head

	for curr.hash < keyHash {
		curr = curr.next
	}

	for curr.hash == keyHash {
		if curr.hash == keyHash && curr.key == key {
			l.list_lock.Unlock()
			return curr.val, true
		}
		curr = curr.next
	}

	l.list_lock.Unlock()

	return nil, false
}

func (l *CGList) Remove(key interface{}) bool {
	l.list_lock.Lock()

	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	pred := l.head
	curr := l.head.next

	for curr.hash < keyHash {
		pred = curr
		curr = curr.next
	}

	for curr.hash == keyHash {
		if curr.hash == keyHash && curr.key == key {
			pred.next = curr.next
			l.list_lock.Unlock()
			return true
		}
		pred = curr
		curr = curr.next
	}

	l.list_lock.Unlock()
	return false
}

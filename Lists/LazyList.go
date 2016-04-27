package Lists

import (
	"fmt"
	"sync"
)

//Class LazyList
//Constructor: NewLazyList()

type NodeLL struct {
	key    int
	next   *NodeLL
	marked bool
	lock   *sync.Mutex
}

type LazyList struct {
	head *NodeLL
	tail *NodeLL
}

func make_nodeLL(key int, next *NodeLL) *NodeLL {
	n := new(NodeLL)
	n.key = key
	n.next = next
	n.marked = false
	n.lock = &sync.Mutex{}
	return n
}

func (l *LazyList) Init() {
	l.tail = make_nodeLL(2147483647, nil)
	l.head = make_nodeLL(-2147483648, l.tail)
}

func (l *LazyList) Printlist() {
	t := l.head
	for t != nil {
		fmt.Println(t.key)
		t = t.next
	}
}

func validate(pred *NodeLL, curr *NodeLL) bool {
	return !pred.marked && !curr.marked && pred.next == curr
}

//Member funcs for LazyList

func (l *LazyList) Insert(key int) bool {
	var returnval bool
	for {
		pred := l.head
		curr := pred.next

		for curr.key < key {
			pred = curr
			curr = curr.next
		}

		pred.lock.Lock()

		if validate(pred, curr) {
			if curr.key == key {
				returnval = false
			} else {
				new_node := make_nodeLL(key, curr)
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

func (l *LazyList) Contains(key int) bool {
	var curr *NodeLL = l.head
	for curr.key < key {
		curr = curr.next
	}

	return (curr.key == key) && (!curr.marked)
}

func (l *LazyList) Remove(key int) bool {
	var returnval bool
	for {
		pred := l.head
		curr := l.head.next

		for curr.key < key {
			pred = curr
			curr = curr.next
		}

		pred.lock.Lock()
		curr.lock.Lock()

		if validate(pred, curr) {
			if curr.key != key {
				returnval = false
			} else {
				curr.marked = true
				pred.next = curr.next
				returnval = true
			}
			pred.lock.Unlock()
			curr.lock.Unlock()
			break
		}

		pred.lock.Unlock()
		curr.lock.Unlock()
	}
	return returnval
}

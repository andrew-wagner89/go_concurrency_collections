package Lists

import (
	"fmt"
	"sync"
)

//Class LazyList
//Constructor: NewLazyList()

type Node struct {
	key  int
	next *Node
}

type CGList struct {
	head      *Node
	tail      *Node
	list_lock *sync.Mutex
}

func make_node(key int, next *Node) *Node {
	n := new(Node)
	n.key = key
	n.next = next
	return n
}

func (l *CGList) Init() {
	l.tail = make_node(2147483647, nil)
	l.head = make_node(-2147483648, l.tail)
	l.list_lock = &sync.Mutex{}
}

func (l *CGList) Printlist() {
	l.list_lock.Lock()

	t := l.head
	for t != nil {
		fmt.Println(t.key)
		t = t.next
	}

	l.list_lock.Unlock()
}

//Member funcs for List

func (l *CGList) Insert(key int) bool {
	var returnval bool
	pred := l.head
	curr := pred.next

	for curr.key < key {
		pred = curr
		curr = curr.next
	}

	l.list_lock.Lock()

	if curr.key == key {
		returnval = false
	} else {
		new_node := make_node(key, curr)
		pred.next = new_node
		returnval = true
	}

	l.list_lock.Unlock()

	return returnval

}

func (l *CGList) Contains(key int) bool {
	l.list_lock.Lock()

	var curr *Node = l.head

	for curr.key < key {
		curr = curr.next
	}

	l.list_lock.Unlock()

	return curr.key == key
}

func (l *CGList) Remove(key int) bool {
	var returnval bool
	l.list_lock.Lock()

	pred := l.head
	curr := l.head.next

	for curr.key < key {
		pred = curr
		curr = curr.next
	}

	if curr.key != key {
		returnval = false
	} else {
		pred.next = curr.next
		returnval = true
	}

	l.list_lock.Unlock()
	return returnval
}

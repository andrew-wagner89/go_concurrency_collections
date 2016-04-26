package main

import (
	"fmt"
	"sync"
)


type Node struct {
	key    int
	next   *Node
	marked bool
	lock *sync.Mutex
}

type Lazy_Locking_List struct {
	head *Node
	tail *Node
}

func make_node(key int, next *Node) *Node {
	n := new(Node)
	n.key = key
	n.next = next
	n.marked = false
	n.lock = &sync.Mutex{}
	return n
}

func NewList() *Lazy_Locking_List {
	l := new(Lazy_Locking_List)
	l.tail = make_node(2147483647, nil)
	l.head = make_node(-2147483648, l.tail)
	return l
}

func (l *Lazy_Locking_List) printlist() {
	t := l.head
	for t != nil {
		fmt.Println(t.key)
		t = t.next
	}
}

func validate(pred *Node, curr *Node) bool {
	return !pred.marked && !curr.marked && pred.next == curr
}

//Member funcs for List


func (l *Lazy_Locking_List) insert(key int) bool {
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
				new_node := make_node(key, curr)
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

func (l *Lazy_Locking_List) contains(key int) bool {
	var curr *Node = l.head
	for curr.key < key {
		curr = curr.next
	}

	return (curr.key == key) && (!curr.marked)
}

func (l *Lazy_Locking_List) remove(key int) bool {
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

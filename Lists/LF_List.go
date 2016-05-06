package Lists

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

//Class LFList
//Constructor: NewLFList()
/*
Harris lock-free list
From :
A Pragmatic Implementation of Non-Blocking Linked-Lists
by Timothy L. Harris (2001)
*/

type Mark int32

const (
	MARKED Mark = 1 + iota
	UNMARKED
)

type NodeLF struct {
	next   *NodeLF
	key    interface{}
	val    interface{}
	hash   uint64
	marked Mark
}

type LFList struct {
	head *NodeLF
	tail *NodeLF
}

func make_nodeLF(key interface{}, val interface{}, next *NodeLF) *NodeLF {
	n := new(NodeLF)
	n.key = key
	n.val = val
	n.next = next
	hash, _ := getHash(key)
	n.hash = uint64(hash)
	n.marked = UNMARKED
	return n
}

func (l *LFList) Init() {
	l.tail = make_nodeLF(0, nil, nil)
	l.head = make_nodeLF(0, nil, l.tail)
	l.head.hash = MIN_UINT64
	l.tail.hash = MAX_UINT64
}

func (l *LFList) Printlist() {
	t := l.head
	for t != nil {
		fmt.Printf("%+v: %+v", t.key, t.val)
		t = t.next
	}
}

//Member funcs for LFList
func (l *LFList) search(key interface{}, left_node **NodeLF) *NodeLF {
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	var left_node_next *NodeLF
	var right_node *NodeLF

search_again:
	for {
		t := l.head
		t_next := l.head.next

		/* 1: Find left_node and right_node */
	inner:
		for ok := true; ok; ok = (t_next.marked == MARKED || (t.hash <= keyHash)) {
			if t_next.marked != MARKED { //Not marked for deletion
				(*left_node) = t
				left_node_next = t_next
			}
			t = t_next
			if t == l.tail {
				break inner
			}
			t_next = t.next

			//Check key equality
			if t.key == key {
				break inner
			}
		}
		right_node = t

		/* 2: Check nodes are adjacent */
		if left_node_next == right_node {
			if (right_node != l.tail) && (right_node.next.marked == MARKED) {
				goto search_again //Marked for deletion, try again
			} else {
				return right_node //Success
			}
		}

		/* 3: Remove one or more marked nodes */
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&(*left_node).next)),
			unsafe.Pointer(left_node_next),
			unsafe.Pointer(right_node)) {
			if (right_node != l.tail) && (right_node.next.marked == MARKED) {
				goto search_again //Should delete right node
			} else {
				return right_node //Sucess
			}
		}
	}
}

func (l *LFList) Insert(key interface{}, val interface{}) bool {
	new_node := make_nodeLF(key, val, nil)
	var right_node *NodeLF
	var left_node *NodeLF
	for {
		right_node = l.search(key, &left_node)
		if (right_node != l.tail) && (right_node.key == key) {
			return false //Already in list
		}
		new_node.next = right_node
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&left_node.next)),
			unsafe.Pointer(right_node),
			unsafe.Pointer(new_node)) {
			return true
		}
	}
}

func (l *LFList) Get(key interface{}) (interface{}, bool) {
	var right_node *NodeLF
	var left_node *NodeLF
	right_node = l.search(key, &left_node)
	if (right_node == l.tail) || (right_node.key != key) {
		return right_node.val, false
	} else {
		return right_node.val, true
	}
}

func (l *LFList) Remove(key interface{}) bool {
	var right_node *NodeLF
	var right_node_next *NodeLF
	var left_node *NodeLF
	for {
		right_node = l.search(key, &left_node)
		if (right_node == l.tail) || (right_node.key != key) {
			return false //Not in the list
		}
		right_node_next = right_node.next
		if right_node.marked == UNMARKED {
			//Set marked
			if atomic.CompareAndSwapInt32(
				(*int32)(&right_node.marked),
				(int32)(UNMARKED),
				(int32)(MARKED)) {
				break //Succesful mark to delete
			}
		}
	}
	//Try to get rid of right_node
	if atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&left_node.next)),
		unsafe.Pointer(right_node),
		unsafe.Pointer(right_node_next)) {
		//Find marked nodes and delete them
		right_node = l.search(right_node.key, &left_node)
	}
	return true
}

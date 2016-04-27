package Lists

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

//Class LFList
//Constructor: NewLFList()

type Mark int32

const (
	MARKED Mark = 1 + iota
	UNMARKED
)

type NodeLF struct {
	key    int
	next   *NodeLF
	marked Mark
}

type LFList struct {
	head *NodeLF
	tail *NodeLF
}

func make_nodeLF(key int, next *NodeLF) *NodeLF {
	n := new(NodeLF)
	n.key = key
	n.next = next
	n.marked = UNMARKED
	return n
}

func (l *LFList) Init() {
	l.tail = make_nodeLF(2147483647, nil)
	l.head = make_nodeLF(-2147483648, l.tail)
}

func (l *LFList) Printlist() {
	t := l.head
	for t != nil {
		fmt.Println(t.key)
		t = t.next
	}
}

//Member funcs for LFList
func (l *LFList) search(key int, left_node **NodeLF) *NodeLF {
	var left_node_next *NodeLF
	var right_node *NodeLF

search_again:
	for {
		t := l.head
		t_next := l.head.next

		/* 1: Find left_node and right_node */
	inner:
		for ok := true; ok; ok = (t_next.marked == MARKED || (t.key < key)) {
			if t_next.marked != MARKED { //Not marked for deletion
				(*left_node) = t
				left_node_next = t_next
			}
			t = t_next
			if t == l.tail {
				break inner
			}
			t_next = t.next
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

func (l *LFList) Insert(key int) bool {
	new_node := make_nodeLF(key, nil)
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

func (l *LFList) Contains(key int) bool {
	var right_node *NodeLF
	var left_node *NodeLF
	right_node = l.search(key, &left_node)
	if (right_node == l.tail) || (right_node.key != key) {
		return false
	} else {
		return true
	}
}

func (l *LFList) Remove(key int) bool {
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

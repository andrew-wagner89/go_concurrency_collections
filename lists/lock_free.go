package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Mark int32

const (
	MARKED Mark = 1 + iota
	UNMARKED
)

type Node struct {
	key    int
	next   *Node
	marked Mark
}

type Lock_Free_List struct {
	head *Node
	tail *Node
}

func make_node(key int, next *Node) *Node {
	n := new(Node)
	n.key = key
	n.next = next
	n.marked = UNMARKED
	return n
}

func NewList() *Lock_Free_List {
	l := new(Lock_Free_List)
	l.tail = make_node(2147483647, nil)
	l.head = make_node(-2147483648, l.tail)
	return l
}

func (l *Lock_Free_List) printlist() {
	t := l.head
	for t != nil {
		fmt.Println(t.key)
		t = t.next
	}
}

//Member funcs for List
func (l *Lock_Free_List) search(key int, left_node **Node) *Node {
	var left_node_next *Node
	var right_node *Node

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

func (l *Lock_Free_List) insert(key int) bool {
	new_node := make_node(key, nil)
	var right_node *Node
	var left_node *Node
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

func (l *Lock_Free_List) contains(key int) bool {
	var right_node *Node
	var left_node *Node
	right_node = l.search(key, &left_node)
	if (right_node == l.tail) || (right_node.key != key) {
		return false
	} else {
		return true
	}
}

func (l *Lock_Free_List) remove(key int) bool {
	var right_node *Node
	var right_node_next *Node
	var left_node *Node
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

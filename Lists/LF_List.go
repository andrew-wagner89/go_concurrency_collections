package Lists

import (
	"container/list"
	"fmt"
	"sync/atomic"
	"unsafe"
)

//Class LFList

func isMarked(ptr *NodeLF) bool {
	addr := (uintptr)(unsafe.Pointer(ptr))
	ret := (addr & 1) == 1
	return ret
}
func setMarked(ptr *NodeLF) *NodeLF {
	addr := (uintptr)(unsafe.Pointer(ptr))
	newaddr := unsafe.Pointer(addr | 1)
	return (*NodeLF)(newaddr)
}
func setUnmarked(ptr *NodeLF) *NodeLF {
	addr := (uintptr)(unsafe.Pointer(ptr))
	//If on 32 bit machine, change this to be a 32 bit value
	newaddr := unsafe.Pointer(addr & (uintptr)(0xFFFFFFFFFFFFFFFE))
	return (*NodeLF)(newaddr)
}
func TestMarks() {
	node := new(NodeLF)
	fmt.Printf("Start: %p\n", node)
	fmt.Printf("Marked? : %t\n", isMarked(node))
	fmt.Printf("unmarked(node) : %p\n", setUnmarked(node))
	node = setMarked(node)
	fmt.Printf("Marked node: %p\n", node)
	fmt.Printf("Marked? : %t\n", isMarked(node))
	node = setUnmarked(node)
	fmt.Printf("Unmarked node: %p\n", node)
	fmt.Printf("Marked? : %t\n", isMarked(node))
}

type NodeLF struct {
	next *NodeLF
	key  interface{}
	val  interface{}
	hash uint64
}

/*
Harris lock-free list
From :
A Pragmatic Implementation of Non-Blocking Linked-Lists
by Timothy L. Harris (2001)
*/
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
	n.hash = hash
	return n
}

func (l *LFList) Init() {
	l.tail = make_nodeLF(0, nil, nil)
	l.head = make_nodeLF(0, nil, l.tail)
	l.head.hash = MIN_UINT64
	l.tail.hash = MAX_UINT64
}

func (l *LFList) KeysAndValues() (*list.List, *list.List) {
	keys := list.New()
	values := list.New()

	t := l.head.next
	for t != l.tail {
		keys.PushBack(t.key)
		values.PushBack(t.val)
		t = t.next
	}

	return keys, values

}

func (l *LFList) Printlist() {
	t := l.head
	for t != nil {
		fmt.Printf("%+v (%d): %+v\n", t.key, t.hash, t.val)
		t = t.next
	}
}

//Helper func for LFList
//Returns node either == to key, or just > than (hash)
//Also sets left_node to be just to the left of returned
func (l *LFList) search(key interface{}, keyHash uint64, left_node **NodeLF) *NodeLF {
	var left_node_next *NodeLF
	var right_node *NodeLF

search_again:
	for {
		t := l.head
		t_next := l.head.next

		/* 1: Find left_node and right_node */
	inner:
		//for ok := true; ok; ok = (isMarked(t_next) || (t.hash <= keyHash)) {
		for {
			if !isMarked(t_next) { //Not marked for deletion
				(*left_node) = t
				left_node_next = t_next
			}
			t = setUnmarked(t_next)
			if t == l.tail {
				break inner
			}
			t_next = t.next

			//For loop condition
			if isMarked(t_next) || (t.hash <= keyHash) {
				//Check key equality
				if t.hash == keyHash && t.key == key && !isMarked(t_next) {
					break inner
				}
			} else {
				break
			}
		}
		right_node = t

		/* 2: Check nodes are adjacent */
		if left_node_next == right_node {
			if (right_node != l.tail) && isMarked(right_node.next) {
				//fmt.Println("Marked node")
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
			if (right_node != l.tail) && isMarked(right_node.next) {
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
		right_node = l.search(key, new_node.hash, &left_node)
		if (right_node != l.tail) && (right_node.key == key) { //Update val
			right_node.val = val
			return false
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
	hash, _ := getHash(key)
	right_node = l.search(key, hash, &left_node)
	if (right_node == l.tail) || (right_node.key != key) {
		return nil, false
	} else {
		return right_node.val, true
	}
}

func (l *LFList) Remove(key interface{}) bool {
	var right_node *NodeLF
	var right_node_next *NodeLF
	var left_node *NodeLF
	hash, _ := getHash(key)
	for {
		right_node = l.search(key, hash, &left_node)
		if (right_node == l.tail) || (right_node.key != key) {
			return false //Not in the list
		}
		right_node_next = right_node.next
		if !isMarked(right_node_next) { //If unmarked
			//Set marked
			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&right_node.next)),
				(unsafe.Pointer)(right_node_next),
				(unsafe.Pointer)(setMarked(right_node_next))) {
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
		right_node = l.search(right_node.key, hash, &left_node)
	}
	return true
}

func (l *LFList) TestCollision() {
	new_node2 := make_nodeLF(2, 4, l.tail)
	new_node1 := make_nodeLF(2, 3, new_node2)
	l.head.next = new_node1
	new_node1.key = 3 //Change [Either] key w/o changing hash
	l.Printlist()
	val, _ := l.Get(2)
	fmt.Println(val)
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var numthreads = 8
var itersperthread = 1024 * 1024

func main() {
	list := NewList()
	var wg sync.WaitGroup
	wg.Add(numthreads)

	start := time.Now()
	for i := 0; i < numthreads; i++ {
		go testList(list, i, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Finished testing %d threads with %d iterations per thread in:\n%s\n", numthreads, itersperthread, elapsed)
}

func testOneList() {
	list := NewList()
	list.insert(1)
	list.insert(7)
	list.insert(2)
	list.insert(16)
	fmt.Println(list.insert(2))
	fmt.Println(list.contains(3))
	fmt.Println(list.remove(20))
	fmt.Println(list.contains(1))
	fmt.Println(list.remove(2))
	list.printlist()
}

func testList(list *List, seed int, wg *sync.WaitGroup) {
	fmt.Printf("Testing with thread %d\n", seed)
	rand.Seed((int64)(seed))
	method := rand.Intn(3)
	key := rand.Intn(256)
	for i := 0; i < itersperthread; i++ {
		if method == 0 {
			list.insert(key)
		} else if method == 1 {
			list.remove(key)
		} else {
			list.contains(key)
		}
	}
	wg.Done()
}



type Node struct {
	key    int
	next   *Node
	marked bool
	lock *sync.Mutex
}

type List struct {
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

func NewList() *List {
	l := new(List)
	l.tail = make_node(2147483647, nil)
	l.head = make_node(-2147483648, l.tail)
	return l
}

func (l *List) printlist() {
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


func (l *List) insert(key int) bool {
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

func (l *List) contains(key int) bool {
	var curr *Node = l.head
	for curr.key < key {
		curr = curr.next
	}

	return (curr.key == key) && (!curr.marked)
}

func (l *List) remove(key int) bool {
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

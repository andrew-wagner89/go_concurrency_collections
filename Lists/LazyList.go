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
	marked bool //marks whether the node is removed from the list
	hash   uint64
	lock   *sync.Mutex //locks individual nodes
}

type LazyList struct {
	head *NodeLL
	tail *NodeLL
}

//node constructor
func make_nodeLL(key interface{}, val interface{}, next *NodeLL) *NodeLL {
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

//initializes the list
func (l *LazyList) Init() {
	//sentinel key and val don't matter, just the hash value
	l.tail = make_nodeLL(0, nil, nil)
	l.head = make_nodeLL(0, nil, l.tail)
	l.head.hash = MIN_UINT64 //head has minimum hash
	l.tail.hash = MAX_UINT64 //tail has maximum hash

}

//prints out the list
func (l *LazyList) Printlist() {

	t := l.head
	for t != nil {
		fmt.Printf("%+v: %+v", t.key, t.val)
		t = t.next
	}
}

//validates whether the two nodes are still in the list and curr is the node after curr
func validate(pred *NodeLL, curr *NodeLL) bool {
	return !pred.marked && !curr.marked && pred.next == curr
}

//Member funcs for LazyList

//inserts a key and val into a list, if the key is already in the list, the val is updated
func (l *LazyList) Insert(key interface{}, val interface{}) bool {

	//hash the key for the search process
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

		//lock the pred and curr since pred's next mey be changed, and currs val may be updated
		pred.lock.Lock()
		curr.lock.Lock()

		if validate(pred, curr) { //make sure pred and curr are still valid
			if curr.hash == keyHash {
				updated := false
				//handle hash collisions by checking all nodes with the search hash until the correct key is found
				for curr.hash == keyHash {
					if curr.hash == keyHash && curr.key == key { //if the key is already in the list, update the val
						curr.val = val
						updated = true
						break
					}
					//if the current node is not the correct node, look at the next node
					pred.lock.Unlock()
					pred = curr
					curr = curr.next
					curr.lock.Lock()
					validate(pred, curr) //revalidate the nodes
				}
				if !updated { //if a val has not been updated, add a new node to the list
					new_node := make_nodeLL(key, val, curr)
					pred.next = new_node
				}
			} else { //if the current hash is not equal to the hash of the key, make a new node
				new_node := make_nodeLL(key, val, curr)
				pred.next = new_node
			}

			pred.lock.Unlock()
			curr.lock.Unlock()
			break
		}

		//if the two nodes are not valid, restart the insert process
		curr.lock.Unlock()
		pred.lock.Unlock()
	}
	return true
}

//get the val that corresponds to the key
//returns two values in the form (val, keyFound)
//keyFound is a boolean specifying if the key was found in the list
func (l *LazyList) Get(key interface{}) (interface{}, bool) {

	//hash the key for the search process
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	//starting from the head, iterate through the list, looking for the hash
	var curr *NodeLL = l.head

	for curr.hash < keyHash {
		curr = curr.next
	}

	//handle collisions with a while loop going through all nodes with the correct hash
	for curr.hash == keyHash {
		if curr.hash == keyHash && curr.key == key && !curr.marked { //if the current node has the right key return it
			return curr.val, true
		}
		curr = curr.next
	}

	return nil, false //if the correct key wasn't found return false
}

//removes the key from the list, returns whether the key was in the list or not
func (l *LazyList) Remove(key interface{}) bool {

	//hash the value for searching purposes
	var keyHash uint64
	hash32, _ := getHash(key)
	keyHash = uint64(hash32)

	for {
		//iterate through the list searching for the hash value
		pred := l.head
		curr := l.head.next

		for curr.hash < keyHash {
			pred = curr
			curr = curr.next
		}

		pred.lock.Lock()
		curr.lock.Lock()

		restart := false //restart holds whether an invalid node pair was found
		if curr.hash == keyHash {
			for curr.hash == keyHash { //loop through all nodes with the correct hash
				if validate(pred, curr) {
					if curr.key == key { //if the key is in the list remove it and return true
						curr.marked = true
						pred.next = curr.next
						pred.lock.Unlock()
						curr.lock.Unlock()
						return true
					} else { // if this node doesn't have the key, go to the next one
						pred.lock.Unlock()
						pred = curr
						curr = curr.next
						curr.lock.Lock()
					}
				} else { //validate failed, stop the remove process
					restart = true
					break
				}
			}
			pred.lock.Unlock()
			curr.lock.Unlock()

			//if a validate didn't work, restart the remove process
			if restart {
				continue
			} else { //otherwise the hash is in the list, but the key isn't, return false
				return false
			}
		} else { //if the hash value wasn't in the list return false
			pred.lock.Unlock()
			curr.lock.Unlock()
			return false
		}
	}
}

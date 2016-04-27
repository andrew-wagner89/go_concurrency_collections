package concurrent_list

type Concurrent_List interface {
	printlist()
	insert(int) bool
	contains(int) bool
	remove(int) bool
}
package Lists

type List interface {
	Insert(key int) bool
	Remove(key int) bool
	Contains(key int) bool
	Init()
	Printlist()
}
package Lists

type List interface {
	Insert(key interface{},val interface{}) bool
	Remove(key interface{}) bool
	Contains(key interface{}) (interface{}, bool)
	Init()
	Printlist()
}
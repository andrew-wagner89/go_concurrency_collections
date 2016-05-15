package Lists

import(
	"sync"
	"fmt"
)


type GoMapRW struct {
	goMap map[interface{}]interface{}
	RWLock *sync.RWMutex
	
}

func (l *GoMapRW) Init() {
	l.goMap = make(map[interface{}]interface{})
	l.RWLock = &sync.RWMutex{}
}

func (l *GoMapRW) Printlist() {
	l.RWLock.RLock()
	for k, v := range l.goMap {
		fmt.Printf("%+v: %+v", k, v)
	}
	l.RWLock.RUnlock()
}

//Member funcs for List

func (l *GoMapRW) Insert(key interface{}, val interface{}) bool {
	l.RWLock.Lock()
	l.goMap[key] = val
	l.RWLock.Unlock()
	return true
}

func (l *GoMapRW) Get(key interface{}) (interface{}, bool) {
	l.RWLock.RLock()
	val := l.goMap[key]
	l.RWLock.RUnlock()
	if val == nil {
		return nil, false
	} else {
		return val, true
	}
	
}

func (l *GoMapRW) Remove(key interface{}) bool {
	l.RWLock.RLock()
	val := l.goMap[key]
	l.RWLock.RUnlock()
	if val == nil {
		return false
	} else {
		l.RWLock.Lock()
		delete(l.goMap, key)
		l.RWLock.Unlock()
		return true
	}
}

func (l *GoMapRW) TestCollision() {
	return
}
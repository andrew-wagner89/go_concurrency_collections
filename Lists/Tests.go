package Lists

import (
	"fmt"
)

func Runtests(list List) {
	list.Insert("Hello", 5)
	key, err := list.Get("Hello")
	if err != false && key != 5 {
		fmt.Println("Insert and Get failed!")
	}

}

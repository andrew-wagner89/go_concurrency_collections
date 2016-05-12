package Lists

import (
	"fmt"
	"os"
)

func assert(test bool, message string) {
	if test == false {
		fmt.Println(message)
		os.Exit(1)
	}
}

func Runtests(list List) {
	assert(list.Insert("Hello", 5), "Insert incorrectly returned false")
	key, found := list.Get("Hello")
	assert(found == true, "Get returned false, should have retruned true")
	assert(key == 5, "Get retruned incorrect key")
	//assert(list.Insert("Hello", 3) == false, "Insert incorrectly returned true")
	assert(list.Remove("Garbage") == false, "Remove incorrectly returned true")
	_, found = list.Get("Garbage")
	assert(found == false, "Get incorrectly returned true")
	assert(list.Remove("Hello") == true, "Remove incorrectly returned false")

	//list.TestCollision()

}

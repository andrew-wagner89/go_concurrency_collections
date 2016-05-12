package Lists

type List interface {
	/* Inserts the key,val pair into the list
	Returns true when key was not in the list
	Returns false if key is already in the list
		Still updates the val in this case
	*/
	Insert(key interface{}, val interface{}) bool

	/* Removes the key,val pair indexed by key
	Returns true on succesful remove
	Returns false if key was not found
	*/
	Remove(key interface{}) bool

	/* Gets the val for the corresponding key
	Returns val,true if key,val is found
	Returns val,false if key,val is NOT found
	*/
	Get(key interface{}) (interface{}, bool)

	Init()
	Printlist()

	//TestCollision()
}

package btree


type BTree interface {
	Put(key int)
	find(key int)*Node
	Print()[]int
}
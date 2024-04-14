package main

import (
	"log"

	"github.com/Shyamnatesan/btree/btree"
)




func main(){
	btree := btree.NewTree()

	for i := 10; i < 100000;  {
		btree.Put(i)
		i += 10
	}

	inorder := btree.Print()
	log.Println("inorder traversal of the whole btree =>>>> ")
	log.Println(inorder)
	log.Println("===================================>")
	// parentMap := map[int][]int{}
	// for i := 10; i < 351; {
	// 	result := btree.Find(i)
	// 	log.Println("result: ")
	// 	log.Println("result node keys", result.keys)
	// 	log.Println("result node isLeaf", result.isLeaf)
	// 	log.Println("result number of keys", result.numKeys)
	// 	log.Println()
	// 	if result.parent == nil {
	// 		parentMap[i] = [M]int{}
	// 	}else{
	// 		parentMap[i] = result.parent.keys
	// 	}
	// 	i+= 10
	// }
	// log.Println("child to parent map => ",parentMap)
	
}
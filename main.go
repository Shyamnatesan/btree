package main

import (
	"log"

	"github.com/Shyamnatesan/btree/btree"
)




func main(){
	btree := btree.NewTree()

	for i := 10; i <= 100000;  {
		btree.Put(i)
		i += 10
	}

	inorder := btree.Print()
	log.Println("inorder traversal of the whole btree =>>>> ")
	log.Println(inorder)
	log.Println("===================================>")
	
}
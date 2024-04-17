package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/Shyamnatesan/btree/btree"
)

func generateUniqueRandomNumbers(n, min, max int) []int {
	if max - min < n {
		fmt.Println("Error: Cannot generate unique numbers with the given range and count.")
		return nil
	}
	nums := make(map[int]bool)
	result := make([]int, 0, n)
	for len(nums) < n {
		num := rand.Intn(max-min) + min
		if !nums[num] {
			nums[num] = true
			result = append(result, num)
		}
	}
	return result
}


func main(){
	btree := btree.NewTree()
	// INSERTION
	insertInput := generateUniqueRandomNumbers(10000, 1, 10001)
	log.Println("insertInput => ", insertInput)
	for _, in := range insertInput{
		btree.Put(in)
	}

	inorder := btree.Print()
	log.Println("inorder traversal of the whole btree ==============================================>>>>>>>>> ")
	log.Println(inorder)
	log.Println("===========================================================================================> ")

	// DELETION
	deleteInput := generateUniqueRandomNumbers(9900, 101, 10001)
	for _, in := range deleteInput {
		btree.Del(in)
	}

	log.Println("btree after deletion")
	inorder = btree.Print()
	log.Println("inorder traversal of the whole btree ==============================================>>>>>>>>> ")
	log.Println(inorder)
	log.Println("===========================================================================================> ")

}
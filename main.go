package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Shyamnatesan/btree/btree"
)

func generateUniqueRandomNumbers(n, min, max int) []int {
	if max-min < n {
		fmt.Println("Error: Cannot generate unique numbers with the given range and count.")
		return nil
	}

	rand.Seed(time.Now().UnixNano())
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
	input := generateUniqueRandomNumbers(10000, 0, 10000)
	for _, in := range input{
		btree.Put(in)
	}

	// for i := 100000; i >= 10;  {
	// 	btree.Put(i)
	// 	i -= 10
	// }

	inorder := btree.Print()
	log.Println("inorder traversal of the whole btree =>>>> ")
	log.Println(inorder)
	log.Println("===================================>")
	
}
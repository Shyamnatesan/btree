package btree

import (
	"log"
)


const (
	M = 5 // ORDER OF A TREE
	MAX_NUM_OF_KEYS = M - 1 // MAXIMUM NUMBER OF KEYS ALLOWED IN A NODE
	ceilOfM = M / 2
	MIN_NUM_OF_KEYS = ceilOfM - 1 // MINIMUM NUMBER OF KEYS ALLOWED IN A NODE EXCEPT ROOT
)

type Node struct {
	numKeys int
	keys [M]int
	children [M + 1]*Node
	parent *Node
	isLeaf bool
}

func (node *Node) searchNode(key int, pos *int) *Node {
	// if found return the node
	// else return nil
	log.Println("searching for key ", key, " in the node ", node.keys)
	for ; *pos < node.numKeys; {
		if key > node.keys[*pos] {
			*pos++;
		}else if key ==  node.keys[*pos]{
			log.Println("found a match at position ", *pos, "in node ", node.keys)
			return node;
		}else{
			break
		}
	}
	return nil
}

func (node *Node) search(key int) *Node {
	if node == nil {
		return nil
	}
	pos := 0

	found := node.searchNode(key, &pos)
	if found != nil {
		return found
	}

	return node.children[pos].search(key)
}

func (node *Node) insertIntoNode(key int) int {
	// find the position "pos" => (which index) where to insert
	log.Println("inserting key", key, " in node", node.keys)
	pos := 0
	for ; pos < node.numKeys;  {
		if key > node.keys[pos] {
			pos++;
		}else if key == node.keys[pos] {
			log.Println("key already exists. No duplicate allowed")
			return -1
		}else{
			break
		}
	}


	// shift the elements to the right based on pos
	for i := node.numKeys - 1; i >= pos; i-- {
		node.keys[i + 1] = node.keys[i]
	}
	// insert
	node.keys[pos] = key
	node.numKeys++
	log.Println("successfully inserted key = ", key, "at position ", pos)
	return pos
}

func (node *Node) maxKeyThresholdReached(key int, rightChild *Node, tree *Tree, pos *int) {
	// this is a recursive function

	// rightChild is initially nil
	// when recursively called for parents, we pass the rightnode


	// if rightnode is not nil, we have to add it to the node.children
	// if rightnode is nil, we do nothing

	// initially, check if the node.numKeys != maxnumberofkeys, 
	// if so return from this function
	if node.numKeys != MAX_NUM_OF_KEYS {
		index := node.insertIntoNode(key)
		log.Println("rightChild passed is", rightChild)
		if rightChild != nil {
            for i := node.numKeys; i > index + 1; i-- {
                node.children[i] = node.children[i - 1]
            }
            node.children[index + 1] = rightChild
            rightChild.parent = node
        }
		return
	}else{
		// else, insert the key in the keys array in a sorted manner
		index := node.insertIntoNode(key)
		if rightChild != nil {
            for i := node.numKeys; i > index + 1; i-- {
                node.children[i] = node.children[i - 1]
            }
            node.children[index + 1] = rightChild
            rightChild.parent = node
        }
		// split it. now you'll have leftnode, rightnode, and median(to pass to the parent node)
		rightNode, median := node.splitNode()
		log.Println("rightnode", rightNode.keys)
		log.Println("median", median)
		if node.parent == nil {
			// log.Println("no parent")
			// if this current node's parent is nil,
			// the currentnode does not have a parent, so we create a new node
			parentNode := NewNode()
			// set the isLeaf to false, since parentNode has leftnode and rightnode as children
			parentNode.isLeaf = false
			// then create a new node and update its children(add leftnode and righnode as children)
			parentNode.insertIntoNode(median)
			// add the leftnode(node) and rightnode as children to this parent node
			parentNode.children[0] = node
			parentNode.children[1] = rightNode
			// also update the parent field in the leftnode(node) and the rightnode
			node.parent = parentNode
			rightNode.parent = parentNode

			// since we create a new node and add a key,
			// we have to make this parentNode as root, and then return.
			// or somehow make this parent node a root
			tree.root = parentNode
			return;
		}else{
			// else, if current node's parent is not nil,
			// then after splitting, we'll have a median right,
			// call this function recursively for the current node's parent and pass the median as key
			log.Println("current node", node.keys)
			log.Println("current node's parent", node.parent.keys)
			node.parent.maxKeyThresholdReached(median, rightNode, tree, pos)

		}
	}
}

func (node *Node) splitNode() (*Node, int) {
	// log.Println("splitting the node ", node.keys)
	medianIndex := node.numKeys / 2
	median := node.keys[medianIndex]

	// create a new node called rightnode
	rightNode := NewNode()
	// rightnode.isLeaf = node.isLeaf (because only if node has children,
	//  they will be copied to the rightnode)
	rightNode.isLeaf = node.isLeaf
	// remaining keys after median, add it to the rightnode
	
	j := 0
	for i := medianIndex + 1; i < node.numKeys; i++ {
		rightNode.keys[j] = node.keys[i]
		rightNode.numKeys++
		j++
	}

	// keep all the keys in node.keys till the median,
	for i := medianIndex; i < node.numKeys; i++ {
		node.keys[i] = 0 // Assuming keys are int type; otherwise, use zero value of the key type
	}
	node.numKeys = medianIndex


	// remaining children after median + 1 send it to the rightnode
	if !node.isLeaf {
		// copy(rightNode.children[:], node.children[medianIndex+1:node.numKeys+1])
		// log.Println("moving children to rightnode", rightNode)
		j := 0
		for i := medianIndex + 1; i < M + 1; i++ {
			// log.Println(node.children[i].keys)
			rightNode.children[j] = node.children[i]
			node.children[i] = nil
			rightNode.children[j].parent = rightNode
			j++
		}
	}
	// if children in node, then keep all children in node.children till, median(index)
	// copy(node.children[:], node.children[:medianIndex + 1])
	
	// else no children, then do nothing

	log.Println("current node", node.keys)
	log.Println("current node's children")
	for i := 0; i < node.numKeys + 1; i++ {
		if node.children[i] == nil {
			log.Println("nil")
		}else{
			log.Println(node.children[i].keys)
		}
	}

	log.Println("right node", rightNode.keys)
	log.Println("right node's children")
	for i := 0; i < rightNode.numKeys + 1; i++ {
		if rightNode.children[i] == nil {
			log.Println("nil")
		}else{
			log.Println(rightNode.children[i].keys)
		}
	}


	return rightNode, median
}

func (node *Node) insert(key int, tree *Tree) {
	// if the node has space, 
	// insert into the node by shifting elements accordingly

	pos := 0
	if node.isLeaf {
		if node.numKeys == MAX_NUM_OF_KEYS {
			log.Println("max limit reached")
			// this is a leaf node with maximum number of keys
			// we have to insert and split and check them recursively for maxNumberOfKeys
			// so, every node should have access to its parent node
			node.maxKeyThresholdReached(key, nil, tree, &pos)
		}else{
			// this is a leaf node with space to add another key, so we just insert it
			node.insertIntoNode(key);
			return
		}
	}else{
		node.searchNode(key, &pos)
		log.Println("going to child node at index", pos)
		node.children[pos].insert(key, tree)
	}
}

func NewNode() *Node {
	return &Node{
		numKeys: 0,
		keys: [M]int{},
		children: [M + 1]*Node{},
		parent: nil,
		isLeaf: true,
	}
}

type Tree struct {
	root *Node
	maxKeys int
}

func (tree *Tree) Find(key int) *Node {
	currentNode := tree.root
	result := currentNode.search(key)
	return result
}

func (tree *Tree) Put(key int) {
	if tree.root == nil {
		tree.root = NewNode()
	}
	currentNode := tree.root
	// pos := 0
	currentNode.insert(key, tree)
}

func (tree *Tree) Del(key int) {
	// TODO:
	//  THIS IS DELETION FROM A LEAF NODE
	// LOOK OUT : updating parent pointers, shifting children if necessary, updating numOfKeys
	// 1. first find the node and delete the node,
	// 		and do the appropriate element shift, and update the numOfKeys in the node.
	// 2. check if node.numOfKeys >= MIN_NUM_OF_KEYS
	// 		if node.numOfKeys >= MIN_NUM_OF_KEYS {
	// 			in this case, just return. it follows the btree property
	// 		}else{
	// 			=> we try to borrow
	// 			=> let us assume, currentNode is our node, in which we deleted the key,
	// 				and is now not following the propety.
	// 			=> let us assume,
	// 				 we have the currentNode's position in its parent's children array
	// 				 For eg: assume root = 40. and the MIN_NUM_OF_KEYS = 2
	// 						50 => root
	// 					   /  \
	// 			     |20, 40|  |60, 70|
	// 			 	 and we want to delete the key, 60.
	// 				 NOW, 60 is at "position" 1 of the children array of root.
	// 				  root.children[1] = *ptr to the node[40, 60]
	// 					above is the position
	// 
	// 			=> so we have currentNode and the position
	// 			=> now we check the currentNode's left sibling
	// 
	// 			=> if currentNode.parent.children[position - 1].numOfKeys > MIN_NUM_OF_KEYS {
	//  			it means the left sibling has more than required keys. so it can spare
	// 				so, we take the separater(the key in the parentNode which separates the siblings) and,
	// 					 put it in the currentNode and shift elements accordingly, and update the numOfKeys.
	// 				and then we take the rightmost key in the left sibling and put it in the separater's place,
	// 				and update the numOfKeys in the left sibling
	// 				return;
	// 			   }
	// 			=> left sibling does not have keys to spare, so we check for right sibling
	// 
	// 			   else if currentNode.parent.children[position + 1].numOfKeys > MIN_NUM_OF_KEYS {
	//				it means the right sibling has more than required keys. so it can spare
	// 			    so, we take the separater(the key in the parentNode which separates the siblings) and,
	// 					put it in the currentNode and shift elements accordingly, and update the numOfKeys.
	// 				and then we take the leftmost key in the right sibling and put it in the separater's place, 
	// 				and update the numOfKeys in the right sibling
	// 				return;
	// 			   }
	//			=> now, both left and right sibling does not have keys to spare, so
	// 				if left sibling != nil {
	// 				 we combine our currentNode with the left sibling, with the separater
	// 				 eg: currentNode = left sibling + separater + currentNode
	// 				 and update the numOfKeys of the parent, 
	// 				 because we just brought separater down to currentNode.
	//              }else if(right sibling != nil) {
	// 				  we combine our currentNode with the right sibling, with the separater
	// 				  eg: currentNode = currentNode + separater + right sibling
	// 				  and update the numOfKeys of the parent, 
	// 				  because we just brought separater down to currentNode sibling.
	// 				}
	// 				
	// 			
	// 			=>	lets asumme, i call the node, newCurrentNode which is the parent of
	// 				 (left sibling and current) or (right sibling and current) from which we just brought separater down.
	// 				now, that newCurrentNode has lost its separater key.
	// 			=> after joining with either the left or right sibling, check is newCurrentNode.parent == nil
	// 				if nil, make the newCurrentNode as root.
	// 				recursively call this same function() with node as newCurrentNode
	// 		}

	// THIS IS DELETION FROM AN INTERNAL NODE
	// 1. Find the successor key
	// 		successor key is the smallest key in the right subtree(always in leaf node)
	// 		steps :  => move to the right child
	// 				 => keep on moving to the 0th child of each node, till we reach a leaf node
	// 				 => leftmost key in that leaf node is the successor key
	// 2. Copy the successor key at the place of the key to be deleted
	// 3. Delete the successor key (deletion from a leaf node)
}

func NewTree() *Tree {
	return &Tree{
		root: nil,
		maxKeys: MAX_NUM_OF_KEYS,
	}
}


func (node *Node) inorder(result *[]int) {
	if node == nil {
	  return
	}
  
	if node.isLeaf {
	  // Print the keys in the leaf node
	
	  for i := 0; i < node.numKeys; i++ {
		*result = append(*result, node.keys[i])
	  }

	  return
	}
  
	// Traverse child nodes and keys
	for i := 0; i <= node.numKeys; i++ { // process all children (0 to numKeys+1)
	  // Recursively traverse the i-th child
	  node.children[i].inorder(result)
  
	  // Append the i-th key (if applicable)
	  if i < node.numKeys {
		*result = append(*result, node.keys[i])
	  }
	}
  }

func (tree *Tree) Print() []int {
    if tree.root == nil {
        return []int{}
    }

    result := []int{}
    tree.root.inorder(&result)
    return result
}
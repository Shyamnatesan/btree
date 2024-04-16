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

func (node *Node) searchNode(key int, pos *int) (*Node, int) {
	// if found return the node
	// else return nil
	log.Println("searching for key ", key, " in the node ", node.keys)
	for ; *pos < node.numKeys; {
		if key > node.keys[*pos] {
			*pos++;
		}else if key ==  node.keys[*pos]{
			log.Println("found a match at position ", *pos, "in node ", node.keys)
			return node, *pos;
		}else{
			break
		}
	}
	return nil, -1
}

func (node *Node) search(key int) (*Node, int) {
	if node == nil {
		return nil, -1
	}
	pos := 0

	found, indexOfKey := node.searchNode(key, &pos)
	if found != nil {
		return found, indexOfKey
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
			node.parent.maxKeyThresholdReached(median, rightNode, tree, pos)

		}
	}
}

func (node *Node) splitNode() (*Node, int) {
	// log.Println("splitting the node ", node.keys)
	medianIndex := node.numKeys / 2
	if medianIndex % 2 == 0 {
		medianIndex--
	}
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

func (tree *Tree) Find(key int) (*Node, int) {
	currentNode := tree.root
	result, posOfKey := currentNode.search(key)
	return result, posOfKey
}

func (tree *Tree) Put(key int) {
	if tree.root == nil {
		tree.root = NewNode()
	}
	currentNode := tree.root
	// pos := 0
	currentNode.insert(key, tree)
}

// [1, 2, 3, 4, 5]
func (node *Node) deleteKeyInNode(pos int) {
	n := node.numKeys
	i := pos + 1
	for ; i < n; i++ {
		node.keys[i - 1] = node.keys[i]
	}
	node.keys[i - 1] = 0
	node.numKeys--
}

func (node *Node) getSiblings() (*Node, int, string) {
	parentNode := node.parent
	var leftSibling *Node
	var separaterIndex int
	var nodeIndexInParent int

	// Find the index of the node in the parent's children
	for i := 0; i < parentNode.numKeys + 1; i++ {
		if parentNode.children[i] == node {
			nodeIndexInParent = i
			break
		}
	}
	// If the node is not the first child, then left sibling exists
	if nodeIndexInParent > 0 {
		// left sibling exists
		leftSibling = parentNode.children[nodeIndexInParent - 1]
		separaterIndex = nodeIndexInParent - 1
		return leftSibling, separaterIndex, "left"
		
	}
	var rightSibling *Node
	// If the node is not the last child, then right sibling exists
	if nodeIndexInParent < parentNode.numKeys {
		// right sibling exists
		rightSibling = parentNode.children[nodeIndexInParent + 1]
		separaterIndex = nodeIndexInParent
		return rightSibling, separaterIndex, "right"	
		
	}

	return nil, -1, ""
}

func (node *Node) borrowFromLeftSibling(leftSibling *Node, separaterIndex int) {
	parentNode := node.parent
	rightMostkeyInLeftSibling := leftSibling.keys[leftSibling.numKeys - 1]
	separater := parentNode.keys[separaterIndex]
	node.insertIntoNode(separater)
	parentNode.keys[separaterIndex] = rightMostkeyInLeftSibling
	leftSibling.deleteKeyInNode(leftSibling.numKeys - 1)
}

func (node *Node) borrowFromRightSibling(rightSibling *Node, separaterIndex int) {
	parentNode := node.parent
	leftMostKeyInRightSibling := rightSibling.keys[0]
	separater := parentNode.keys[separaterIndex]
	node.insertIntoNode(separater)
	parentNode.keys[separaterIndex] = leftMostKeyInRightSibling
	rightSibling.deleteKeyInNode(0)
}

func (node *Node) rebalancing() {
	if node.numKeys >= MIN_NUM_OF_KEYS {
		return
	}

	sibling, separaterIndex, siblingType := node.getSiblings()
	if sibling == nil {
		log.Println("no sibling to borrow from")
		// we merge with parent here
	}
	borrowed := false
	switch siblingType {
		case "left":
			if sibling.numKeys > MIN_NUM_OF_KEYS {
				node.borrowFromLeftSibling(sibling, separaterIndex)
				borrowed = true
				break
			}
		case "right":
			if sibling.numKeys > MIN_NUM_OF_KEYS {
				node.borrowFromRightSibling(sibling, separaterIndex)
				borrowed = true
				break
			}
	}
	if !borrowed {
		// if not borrowed we try merging with one of the sibling
	}
}

func (tree *Tree) Del(key int) {
	node, posOfKey := tree.root.search(key)
	if node == nil {
		log.Println("key does not exist in the btree")
		return
	}
	if node.isLeaf {
		// key is in a leaf node
		// 1. delete the key
		node.deleteKeyInNode(posOfKey)
		if node.numKeys >= MIN_NUM_OF_KEYS {
			// no need to rebalance, since the noOfKeys is >= minimum number of keys in a node threshold
			return
		}
		node.rebalancing()
		// 2. check the node for rebalancing
		//     - borrow from left sibling, if left sibling can spare keys (or)
		// 	   - borrow from right sibling, if right sibling can spare keys
		// 	   - if left sibling and right sibling cannot spare keys, then
		// 	   - merge with left sibling, if left sibling is not nil (or)
		// 	   - merge with right sibling, if right sibling is not nil
		//     - when merging with left or right, the parent separater key comes down(removed from the parent node)
		// 	   - now recursively call this parentNode for rebalancing

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

func NewTree() *Tree {
	return &Tree{
		root: nil,
		maxKeys: MAX_NUM_OF_KEYS,
	}
}
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
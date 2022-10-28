package go_redblacktree

type RedBlackTree struct {
	root       *RedBlackTreeNode
	Comparator NodeValueEntryComparator
}

// findReturnLastNode: if it finds a node has the same value, return this node. Otherwise, return
// a leaf node which can be the parent of this value
func (tree *RedBlackTree) findReturnLastNode(val RedBlackTreeNodeValEntry) (node *RedBlackTreeNode, ok bool) {
	n := tree.root
OuterLoop:
	for {
		switch {
		case tree.Comparator.Equal(n.val, val):
			return n, true
		case tree.Comparator.Smaller(n.val, val):
			if nil == n.rightChild {
				break OuterLoop
			}
			n = n.rightChild
		default:
			if nil == n.leftChild {
				break OuterLoop
			}
			n = n.leftChild
		}
	}
	return n, false
}

// FindNode return a node has the same value. If not found, return nil
func (tree *RedBlackTree) FindNode(val RedBlackTreeNodeValEntry) (node *RedBlackTreeNode, ok bool) {
	if nil == tree.root || nil == tree.Comparator || nil == val {
		return nil, false
	}

	if n, ok := tree.findReturnLastNode(val); ok {
		return n, true
	} else {
		return nil, false
	}
}

// AddNode add a new node to the tree, then return the new node and whether add operation is successful or not.
// If the value is already exit, return the exit node
func (tree *RedBlackTree) AddNode(val RedBlackTreeNodeValEntry) (node *RedBlackTreeNode, ok bool) {
	if nil == tree.Comparator {
		return nil, false
	}

	if nil == tree.root {
		tree.root = &RedBlackTreeNode{
			val:   val,
			color: RedBlackTreeNodeColorBlack,
		}
		return tree.root, true
	}

	n, isOK := tree.findReturnLastNode(val)
	if isOK {
		return n, true
	}
	newN := RedBlackTreeNode{
		parent: n,
		val:    val,
		color:  RedBlackTreeNodeColorRed,
	}
	if tree.Comparator.Smaller(n.val, val) {
		n.rightChild = &newN
	} else {
		n.leftChild = &newN
	}
	tree.adjustAfterAdd(&newN)
	return &newN, true
}

func (tree *RedBlackTree) adjustAfterAdd(node *RedBlackTreeNode) {
	if nil != node.parent &&
		RedBlackTreeNodeColorRed == node.color && RedBlackTreeNodeColorRed == node.parent.color {
		var parent *RedBlackTreeNode = node.parent
		var grand *RedBlackTreeNode = node.parent.parent
		var uncle *RedBlackTreeNode = nil
		isNodeLeft := false
		isParentLeft := false
		if nil == grand {
			// Grand is impossible to be nil. If occured, there must be
			// something wrong, panic!
			panic(false)
		}

		if parent.leftChild == node {
			isNodeLeft = true
		} else {
			isNodeLeft = false
		}
		if grand.leftChild == parent {
			uncle = grand.rightChild
			isParentLeft = true
		} else {
			uncle = grand.leftChild
			isParentLeft = false
		}
		switch {
		case nil != uncle && RedBlackTreeNodeColorRed == uncle.color:
			parent.color = RedBlackTreeNodeColorBlack
			uncle.color = RedBlackTreeNodeColorBlack
			grand.color = RedBlackTreeNodeColorRed
			if tree.root == grand {
				grand.color = RedBlackTreeNodeColorBlack
			} else {
				tree.adjustAfterAdd(grand)
			}
		default:
			switch {
			case isParentLeft && isNodeLeft:
				parent.color = RedBlackTreeNodeColorBlack
				grand.color = RedBlackTreeNodeColorRed
				tree.rightRotation(grand)
			case isParentLeft && !isNodeLeft:
				node.color = RedBlackTreeNodeColorBlack
				grand.color = RedBlackTreeNodeColorRed
				tree.leftRotation(parent)
				tree.rightRotation(grand)
			case !isParentLeft && isNodeLeft:
				node.color = RedBlackTreeNodeColorBlack
				grand.color = RedBlackTreeNodeColorRed
				tree.rightRotation(parent)
				tree.leftRotation(grand)
			case !isParentLeft && !isNodeLeft:
				parent.color = RedBlackTreeNodeColorBlack
				grand.color = RedBlackTreeNodeColorRed
				tree.leftRotation(grand)
			}
		}
	}
}

// RemoveNodeByVal remove a node which has the same value. If this node is not exit, do nothing
func (tree *RedBlackTree) RemoveNodeByVal(val RedBlackTreeNodeValEntry) {
	if node, ok := tree.FindNode(val); ok {
		tree.RemoveNode(node)
	}
}

// RemoveNode remove a node from tree
func (tree *RedBlackTree) RemoveNode(node *RedBlackTreeNode) {
	if nil == node {
		return
	}
	parent := node.parent
	isLeafNode := (nil == node.rightChild && nil == node.leftChild)
	if isLeafNode {
		switch {
		case nil == parent:
			// This case means the only one node in the tree is the ROOT node,
			// remove it directly.
			tree.root = nil
		case RedBlackTreeNodeColorRed == node.color:
			if parent.leftChild == node {
				parent.leftChild = nil
			} else {
				parent.rightChild = nil
			}
		case RedBlackTreeNodeColorBlack == node.color:
			isNodeLeftChild := false
			var brother *RedBlackTreeNode = nil
			if parent.leftChild == node {
				isNodeLeftChild = true
				brother = parent.rightChild
			} else {
				isNodeLeftChild = false
				brother = parent.leftChild
			}
			if nil != brother && RedBlackTreeNodeColorRed == brother.color {
				brother.color, parent.color = RedBlackTreeNodeColorBlack, RedBlackTreeNodeColorRed
				if isNodeLeftChild {
					tree.leftRotation(parent)
				} else {
					tree.rightRotation(parent)
				}
			}
			tree.adjustAfterChildRemoved(parent, isNodeLeftChild)
			if isNodeLeftChild {
				parent.leftChild = nil
			} else {
				parent.rightChild = nil
			}
		}
	} else {
		switch {
		case
			RedBlackTreeNodeColorBlack == node.color &&
				(nil != node.leftChild && RedBlackTreeNodeColorRed == node.leftChild.color &&
					nil == node.leftChild.leftChild && nil == node.leftChild.rightChild &&
					nil == node.rightChild ||
					nil != node.rightChild && RedBlackTreeNodeColorRed == node.rightChild.color &&
						nil == node.rightChild.leftChild && nil == node.rightChild.rightChild &&
						nil == node.leftChild):
			var child *RedBlackTreeNode = nil
			if nil != node.leftChild {
				child = node.leftChild
			} else {
				child = node.rightChild
			}
			if nil != parent {
				if parent.leftChild == node {
					parent.leftChild = child
				} else {
					parent.rightChild = child
				}
			}
			child.parent = parent
			child.color = RedBlackTreeNodeColorBlack
		default:
			// Use the NEXT node first
			var removeNode *RedBlackTreeNode = nil
			if removeNode = node.NextNode(); nil == removeNode {
				removeNode = node.PrevNode()
			}
			node.val = removeNode.val
			tree.RemoveNode(removeNode)
		}
	}
}

func (tree *RedBlackTree) adjustAfterChildRemoved(parent *RedBlackTreeNode, isRemovedNodeLeftChild bool) {
	var brotherOfRemovedNode *RedBlackTreeNode = nil
	var isBrotherLeftChild bool = !isRemovedNodeLeftChild
	if isRemovedNodeLeftChild {
		brotherOfRemovedNode = parent.rightChild
	} else {
		brotherOfRemovedNode = parent.leftChild
	}

	switch {
	case RedBlackTreeNodeColorRed == brotherOfRemovedNode.color:
		brotherOfRemovedNode.color, parent.color = RedBlackTreeNodeColorBlack, RedBlackTreeNodeColorRed
		if isRemovedNodeLeftChild {
			tree.leftRotation(parent)
		} else {
			tree.rightRotation(parent)
		}
		tree.adjustAfterChildRemoved(parent, isRemovedNodeLeftChild)
	case
		(nil == brotherOfRemovedNode.leftChild || nil != brotherOfRemovedNode.leftChild &&
			RedBlackTreeNodeColorBlack == brotherOfRemovedNode.leftChild.color) &&
			(nil == brotherOfRemovedNode.rightChild || nil != brotherOfRemovedNode.rightChild &&
				RedBlackTreeNodeColorBlack == brotherOfRemovedNode.rightChild.color):
		oldParentColor := parent.color
		brotherOfRemovedNode.color, parent.color = RedBlackTreeNodeColorRed, RedBlackTreeNodeColorBlack
		if RedBlackTreeNodeColorBlack == oldParentColor {
			if nil != parent.parent {
				parentLeftChild := false
				if parent.parent.leftChild == parent {
					parentLeftChild = true
				} else {
					parentLeftChild = false
				}
				tree.adjustAfterChildRemoved(parent.parent, parentLeftChild)
			}
		}
	default:
		var cencerNodeAfterRotation *RedBlackTreeNode = nil
		cencerNodeColor := parent.color
		brotherHasLeftRedChild := false
		brotherHasRightRedChild := false
		if nil != brotherOfRemovedNode.leftChild && RedBlackTreeNodeColorRed == brotherOfRemovedNode.leftChild.color {
			brotherHasLeftRedChild = true
		}
		if nil != brotherOfRemovedNode.rightChild && RedBlackTreeNodeColorRed == brotherOfRemovedNode.rightChild.color {
			brotherHasRightRedChild = true
		}
		switch {
		case brotherHasLeftRedChild && brotherHasRightRedChild:
			cencerNodeAfterRotation = brotherOfRemovedNode
			if isBrotherLeftChild {
				tree.rightRotation(parent)
			} else {
				tree.leftRotation(parent)
			}
		case brotherHasLeftRedChild && !brotherHasRightRedChild:
			if isBrotherLeftChild {
				cencerNodeAfterRotation = brotherOfRemovedNode
				tree.rightRotation(parent)
			} else {
				tree.rightRotation(brotherOfRemovedNode)
				cencerNodeAfterRotation = parent.rightChild
				tree.leftRotation(parent)
			}
		case !brotherHasLeftRedChild && brotherHasRightRedChild:
			if isBrotherLeftChild {
				tree.leftRotation(brotherOfRemovedNode)
				cencerNodeAfterRotation = parent.leftChild
				tree.rightRotation(parent)
			} else {
				cencerNodeAfterRotation = brotherOfRemovedNode
				tree.leftRotation(parent)
			}
		}
		cencerNodeAfterRotation.color = cencerNodeColor
		if nil != cencerNodeAfterRotation.leftChild {
			cencerNodeAfterRotation.leftChild.color = RedBlackTreeNodeColorBlack
		}
		if nil != cencerNodeAfterRotation.rightChild {
			cencerNodeAfterRotation.rightChild.color = RedBlackTreeNodeColorBlack
		}
	}
}

func (tree *RedBlackTree) leftRotation(node *RedBlackTreeNode) {
	if nil == node || nil == node.rightChild {
		return
	}

	if tree.root == node {
		tree.root = node.rightChild
	}
	parent := node.parent
	rightChild := node.rightChild
	rightChildLeftChild := node.rightChild.leftChild
	if nil != parent {
		if parent.leftChild == node {
			parent.leftChild = rightChild
		} else {
			parent.rightChild = rightChild
		}
	}
	rightChild.parent = parent
	node.rightChild = rightChildLeftChild
	if nil != rightChildLeftChild {
		rightChildLeftChild.parent = node
	}
	rightChild.leftChild, node.parent = node, rightChild
}

func (tree *RedBlackTree) rightRotation(node *RedBlackTreeNode) {
	if nil == node || nil == node.leftChild {
		return
	}

	if tree.root == node {
		tree.root = node.leftChild
	}
	parent := node.parent
	leftChild := node.leftChild
	leftChildRightChild := node.leftChild.rightChild
	if nil != parent {
		if parent.leftChild == node {
			parent.leftChild = leftChild
		} else {
			parent.rightChild = leftChild
		}
	}
	leftChild.parent = parent
	node.leftChild = leftChildRightChild
	if nil != leftChildRightChild {
		leftChildRightChild.parent = node
	}
	leftChild.rightChild, node.parent = node, leftChild
}

package go_redblacktree

type RedBlackTreeNodeColor byte

const (
	RedBlackTreeNodeColorInvalid RedBlackTreeNodeColor = iota
	RedBlackTreeNodeColorBlack
	RedBlackTreeNodeColorRed
)

// RedBlackTreeNodeValEntry is the value wrapper of RedBlackTreeNode, whitch need to
// implement the comparison operator
type RedBlackTreeNodeValEntry interface {
	Equal(b RedBlackTreeNodeValEntry) bool
	Smaller(b RedBlackTreeNodeValEntry) bool
}

// RedBlackTreeNode the node of Red-Black-Tree
type RedBlackTreeNode struct {
	leftChild  *RedBlackTreeNode
	rightChild *RedBlackTreeNode
	val        RedBlackTreeNodeValEntry
	color      RedBlackTreeNodeColor
	parent     *RedBlackTreeNode
}

// IsValidNode checks if this node valid
func (node RedBlackTreeNode) IsValidNode() bool {
	if RedBlackTreeNodeColorRed == node.color || RedBlackTreeNodeColorBlack == node.color {
		return true
	}
	return false
}

// Value return value of the node if node is valid
func (node RedBlackTreeNode) Value() RedBlackTreeNodeValEntry {
	if !node.IsValidNode() {
		return nil
	}
	return node.val
}

// PrevNode return the previous node by Inorder Traversal
func (node RedBlackTreeNode) PrevNode() *RedBlackTreeNode {
	if !node.IsValidNode() {
		return nil
	}
	var prevNode *RedBlackTreeNode = nil
	switch {
	case nil != node.leftChild:
		prevNode = node.leftChild
		for nil != prevNode.rightChild {
			prevNode = prevNode.rightChild
		}
	case nil != node.parent && &node == node.parent.rightChild:
		prevNode = node.parent
	default:
		n := &node
		for nil != n.parent {
			if n == n.parent.rightChild {
				prevNode = n.parent
				break
			}
			n = n.parent
		}
	}
	return prevNode
}

// NextNode return the next node by Inorder Traversal
func (node RedBlackTreeNode) NextNode() *RedBlackTreeNode {
	if !node.IsValidNode() {
		return nil
	}
	var nextNode *RedBlackTreeNode = nil
	switch {
	case nil != node.rightChild:
		nextNode = node.rightChild
		for nil != nextNode.leftChild {
			nextNode = nextNode.leftChild
		}
	default:
		n := &node
		for nil != n.parent {
			if n == n.parent.leftChild {
				nextNode = n.parent
				break
			}
			n = n.parent
		}
	}
	return nextNode
}

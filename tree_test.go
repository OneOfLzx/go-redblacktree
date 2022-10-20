package go_redblacktree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type RedBlackTreeError byte

const (
	RedBlackTreeErrorInvalidColor RedBlackTreeError = iota
	RedBlackTreeErrorLeftChildValNotSmaller
	RedBlackTreeErrorRightChildValNotBigger
	RedBlackTreeErrorDoubleLinkedRedNode
	RedBlackTreeErrorDifferentBlackPathLen
	RedBlackTreeErrorRedRootNode
)

var RedBlackTreeErrorToString = [...]string{
	RedBlackTreeErrorInvalidColor:           "Invalid Color",
	RedBlackTreeErrorLeftChildValNotSmaller: "Left Child Val Not Smaller",
	RedBlackTreeErrorRightChildValNotBigger: "Right Child Val Not Bigger",
	RedBlackTreeErrorDoubleLinkedRedNode:    "Double Linked Red Node",
	RedBlackTreeErrorDifferentBlackPathLen:  "Different Black Path Len",
	RedBlackTreeErrorRedRootNode:            "Red Root Node",
}

func (e RedBlackTreeError) String() string {
	idx := int(e)
	if len(RedBlackTreeErrorToString) <= idx {
		return fmt.Sprintf("Unknown error: %d", idx)
	}
	return RedBlackTreeErrorToString[idx]
}

func (e RedBlackTreeError) Error() string {
	return e.String()
}

type testRedBlackTreeNodeValEntry struct {
	RedBlackTreeNodeValEntry
	val int
}

func (e testRedBlackTreeNodeValEntry) Equal(b RedBlackTreeNodeValEntry) bool {
	return e.val == (b.(testRedBlackTreeNodeValEntry)).val
}

func (e testRedBlackTreeNodeValEntry) Smaller(b RedBlackTreeNodeValEntry) bool {
	return e.val < (b.(testRedBlackTreeNodeValEntry)).val
}

func TestRedBlackTree(t *testing.T) {
	const (
		testTreeNum  = 5
		minTestNum   = -100
		maxTestNum   = 100
		minOpPerTree = (maxTestNum - minTestNum) * 30
		maxOpPerTree = (maxTestNum - minTestNum) * 60
	)
	rand.Seed(time.Now().Unix())
	opCounts := maxOpPerTree - minOpPerTree
	opCounts = rand.Intn(opCounts) + minOpPerTree
OuterLoop:
	for tIdx := 0; tIdx < testTreeNum; tIdx++ {
		tree := RedBlackTree{}
		addTimes, removeTimes, findTimes := 0, 0, 0
		for i := 0; i < opCounts; i++ {
			entry := testRedBlackTreeNodeValEntry{val: rand.Intn(maxTestNum) + minTestNum}
			switch op := rand.Intn(3); op {
			case 0:
				tree.AddNode(entry)
				addTimes++
			case 1:
				tree.RemoveNodeByVal(entry)
				removeTimes++
			case 2:
				n, ok := tree.FindNode(entry)
				if ok {
					if !n.val.Equal(entry) {
						t.Error("tree ", tIdx, " Find return ok but val is not equal")
					}
				} else {
					if nil != n {
						t.Error("tree ", tIdx, " Find return fail but node is not nil")
					}
				}
				findTimes++
			}

			if err := isRedBlackTreeValid(&tree); nil != err {
				t.Error("tree ", tIdx, " BAD tree, error: ", err)
				break OuterLoop
			}
		}

		fmt.Println("tree ", tIdx, " done, root: ", (tree.root), " opCounts: ", opCounts,
			" addTimes: ", addTimes, " removeTimes: ", removeTimes, " findTimes: ", findTimes)
	}
}

func checkRedBlackTreeBlackChildCountsAndRedChild(root *RedBlackTreeNode) (int, error) {
	leftChildTreeBlackCounts := 0
	rightChildTreeBlackCounts := 0
	var err error = nil
	if RedBlackTreeNodeColorRed != root.color && RedBlackTreeNodeColorBlack != root.color {
		return -1, RedBlackTreeErrorInvalidColor
	}
	if nil != root.leftChild {
		if !root.leftChild.val.Smaller(root.val) {
			return -1, RedBlackTreeErrorLeftChildValNotSmaller
		}

		if RedBlackTreeNodeColorRed == root.color && RedBlackTreeNodeColorRed == root.leftChild.color {
			return -1, RedBlackTreeErrorDoubleLinkedRedNode
		}

		leftChildTreeBlackCounts, err = checkRedBlackTreeBlackChildCountsAndRedChild(root.leftChild)
		if nil != err {
			return -1, err
		}
	}
	if nil != root.rightChild {
		if !root.val.Smaller(root.rightChild.val) {
			return -1, RedBlackTreeErrorRightChildValNotBigger
		}

		if RedBlackTreeNodeColorRed == root.color && RedBlackTreeNodeColorRed == root.rightChild.color {
			return -1, RedBlackTreeErrorDoubleLinkedRedNode
		}

		rightChildTreeBlackCounts, err = checkRedBlackTreeBlackChildCountsAndRedChild(root.rightChild)
		if nil != err {
			return -1, err
		}
	}
	if leftChildTreeBlackCounts == rightChildTreeBlackCounts {
		if RedBlackTreeNodeColorBlack == root.color {
			leftChildTreeBlackCounts++
		}
		return leftChildTreeBlackCounts, nil
	} else {
		return -1, RedBlackTreeErrorDifferentBlackPathLen
	}
}

func isRedBlackTreeValid(tree *RedBlackTree) error {
	if nil == tree.root {
		return nil
	}
	if RedBlackTreeNodeColorBlack != tree.root.color {
		return RedBlackTreeErrorRedRootNode
	}
	_, err := checkRedBlackTreeBlackChildCountsAndRedChild(tree.root)
	return err
}

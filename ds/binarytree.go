package main

import "fmt"

// BTree represents a binary tree
type BTree struct {
	val int

	left  *BTree
	right *BTree
}

// NewBTree returns a new binary tree
func NewBTree(val int) *BTree {
	return &BTree{
		val: val,

		left:  nil,
		right: nil,
	}
}

// Print sends the values of the binary tree to stdout
func (bt *BTree) Print() {
	if bt == nil {
		return
	}
	fmt.Println(bt.val)

	bt.left.Print()
	bt.right.Print()
}

// Insert adds any value to the BTree
func (bt *BTree) Insert(val int) {
	if val >= bt.val {
		if bt.right == nil {
			bt.right = NewBTree(val)
		} else {
			bt.right.Insert(val)
		}
	} else {
		if bt.left == nil {
			bt.left = NewBTree(val)
		} else {
			bt.left.Insert(val)
		}
	}
}

// printBTree creates a sample btree and prints it to stdout
func printBTree() {
	bt := NewBTree(0)

	bt.Insert(5)
	bt.Insert(1)
	bt.Insert(99)
	bt.Insert(23)

	bt.Print()
}

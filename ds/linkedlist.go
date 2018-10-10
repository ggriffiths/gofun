package main

import "fmt"

// Node represents a list node
type Node struct {
	Val  interface{}
	Next *Node
	Prev *Node
}

// NewNode represents a new node
func NewNode(val interface{}) *Node {
	return &Node{
		Val:  val,
		Next: nil,
		Prev: nil,
	}
}

// LinkedList represents a double linked list
type LinkedList struct {
	head *Node
	tail *Node
}

// NewList creates a new linked list
func NewList(val interface{}) *LinkedList {
	n := NewNode(val)

	return &LinkedList{
		head: n,
		tail: n,
	}
}

// Append adds a new Node to the end of the list
func (ll *LinkedList) Append(val interface{}) *LinkedList {
	// Create new node object
	n := NewNode(val)

	// Attach to end of list
	ll.tail.Next = n

	// Reference back to old tail
	n.Prev = ll.tail

	// Set new tail
	ll.tail = n

	return ll
}

// Prepend adds a new Node to the beginning of the list
func (ll *LinkedList) Prepend(val interface{}) *LinkedList {
	// Create new node object
	n := NewNode(val)

	// Attach to end of list
	ll.head.Prev = n

	// Reference back to old tail
	n.Next = ll.head

	// Set new tail
	ll.head = n

	return ll
}

// Print sends the contents of the linked list to Stdout from head to tail
func (ll *LinkedList) Print() {
	current := ll.head

	for current != nil {
		fmt.Println(current.Val)
		current = current.Next
	}
}

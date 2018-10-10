package main

import "fmt"

func main() {
	// printBasicNode()
	// printLinkedList()

	printBTree()
}

func printLinkedList() {
	ll := NewList(0)

	nums := []int{192, 3, 5, 1, 5, 12, 3, 3}

	for _, n := range nums {
		ll.Append(n)
	}

	ll.Print()

	fmt.Println("---")

	nums = []int{1, 3, 9999, 1000, 2000, 3000, -1000}
	for _, n := range nums {
		ll.Prepend(n)
	}

	ll.Print()
}

func printBasicNode() {
	nums := []int{1, 6, 1, 2, 1, 201}
	head := NewNode(0)

	var current = head
	var old = head
	for _, n := range nums {
		// Keep track of prev node
		old = current

		// Add link to new node
		current.Next = NewNode(n)

		// Increment ptr
		current = current.Next

		// Add link to previous
		current.Prev = old
	}

	// Verify correctness
	current = head
	for current != nil {
		fmt.Println(current.Val)
		current = current.Next
	}
}

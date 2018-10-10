package main

import "fmt"

func fib(n int, history map[int]int) int {
	fmt.Println("n:", n, "history:", history)
	if n == 1 || n == 2 {
		return 1
	}

	var f1, f2 int
	var found bool
	if f1, found = history[n-1]; !found {
		f1 = fib(n-1, history)
		history[n-1] = f1
	}
	if f2, found = history[n-2]; !found {
		f2 = fib(n-2, history)
		history[n-2] = f2
	}

	return f1 + f2

}

// FibN returns the first N fibonaci numbers
func FibN(n int) []int {
	history := make(map[int]int)
	var fibs []int

	for i := 1; i <= n; i++ {
		fibs = append(fibs, fib(i, history))
	}

	return fibs
}

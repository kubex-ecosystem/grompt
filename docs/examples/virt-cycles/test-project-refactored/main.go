package main

import "fmt"

func main() {
	// Use idiomatic variable names (camelCase).
	num1 := 1
	num2 := 2

	// Use the add function from utils.go to calculate the sum.
	sum := add(num1, num2)

	// Print the sum using fmt.Println.  No error handling needed for this simple case.
	fmt.Println(sum)
}
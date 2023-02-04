// Exercise 5.19 from gopl.io

package main

import (
	"fmt"
)

type nonZeroReturn struct{}

func nonZeroPanicReturn() (result int) {

	defer func() {
		switch p := recover(); p {
		case nil:
		case nonZeroReturn{}:
			result = 3
		default:
			panic(p) // unexpected panic; carry on panicking
		}
	}()

	panic(nonZeroReturn{})
}

func main() {
	fmt.Println(nonZeroPanicReturn())
}

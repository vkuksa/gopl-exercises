// Exercises 1.1 1.2 1.3 from gopl.io

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Echo1(Args []string) string {
	var s, sep string
	for i := 1; i < len(Args); i++ {
		s += sep + Args[i]
		sep = " "
	}

	return s
}

func Echo2(Args []string) string {
	s, sep := "", ""
	for i, arg := range Args[0:] {
		s += sep + strconv.Itoa(i) + sep + arg
		sep = " "
	}

	return s
}

func Echo3(Args []string) string {
	return strings.Join(Args[0:], " ")
}

func main() {
	fmt.Println(Echo1(os.Args))
	fmt.Println(Echo2(os.Args))
	fmt.Println(Echo3(os.Args))
}

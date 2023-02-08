// Exercise 7.3: Write a String method for the *tree type in gopl.io/ch4/treesort (§4.4)
// that reveals the sequence of values in the tree.

package main

import (
	"bytes"
	"fmt"
)

type tree struct {
	value       int
	left, right *tree
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func (t *tree) String() string {
	var buf bytes.Buffer
	depth := 0
	var traverse func(*tree)
	traverse = func(t *tree) {
		if t == nil {
			return
		}

		depth++
		traverse(t.left)
		fmt.Fprintf(&buf, "%*s%d\n", int(1.0/float64(depth)*10), "", t.value)
		traverse(t.right)
		depth--
	}
	traverse(t)

	return buf.String()
}

func main() {
	t := &tree{4, nil, nil}
	add(t, 2)
	add(t, 3)
	add(t, 5)
	add(t, 6)
	add(t, 7)
	add(t, 1)

	fmt.Println(t)
}

package main

import (
	"fmt"
	"log"
)

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	// "linear algebra": {"calculus"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	sorted, err := topoSort(prereqs)
	if err != nil {
		log.Fatalf("topoSort: " + err.Error())
	}

	for i, course := range sorted {
		log.Printf("%d:\t%s", i+1, course)
	}
}

func topoSort(graph map[string][]string) ([]string, error) {
	var order []string
	seen := make(map[string]bool, len(graph))
	var visitAllIn func(string, string) error

	visitAllIn = func(current string, previous string) error {
		if !seen[current] {
			seen[current] = true
			for _, to := range graph[current] {
				if to == previous {
					return fmt.Errorf("cycle in a graph encountered: %s - %s", current, previous)
				}

				err := visitAllIn(to, current)
				if err != nil {
					return err
				}
			}

			order = append(order, current)
		}

		return nil
	}

	for node := range graph {
		err := visitAllIn(node, "")

		if err != nil {
			return nil, err
		}
	}
	return order, nil
}

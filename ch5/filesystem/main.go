// Exercise 5.14 from gopl.io

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func traverseDirectorySubtree(root string) []string {
	var directories []string
	entries, err := ioutil.ReadDir(root)
	if err != nil {
		log.Printf("reading info of %s failed with: %s", root, err.Error())
		return nil
	}

	for _, file := range entries {
		if file.IsDir() {
			dirName := root + file.Name() + "/"
			fmt.Println(dirName)
			directories = append(directories, dirName)
		}
	}

	return directories
}

func main() {
	breadthFirst(traverseDirectorySubtree, []string{os.Args[1]})
}

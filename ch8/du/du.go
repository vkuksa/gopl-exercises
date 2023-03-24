// Exercise 8.9: Write a version of du that computes and periodically displays separate totals for
// each of the root directories.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "show verbose progress messages")

func main() {
	flag.Parse()

	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go computeDiskUsage(root, func() { n.Done() })
	}

	n.Wait()
}

func computeDiskUsage(dir string, onCompletion func()) {
	defer onCompletion()

	var walingWorkers sync.WaitGroup
	fileSizes := make(chan int64)
	walingWorkers.Add(1)
	go walkDir(dir, &walingWorkers, fileSizes)
	go func() {
		walingWorkers.Wait()
		close(fileSizes)
	}()

	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64

loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(dir, nfiles, nbytes)
		}
	}

	printDiskUsage(dir, nfiles, nbytes)
}

func printDiskUsage(dir string, nfiles, nbytes int64) {
	fmt.Printf("%s: %d files  %.1f GB\n", dir, nfiles, float64(nbytes)/1e9)
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}

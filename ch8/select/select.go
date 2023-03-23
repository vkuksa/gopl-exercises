// Exercise 8.8: Using a select statement, add a timeout to the echo server from Section 8.3 so
// that it disconnects any client that shouts nothing within 10 seconds.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

// !+
func handleConn(c net.Conn) {
	defer c.Close()
	fmt.Println("Established connection...")

	messageReceived := make(chan bool)
	go func() {
		input := bufio.NewScanner(c)
		for input.Scan() {
			messageReceived <- true
			go echo(c, input.Text(), 1*time.Second)
		}
		// NOTE: ignoring potential errors from input.Err()
	}()

	for {
		select {
		case <-messageReceived:
		// Do nothing.
		case <-time.After(10 * time.Second):
			fmt.Println("No activity from client. Aborting connection...")
			return
		}
	}
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}

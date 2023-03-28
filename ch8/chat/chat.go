// Exercise 8.12: Make the broadcaster announce the current set of clients to each new arrival.
// This requires that the clients set and the entering and leaving channels record the client
// name too.
// Exercise 8.13: Make the chat server disconnect idle clients, such as those that have sent no
// messages in the last five minutes. Hint: calling conn.Close() in another goroutine unblocks
// active Read calls such as the one done by input.Scan().
// Exercise 8.14: Change the chat server’s network protocol so that each client provides its name
// on entering. Use that name instead of the network address when prefixing each message with
// its sender’s identity.
// Exercise 8.15: Failure of any client program to read data in a timely manner ultimately causes
// all clients to get stuck. Modify the broadcaster to skip a message rather than wait if a client
// writer is not ready to accept it. Alternatively, add buffering to each client’s outgoing message
// channel so that most messages are not dropped; the broadcaster should use a non-blocking
// send to this channel.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// !+broadcaster
type client struct {
	name    string
	msgChan chan<- string // an outgoing message channel
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients

	clientsToString := func() string {
		var result string
		for client := range clients {
			result += "\n" + client.name
		}
		return result
	}

	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for client := range clients {
				client.msgChan <- msg
			}

		case client := <-entering:
			clients[client] = true
			client.msgChan <- "List of active clients: " + clientsToString()

		case client := <-leaving:
			delete(clients, client)
			close(client.msgChan)
		}
	}
}

//!-broadcaster

// !+handleConn
func handleConn(conn net.Conn) {
	who := conn.RemoteAddr().String()
	ch := make(chan string) // outgoing client messages
	client := client{who, ch}

	go clientWriter(conn, ch)

	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- client
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

// !+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main

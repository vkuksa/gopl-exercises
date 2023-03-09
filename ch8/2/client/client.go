// Exercise 8.2: Implement a concurrent File Transfer Protocol (FTP) server. The server should
// interpret commands from each client such as cd to change directory, ls to list a directory, get
// to send the contents of a file, and close to close the connection. You can use the standard ftp
// command as the client, or write your own.

// This client straightforwardly sends minimum set of required commands
// Errors are not handled gracefully

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	BufferCapacity = 1024
)

func processPassiveConnection(port int, writer io.Writer) {
	address := "localhost:" + strconv.Itoa(port)
	passiveConn, err := net.Dial("tcp", address)
	log.Println("DEBUG: dialing " + address)
	if err != nil {
		log.Println("ERROR: passive connection: " + err.Error())
	}
	defer passiveConn.Close()

	log.Println("INFO: connected. writing... ")

	rw := bufio.NewReadWriter(bufio.NewReader(passiveConn), bufio.NewWriter(writer))
	_, err = rw.WriteTo(writer)
	if err != nil {
		log.Printf("ERROR: processing writer: %s", err.Error())
	}
	log.Println("INFO: data from passive connection processed successfully")
}

func main() {
	receiveBuffer := make([]byte, BufferCapacity)

	c, err := net.Dial("tcp", "localhost:11137")
	if err != nil {
		log.Fatalln("FATAL: Connection with server was not established: " + err.Error())
	}
	defer c.Close()

	writeCommand := func(cmd string) error {
		n, err := c.Write([]byte(cmd))
		if err != nil {
			return fmt.Errorf("writing to connection: " + err.Error())
		}

		log.Printf("DEBUG: %d bytes written: %s\n", n, cmd)
		c.Write([]byte("\n"))

		return nil
	}
	parseResponse := func() (code int, msg string, err error) {
		n, err := c.Read(receiveBuffer)
		if err != nil {
			err = fmt.Errorf("reading: %s\n", err)
			return
		}

		log.Printf("DEBUG: %d bytes read: %s\n", n, receiveBuffer[:n])

		if res := strings.SplitN(string(receiveBuffer[:n]), " ", 2); len(res) > 1 {
			code, err = strconv.Atoi(string(res[0]))
			msg = res[1]
			return
		}
		err = fmt.Errorf("Couldn't parse response code")
		return
	}
	expectCodeEqual := func(expected int) (code int, msg string, err error) {
		code, msg, err = parseResponse()
		if err != nil || code != expected {
			err = fmt.Errorf("unexpected code returned %d, expected %d", code, expected)
		}
		return
	}

	if _, _, err = expectCodeEqual(220); err != nil {
		log.Fatalln("FATAL: server handshake failed: %s" + err.Error())
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		res := strings.SplitN(line, " ", 2)
		command := res[0]
		var arg string
		if len(res) > 1 {
			arg = res[1]
		}

		switch command {
		case "USER":
			if err = writeCommand(line); err != nil {
				log.Println("ERROR: username: " + err.Error())
				continue
			}

			if _, _, err = expectCodeEqual(331); err != nil {
				log.Println("ERROR: username: " + err.Error())
			}
		case "PASS":
			if err = writeCommand(line); err != nil {
				log.Println("ERROR: password: " + err.Error())
				continue
			}

			_, _, err = expectCodeEqual(230)
			if err != nil {
				log.Println("ERROR: password: " + err.Error())
			}
		case "LIST":
			if err = writeCommand(line); err != nil {
				log.Println("ERROR: list: " + err.Error())
				continue
			}

			if _, _, err = expectCodeEqual(150); err != nil {
				log.Println("ERROR: list: " + err.Error())
				continue
			}

			processPassiveConnection(52789, os.Stdout)
			if _, _, err = expectCodeEqual(226); err != nil {
				log.Println("ERROR: list: " + err.Error())
			}
		case "CWD":
			if err = writeCommand(line); err != nil {
				log.Println("ERROR: cwd: " + err.Error())
				continue
			}
			if _, _, err := expectCodeEqual(250); err != nil {
				log.Println("ERROR: cwd: " + err.Error())
			}
		case "RETR":
			if err = writeCommand("SIZE " + arg); err != nil {
				log.Println("WARN: size: " + err.Error())
				continue
			}
			_, msg, err := expectCodeEqual(213)
			if err != nil {
				log.Println("WARN: size: " + err.Error())
			}
			size, err := strconv.Atoi(msg)
			if err != nil {
				log.Println("WARN: converting size from string: " + err.Error())
			}
			if size == 0 {
				log.Println("WARN: Filesize check will not be performed.")
			}

			if err = writeCommand("RETR " + arg); err != nil {
				log.Println("WARN: size: " + err.Error())
				continue
			}
			if _, _, err := expectCodeEqual(150); err != nil {
				log.Println("ERROR: retr: " + err.Error())
				continue
			}

			func() {
				file, err := os.Create(arg)
				if err != nil {
					log.Println("ERROR: retr: " + err.Error())
				}
				defer file.Close()

				processPassiveConnection(33477, file)
			}()

			if _, _, err := expectCodeEqual(226); err != nil {
				log.Println("ERROR: retr: " + err.Error())
			}
		case "QUIT":
			writeCommand(line)
			if _, _, err := expectCodeEqual(221); err != nil {
				log.Println("WARN: quit: " + err.Error())
			}
			break
		default:
			fmt.Println("Unsupported command")
		}
	}
	if scanner.Err() != nil {
		log.Println("scanner: ", scanner.Err())
	}
}

// Exercise 8.2: Implement a concurrent File Transfer Protocol (FTP) server. The server should
// interpret commands from each client such as cd to change directory, ls to list a directory, get
// to send the contents of a file, and close to close the connection. You can use the standard ftp
// command as the client, or write your own.

// No implementation of authentication here
// Error handling is minimal

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
)

const (
	BufferCapacity = 1024
)

type connPanic string

var commandRegexp *regexp.Regexp = regexp.MustCompile(`([A-Z]{3,4})(\s)?([A-Za-z/._0-9]*)`)

func extractCommandAndArgs(input string) (cmd, args string) {
	groups := commandRegexp.FindStringSubmatch(input)

	return string(groups[1]), string(groups[3])
}

func transferDataThroughPassiveConnection(port int, dataGenerator func(io.Writer) error) {
	var passiveModeListener net.Listener
	var dataTransferringConn net.Conn

	defer func() {
		dataTransferringConn.Close()
		passiveModeListener.Close()
	}()

	address := "localhost:" + strconv.Itoa(port)
	passiveModeListener, err := net.Listen("tcp", address)
	log.Println("DEBUG: passive listening " + address)
	if err != nil {
		log.Println("ERROR: passive listening: " + err.Error())
	}

	dataTransferringConn, err = passiveModeListener.Accept()
	if err != nil {
		log.Print("ERROR: passive accepting: " + err.Error())
	}

	connWriter := bufio.NewWriter(dataTransferringConn)
	if err = dataGenerator(connWriter); err != nil {
		log.Printf("WARN: generating data: %s", err.Error())
	}

	if err = connWriter.Flush(); err != nil {
		log.Printf("WARN: flushing: %s", err.Error())
	}
	log.Println("DEBUG: finalizing connection to " + address)
}

func hasSubdir(rootDir, subDir string) (hasSubdir bool, err error) {
	entries, err := ioutil.ReadDir(rootDir)
	if err != nil {
		err = fmt.Errorf("reading dir %s: %s", rootDir, err.Error())
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if name == subDir && entry.IsDir() {
			hasSubdir = true
			return
		}
	}

	return
}

func handleConn(c net.Conn) {
	defer func() {
		c.Close()
		log.Println("DEBUG: closing connection")

		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	workingDir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	log.Println("DEBUG: processing client")

	_, err = c.Write([]byte("220 (ftpserver)"))
	if err != nil {
		log.Panic(err)
	}

	scanner := bufio.NewScanner(c)
	writer := bufio.NewWriter(c)

	writeResponse := func(str string) {
		if _, err = writer.WriteString(str); err != nil {
			log.Printf("writing: %s", err.Error())
		}
		if err = writer.Flush(); err != nil {
			log.Printf("flushing: %s", err.Error())
		}
	}

	for scanner.Scan() {
		input := scanner.Text()
		log.Printf("DEBUG: %s", input)

		switch cmd, arg := extractCommandAndArgs(input); cmd {
		case "USER":
			writeResponse("331 Please specify the password.")
		case "PASS":
			writeResponse("230 Login successful.")
		case "LIST":
			// EPSV communication emulation ommited here
			writeResponse("150 Here comes the directory listing.")

			directoryListGenerator := func(writer io.Writer) error {
				entries, err := ioutil.ReadDir(workingDir)
				if err != nil {
					return fmt.Errorf("reading dir %s: %s", workingDir, err.Error())
				}

				for _, entry := range entries {
					name := entry.Name()
					if entry.IsDir() {
						name += "/"
					}
					if _, err = fmt.Fprintf(writer, "%s\t%d\t%s\n", entry.Mode(), entry.Size(), name); err != nil {
						log.Printf("writing entry: %s with %s", entry, err.Error())
					}
				}

				return nil
			}
			transferDataThroughPassiveConnection(52789, directoryListGenerator)
			writeResponse("226 Directory send OK.")
		case "CWD":
			// ... is not supported. Only going in depth
			if ok, err := hasSubdir(workingDir, arg); !ok {
				if err != nil {
					writeResponse("550 No folder found.")
				} else {
					writeResponse("451 Error occured during.")
				}
			} else {
				workingDir += "/" + arg
				writeResponse("250 Directory successfully changed.")
			}
		case "SIZE":
			fstats, err := os.Stat(workingDir + "/" + arg)
			if err != nil {
				writeResponse("550 No file found.")
			}

			size := fstats.Size()
			writeResponse("213 " + strconv.Itoa(int(size)))
		case "RETR":
			writeResponse("150 Opening BINARY mode data connection for transfer.")

			fileContentReader := func(writer io.Writer) error {
				filepath := workingDir + "/" + arg
				file, err := os.Open(filepath)
				if err != nil {
					return fmt.Errorf("opening file %s: %s", filepath, err.Error())
				}
				defer file.Close()

				rw := bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(writer))
				_, err = rw.WriteTo(writer)
				if err != nil {
					return fmt.Errorf("sending file %s: %s", filepath, err.Error())
				}
				log.Printf("DEBUG: file %s was sent succesfully\n", filepath)
				return nil
			}
			transferDataThroughPassiveConnection(33477, fileContentReader)

			writeResponse("226 Transfer complete.")
		case "QUIT":
			writeResponse("221 Goodbye.")
		default:
			writeResponse("502 Not supported.")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner: %s", err.Error())
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:11137")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("ERROR: accepting conn: " + err.Error())
			continue
		}

		go handleConn(conn)
	}
}

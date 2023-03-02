// $ clockwall NewYork=localhost:8010 London=localhost:8020 Tokyo=localhost:8030

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"time"
)

type WatchSettings struct {
	city    string
	address string
}

var argRegexp *regexp.Regexp = regexp.MustCompile(`(\w*)=(localhost:\d{4})`)

func parseWatchSettings() []WatchSettings {
	watchSettings := make([]WatchSettings, 0)
	for i := 1; i < len(os.Args); i++ {
		groups := argRegexp.FindStringSubmatch(os.Args[i])
		watchSettings = append(watchSettings, WatchSettings{groups[1], groups[2]})
	}
	return watchSettings
}

func main() {
	settings := parseWatchSettings()

	for _, s := range settings {
		go pollTime(s)
	}

	for {
		time.Sleep(10 * time.Second)
	}
}

func pollTime(setting WatchSettings) {
	conn, err := net.Dial("tcp", setting.address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error reading from connection:", err)
			}
			return
		}

		data := buf[:n]
		fmt.Fprintf(os.Stdout, "%s\t: %s", setting.city, data)
	}
}

// Exercise 4.12 from gopl.io

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	. "gopl-exercises/ch4/json/xkcd/match"
)

const (
	url            = "https://xkcd.com/%d/info.0.json"
	usage          = `xkcd <search_strings>`
	episodes       = 2730 // How to flexibly get number of episodes?
	matchThreshold = 3

	flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
)

var cacheRoot string
var indexPath string

type Episode struct {
	Month      string
	Num        int
	Link       string
	Year       string
	News       string
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

func retrieveDataOfEpisode(number int) []byte {
	resp, err := http.Get(fmt.Sprintf(url, number))
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	return buf.Bytes()
}

func writeDataToCache(data []byte, filename string) bool {
	cacheEntry, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		log.Printf("Couldn't open cache entry: %s", filename)
		return false
	}
	defer cacheEntry.Close()

	cacheEntry.Write(data)
	return true
}

func readDataFromCacheEntry(entry string) ([]byte, error) {
	cache, err := os.Open(cacheRoot + entry)
	if err != nil {
		return nil, fmt.Errorf("opening cache entry: %s", err)

	}
	defer cache.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(cache)
	if err != nil {
		log.Printf("reading from cache failed: %s", err.Error())
	}
	return buf.Bytes(), nil
}

func init() {
	cacheRoot = os.Getenv("HOME") + "/xkcd/"
	indexPath = cacheRoot + "index"
	log.Println("Cache root " + cacheRoot)

	if _, err := os.Stat(cacheRoot); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(cacheRoot, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		index, err := os.OpenFile(indexPath, flags, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer index.Close()

		// It could be parralelled by number of available threads
		// Index writing should be mutex protected tho
		for i := 1; i <= episodes; i++ {
			data := retrieveDataOfEpisode(i)
			if data == nil {
				continue
			}
			iStr := strconv.Itoa(i)
			log.Println("Writing: " + iStr + " of size " + strconv.Itoa(len(data)))

			filename := cacheRoot + iStr

			_, alt, found := bytes.Cut(data, []byte("alt\": \""))
			if !found {
				log.Fatal("Could not find a title in a JSON.")
			}
			alt, _, _ = bytes.Cut(alt, []byte("\""))

			ok := writeDataToCache(data, filename)
			if !ok {
				continue
			}

			index.WriteString(iStr + ":" + string(alt) + "\n")
			// Only alt is written into index. Maybe some keyword detection algorithm would be better?
		}

		if err := index.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	index, err := os.Open(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	fileScanner := bufio.NewScanner(index)
	fileScanner.Split(bufio.ScanLines)

	// Also, search terms could be expanded and standardized
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if FuzzyMatch(line, os.Args[1], matchThreshold) {
			number, _, _ := bytes.Cut([]byte(line), []byte(":"))
			numberStr := string(number)

			data, err := readDataFromCacheEntry(numberStr)
			if err != nil {
				log.Printf("reading cache entry %s failed with %s", numberStr, err.Error())
				continue
			}

			var episode *Episode
			err = json.Unmarshal(data, &episode)
			if err != nil {
				log.Printf("Failed unmarshalling of episode " + numberStr + " with " + err.Error())
			}
			fmt.Println(numberStr + ": " + episode.Img + "\n" + episode.Transcript)
		}
	}
}

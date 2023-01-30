// Exercise 4.13 from gopl.io

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	usage   = `poster <movie_title>`
	omdbUrl = `http://www.omdbapi.com/?apikey=69c941e9&t=%s`
)

var saveLocation string

func queryMovieData(movie_title string) []byte {
	resp, err := http.Get(fmt.Sprintf(omdbUrl, movie_title))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("No movie info queried")
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func downloadPosterIntoFile(posterUrl, filepath string) {
	resp, err := http.Get(string(posterUrl))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("No poster returned")
	}

	posterFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer posterFile.Close()
	_, err = io.Copy(posterFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func getValueOfKeyFromRawJson(data []byte, key string) []byte {
	searchFor := []byte(key + "\":\"")
	_, cutted, found := bytes.Cut(data, searchFor)
	if !found {
		log.Fatal("Could not find a key in a JSON.")
	}
	result, _, found := bytes.Cut(cutted, []byte("\""))
	if !found {
		log.Fatal("Invalid JSON provided")
	}

	return result
}

func init() {
	saveLocation = os.Getenv("HOME") + "/"
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(usage)
	}

	data := queryMovieData(url.PathEscape(os.Args[1]))

	// Yeah, could've been done with json.Decode
	posterUrl := getValueOfKeyFromRawJson(data, "Poster")
	title := getValueOfKeyFromRawJson(data, "Title")
	// Not very unicode friendly
	filepath := saveLocation + string(title) + ".jpg"

	downloadPosterIntoFile(string(posterUrl), filepath)
	fmt.Println("Saved image into " + filepath)
}

// Exercises 5.5, 5.6 from gopl.io

package main

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

func corner(i, j int) (sx, sy float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	fmt.Printf("i %d, j %d, x %g, y %g, z %g\n", i, j, x, y, z)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx = width/2 + (x-y)*cos30*xyscale
	sy = height/2 + (x+y)*sin30*xyscale - z*zscale
	return
}

func visit(currentWords, currentImages int, n *html.Node) (words, images int) {
	words = currentWords
	images = currentImages

	words += len(strings.Fields(n.Data))
	if n.Type == html.ElementNode && n.Data == "img" {
		images++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		words, images = visit(words, images, c)
	}
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	words, images = visit(0, 0, n)
	return
}

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func main() {
	url := "https://golang.org"

	words, images, err := CountWordsAndImages(url)
	if err != nil {
		fmt.Print("Counting failed: " + err.Error())
	}

	fmt.Printf("Words: %d\nImages:%d\n", words, images)
}

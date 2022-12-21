// Exercises 3.1, 3.2, 3.3, 3.4 from gopl.io

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

const (
	angle = math.Pi / 6 // angle of x, y axes (=30°)
)

var (
	sin30, cos30  = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)
	width, height = 600, 320
	cells         = 100
	xyrange       = 30    // axis ranges (-xyrange..+xyrange)
	xyscale       float64 // pixels per x or y unit
	zscale        float64 // pixels per z unit
)

var figure = flag.String("figure", "hypot", "figure: eggbox, moguls, saddle, hypot(default)")
var web = flag.Bool("web", false, "start program as webserver")

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()

	if *web {
		http.HandleFunc("/", handler)
		http.ListenAndServe("localhost:8000", nil)
	}

	makeSVG(os.Stdout, width, height, "white")
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("sample", "header")
	w.Header().Set("Content-Type", "image/svg+xml")

	var err error
	if r.URL.Query().Has("width") {
		if width, err = strconv.Atoi(r.URL.Query().Get("width")); err != nil {
			log.Printf("Couldnt retrieve width from request: %s. Using default %d", err, width)
		} else {
			log.Printf("Retrieved url parameter for width=%d", width)
		}
	}

	if r.URL.Query().Has("height") {
		if height, err = strconv.Atoi(r.URL.Query().Get("height")); err != nil {
			log.Printf("Couldnt retrieve height from request: %s. Using default %d", err, height)
		} else {
			log.Printf("Retrieved url parameter for height=%d", height)
		}
	}

	if r.URL.Query().Has("xyrange") {
		if xyrange, err = strconv.Atoi(r.URL.Query().Get("xyrange")); err != nil {
			log.Printf("Couldnt retrieve xyrange from request: %s. Using default %d", err, xyrange)
		} else {
			log.Printf("Retrieved url parameter for xyrange=%d", xyrange)
		}
	}

	if r.URL.Query().Has("cells") {
		if cells, err = strconv.Atoi(r.URL.Query().Get("cells")); err != nil {
			log.Printf("Couldnt retrieve cells from request: %s. Using default %d", err, cells)
		} else {
			log.Printf("Retrieved url parameter for cells=%d", cells)
		}
	}

	if r.URL.Query().Has("figure") {
		*figure = r.URL.Query().Get("figure")
	}

	var fill string
	if r.URL.Query().Has("fill") {
		fill = r.URL.Query().Get("fill")
	} else {
		fill = "white"
	}

	makeSVG(w, width, height, fill)
}

func makeSVG(out io.Writer, width, height int, fill string) {
	xyscale = float64(width / 2 / xyrange)
	zscale = float64(height) * 0.4

	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: %s; stroke-width: 0.7' "+
		"width='%d' height='%d'>", fill, width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, aerr := corner(width, height, i+1, j)
			bx, by, berr := corner(width, height, i, j)
			cx, cy, cerr := corner(width, height, i, j+1)
			dx, dy, derr := corner(width, height, i+1, j+1)

			if aerr != nil || berr != nil || cerr != nil || derr != nil {
				log.Printf("Tried to write invalid polygon value")
				continue
			}

			r, b := color(i, j)
			fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill=\"rgb(%d, 0, %d)\"/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, r, b)
		}
	}
	fmt.Fprint(out, "</svg>")
}

func corner(width, height, i, j int) (float64, float64, error) {
	x, y, z := coords(i, j)

	if math.IsInf(z, 0) || math.IsNaN(z) {
		return 0, 0, fmt.Errorf("f(%g, %g) corner has invalid values of i:%d j:%d", x, y, i, j)
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, nil
}

func coords(i, j int) (float64, float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := float64(xyrange) * (float64(i)/float64(cells) - 0.5)
	y := float64(xyrange) * (float64(j)/float64(cells) - 0.5)

	// Compute surface height z.
	var z float64
	switch *figure {
	case "eggbox":
		z = eggbox(x, y)
	case "moguls":
		z = moguls(x, y)
	case "saddle":
		z = saddle(x, y)
	default:
		z = hypot(x, y)
	}

	return x, y, z
}

func color(i, j int) (int, int) {
	_, _, z := coords(i, j)

	if z > 0 {
		return 255, 0
	} else if z < 0 {
		return 0, 255
	} else {
		return 0, 0
	}
}

func hypot(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func eggbox(x, y float64) float64 {
	return (math.Sin(x) * math.Sin(y)) * 0.25
}

func moguls(x, y float64) float64 {
	return math.Pow(2, math.Sin(x)) * math.Pow(2, math.Sin(y)) / 25
}

func saddle(x, y float64) float64 {
	return (3*math.Pow(y, 2) - 2*math.Pow(x, 2)) / 1000
}

//!-

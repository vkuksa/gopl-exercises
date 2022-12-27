// Exercises 3.5, 3.6, 3.7, 3.8, 3.9 from gopl.io

package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/big"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
)

const (
	cx, cy, czoom = 2.0, 2.0, 1.0
	width, height = 1024, 1024
)

var web = flag.Bool("web", false, "start program as webserver")

func main() {
	flag.Parse()

	if *web {
		http.HandleFunc("/", handler)
		http.ListenAndServe("localhost:8000", nil)
	}

	draw(os.Stdout, cx, cy, czoom)
}

func handler(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.ParseFloat(r.URL.Query().Get("x"), 64)
	if err != nil {
		x = cx
	}

	y, err := strconv.ParseFloat(r.URL.Query().Get("y"), 64)
	if err != nil {
		y = cy
	}

	zoom, err := strconv.ParseFloat(r.URL.Query().Get("zoom"), 64)
	if err != nil {
		zoom = czoom
	}

	draw(w, x, y, zoom)
}

func draw(out io.Writer, x, y, zoom float64) {
	img := generatepicture(width*2, height*2, x, y, zoom)

	supersampled := supersample(img, width, height)

	png.Encode(out, supersampled) // NOTE: ignoring errors
}

func generatepicture(width, height int, ix, iy, zoom float64) *image.RGBA {

	s := math.Abs(zoom)
	v := float64(2) / s
	if math.IsNaN(v) {
		v = 2
	}

	var xmin, ymin, xmax, ymax float64 = -v, -v, +v, +v

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin
			z := complex(x, y)

			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot128(z))
		}
	}
	return img
}

func mandelbrot128(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			y := 255 - contrast*n
			cb := uint8(real(v) * 128)
			cr := uint8(imag(v) * 128)

			return color.YCbCr{y, cb, cr}

		}
	}
	return color.Black
}

func mandelbrot64(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex64
	for n := uint8(0); n < iterations; n++ {
		v = v*v + complex64(z)
		if cmplx.Abs(complex128(v)) > 2 {
			y := 255 - contrast*n
			cb := uint8(real(v) * 128)
			cr := uint8(imag(v) * 128)

			return color.YCbCr{y, cb, cr}

		}
	}
	return color.Black
}

func mandelbrotbigfloat(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	var bigmultiplier = &big.Float{}
	bigmultiplier.SetInt64(128)

	vr, vi, zr, zi := &big.Float{}, &big.Float{}, &big.Float{}, &big.Float{}
	zr.SetFloat64(real(z))
	zi.SetFloat64(imag(z))
	for n := uint8(0); n < iterations; n++ {
		// v = v*v + z
		vr2, vi2 := &big.Float{}, &big.Float{}
		vr2.Mul(vr, vr).Sub(vr2, (&big.Float{}).Mul(vi, vi)).Add(vr2, zr)
		vi2.Mul(vr, vi).Mul(vi2, big.NewFloat(2)).Add(vi2, zi)
		vr, vi = vr2, vi2

		vrf, _ := vr.Float64()
		vif, _ := vi.Float64()
		if math.Hypot(vrf, vif) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}

func supersample(img *image.RGBA, width, height int) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			result.Set(x, y, average(img.At(x*2, y*2), img.At((x*2)+1, y*2), img.At(x*2, (y*2)+1), img.At((x*2)+1, (y*2)+1)))
		}
	}
	return result
}

func average(x color.Color, y color.Color, z color.Color, t color.Color) color.Color {
	xr, xg, xb, xa := x.RGBA()
	yr, yg, yb, ya := y.RGBA()
	zr, zg, zb, za := z.RGBA()
	tr, tg, tb, ta := t.RGBA()
	result := color.RGBA{
		uint8((xr + yr + zr + tr) / 4),
		uint8((xg + yg + zg + tg) / 4),
		uint8((xb + yb + zb + tb) / 4),
		uint8((xa + ya + za + ta) / 4)}
	return result
}

func newton(z complex128) color.Color {
	const iterations = 37
	const contrast = 7
	for i := uint8(0); i < iterations; i++ {
		z -= (z - 1/(z*z*z)) / 4
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			return color.Gray{255 - contrast*i}
		}
	}
	return color.Black
}

//!-

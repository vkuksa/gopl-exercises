// Exercise 1.5, 1.6 from gopl.io

package lissajous

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
)

//!-main
// Packages not needed by version in book.

//!+main

var palette = []color.Color{color.Black, color.RGBA{0x05, 0x05, 0x05, 0xff}, color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff}}

const (
	cycles  = 5     // number of complete x oscillator revolutions
	res     = 0.001 // angular resolution
	nframes = 64    // number of animation frames
	delay   = 8     // delay between frames in 10ms units
)

func Lissajous(out io.Writer, size int) {

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5),
				uint8(i)%5)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

const (
	imgPath = "image"
)

func main() {
	drawPNG(imgPath + ".png")
	drawSVG(imgPath + ".svg")
}

func drawPNG(imgPath string) {
	img := image.NewRGBA(image.Rect(0, 0, 600, 400))

	for x := 0; x <= img.Bounds().Dx(); x++ {
		for y := 0; y <= img.Bounds().Dy(); y++ {
			img.Set(x, y, color.Black)
			if x == y {
				img.Set(x, y, color.RGBA{R: 255, A: 255})
			}
		}
	}

	r := image.Rect(10, 10, 30, 390)
	c := color.RGBA{B: 255, A: 255}
	draw.Draw(img, r, &image.Uniform{c}, image.ZP, draw.Src)

	drawCircle := func(x, y, r int, drawFn func(int, int)) {
		for theta := 0.0; theta < 2*math.Pi; theta += 0.001 {
			xi := x + int(float64(r)*math.Cos(theta))
			yi := y - int(float64(r)*math.Sin(theta))
			drawFn(xi, yi)
		}
	}

	drawCircle(450, 150, 50, func(x, y int) {
		img.Set(x, y, color.RGBA{G: 255, A: 255})
	})

	f, err := os.Create(imgPath)
	if err != nil {
		log.Fatalf("Failed to create: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}
}

func drawSVG(imgPath string) {
	f, err := os.Create(imgPath)
	if err != nil {
		log.Fatalf("Failed to create: %v", err)
	}
	defer f.Close()

	img := svg.New(f)
	img.Start(600, 400)
	img.Rect(0, 0, 600, 400)
	img.Line(0, 0, 600, 600, "stroke:red;stroke-width:1")
	img.Rect(10, 10, 30, 390, "fill:blue")
	img.Circle(450, 150, 50, "stroke:green;stroke-width:1")
	img.End()
}

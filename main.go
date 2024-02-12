package main

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	"sync"
)

const (
	MAX_ITER = 100
	HEIGHT   = 800
	WIDTH    = 1024
)

var (
	XStep = 4.0 / WIDTH
	YStep = 4.0 / HEIGHT
)

func Julia(x, y, mouseX, mouseY float64) color.RGBA {
	z := complex(x, y)
	c := complex(mouseX, mouseY)
	for iter := 0; iter < MAX_ITER; iter++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4.0 {
			return colornames.Black
		}
		z = z*z + c
	}
	return colornames.Aquamarine
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Julia!",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	canvas := pixel.MakePictureData(pixel.R(0, 0, WIDTH, HEIGHT))
	var wg sync.WaitGroup
	var mouseX, mouseY float64

	for !win.Closed() {
		mouseX = -2.0 + XStep*win.MousePosition().X
		mouseY = -2.0 + YStep*win.MousePosition().Y

		wg.Add(HEIGHT)
		for y := 0; y < HEIGHT; y++ {
			go func(y int) {
				normX := -2.0
				normY := 2.0 - YStep*float64(y)
				for i := 0; i < WIDTH; i++ {
					normX += XStep
					canvas.Pix[y*WIDTH+i] = Julia(normX, normY, mouseX, mouseY)
				}
				wg.Done()
			}(y)
		}
		wg.Wait()

		sprite := pixel.NewSprite(canvas, canvas.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

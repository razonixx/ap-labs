package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type food struct {
	vec pixel.Vec
}

func run() {
	foods := []food{
		{vec: pixel.V(0, 0)},
		{vec: pixel.V(100, 100)},
		{vec: pixel.V(200, 200)},
		{vec: pixel.V(300, 300)},
		{vec: pixel.V(400, 400)},
		{vec: pixel.V(500, 500)},
	}
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	imd.Color = colornames.Yellow
	for _, singleFood := range foods {
		imd.Push(singleFood.vec)
	}
	imd.Circle(8, 0)

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

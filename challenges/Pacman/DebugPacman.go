package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type food struct {
	vec pixel.Vec
}

type wall struct {
	vec1, vec2 pixel.Vec
	line       pixel.Line
}

const constSpeed = 3.3

var speed = constSpeed
var pause = false
var score = 0
var imd = imdraw.New(nil)

var foods = make(map[int]food)
var walls = make(map[int]wall)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func checkCollision(vec1, vec2 pixel.Vec) bool { //Check if vec1 is in boundary of vec2
	if vec1.X > vec2.X-15 && vec1.X < vec2.X+15 {
		if vec1.Y > vec2.Y-15 && vec1.Y < vec2.Y+15 {
			return true
		}
	}
	return false
}

func checkLine(vec1 pixel.Vec, line pixel.Line) bool { //Check if vec1 is in boundary of rect created by vec2 and vec3
	if vec1.X > line.A.X && vec1.X < line.B.X {
		if vec1.Y >= line.A.Y-5 && vec1.Y <= line.B.Y+5 {
			log.Println("Collition")
			return true
		}
	}
	return false
}

func setUpLevel() {
	{ //Declare food
		foods[0] = food{vec: pixel.V(float64(100), float64(45))}
		foods[1] = food{vec: pixel.V(float64(200), float64(45))}
		foods[2] = food{vec: pixel.V(float64(300), float64(45))}
		foods[3] = food{vec: pixel.V(float64(400), float64(45))}
		foods[4] = food{vec: pixel.V(float64(500), float64(45))}
		foods[5] = food{vec: pixel.V(float64(600), float64(45))}
		foods[6] = food{vec: pixel.V(float64(700), float64(45))}
		foods[7] = food{vec: pixel.V(float64(800), float64(45))}
		foods[8] = food{vec: pixel.V(float64(900), float64(45))}
		foods[9] = food{vec: pixel.V(float64(1000), float64(45))}
	}

	{ //Declare walls
		walls[0] = wall{vec1: pixel.V(float64(112), float64(635)), vec2: pixel.V(float64(210), float64(635))}
		walls[1] = wall{vec1: pixel.V(float64(112), float64(590)), vec2: pixel.V(float64(210), float64(590))}
	}
	imd.Color = colornames.Yellow
	for _, singleFood := range foods {
		imd.Push(singleFood.vec)
	}
	imd.Circle(6, 0)
}

func run() {

	setUpLevel()

	cfg := pixelgl.WindowConfig{
		Title:  "Pacman",
		Bounds: pixel.R(0, 0, 1024, 688),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pacmanDown, err := loadPicture("PacmanDown.png")
	if err != nil {
		panic(err)
	}
	pacmanUp, err := loadPicture("PacmanUp.png")
	if err != nil {
		panic(err)
	}
	pacmanLeft, err := loadPicture("PacmanLeft.png")
	if err != nil {
		panic(err)
	}
	pacmanRight, err := loadPicture("PacmanRight.png")
	if err != nil {
		panic(err)
	}
	background, err := loadPicture("pacmanBackground.png")
	if err != nil {
		panic(err)
	}

	pacmanSpriteUp := pixel.NewSprite(pacmanUp, pacmanUp.Bounds())
	pacmanSpriteDown := pixel.NewSprite(pacmanDown, pacmanDown.Bounds())
	pacmanSpriteLeft := pixel.NewSprite(pacmanLeft, pacmanLeft.Bounds())
	pacmanSpriteRight := pixel.NewSprite(pacmanRight, pacmanRight.Bounds())
	backgroundSprite := pixel.NewSprite(background, background.Bounds())

	win.Clear(colornames.Greenyellow)
	ctrl := 1
	last := time.Now()
	currentPos := pixel.V(350, 45) //Starting position

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Greenyellow)
		mat := pixel.IM
		backgroundMat := pixel.IM
		backgroundMat = mat.Moved(pixel.V(512, 344))
		backgroundSprite.Draw(win, backgroundMat)
		for _, singeWall := range walls {
			singeWall.line = pixel.L(singeWall.vec1, singeWall.vec2)
			if checkLine(currentPos, singeWall.line) {
				pause = true
			}
		}
		for i, singleFood := range foods { //Draw and check collision with food
			if checkCollision(currentPos, singleFood.vec) {
				pause = true
				imd.Color = colornames.Black
				imd.Push(singleFood.vec)
				imd.Circle(6, 0)
				delete(foods, i)
				score += 100
				fmt.Printf("Score: %d\n", score)
			}
		}
		mat = mat.ScaledXY(pixel.ZV, pixel.V(0.55, 0.55))
		mat = mat.Rotated(pixel.ZV, (float64(ctrl) * math.Round(dt*100) / 100))
		imd.Draw(win)

		if win.Pressed(pixelgl.KeyLeft) {
			pause = false
			currentPos.X -= speed
			mat = mat.Moved(currentPos)
			pacmanSpriteLeft.Draw(win, mat)
			if checkCollision(currentPos, pixel.V(10, 370)) {
				currentPos.X = 1024
			}
			log.Printf("X: %d, Y: %d\n", int64(currentPos.X), int64(currentPos.Y)) //Print current position
		} else if win.Pressed(pixelgl.KeyRight) {
			pause = false
			currentPos.X += speed
			mat = mat.Moved(currentPos)
			pacmanSpriteRight.Draw(win, mat)
			if checkCollision(currentPos, pixel.V(1020, 370)) {
				currentPos.X = 0
			}
			log.Printf("X: %d, Y: %d\n", int64(currentPos.X), int64(currentPos.Y)) //Print current position
		} else if win.Pressed(pixelgl.KeyDown) {
			pause = false
			currentPos.Y -= speed
			mat = mat.Moved(currentPos)
			pacmanSpriteDown.Draw(win, mat)
			log.Printf("X: %d, Y: %d\n", int64(currentPos.X), int64(currentPos.Y)) //Print current position
		} else if win.Pressed(pixelgl.KeyUp) {
			pause = false
			currentPos.Y += speed
			mat = mat.Moved(currentPos)
			pacmanSpriteUp.Draw(win, mat)
			log.Printf("X: %d, Y: %d\n", int64(currentPos.X), int64(currentPos.Y)) //Print current position
		} else {
			mat = mat.Moved(currentPos)
			pacmanSpriteUp.Draw(win, mat)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
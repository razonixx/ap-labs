package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type food struct {
	vec pixel.Vec
}

type wall struct {
	rect pixel.Rect
}

const constSpeed = 1.5

var imd = imdraw.New(nil)
var speed = constSpeed
var pause = false
var score = 0
var foods = make(map[int]food)
var walls = []wall{}

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

func checkCollisionWall(vec pixel.Vec, rect pixel.Rect) bool { //Check if vec1 is in boundary of vec2
	newRect := pixel.R(vec.X-15, vec.Y-15, vec.X+15, vec.Y+15)
	if newRect.Min.X < rect.Min.X+rect.W() &&
		newRect.Min.X+newRect.W() > rect.Min.X &&
		newRect.Min.Y < rect.Min.Y+rect.H() &&
		newRect.Min.Y+newRect.H() > rect.Min.Y {
		// collision detected!
		return true
	}
	return false
}

func checkCollisionAllWalls(vec pixel.Vec) bool {
	for _, singeWall := range walls {
		if checkCollisionWall(vec, singeWall.rect) {
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
	imd.Color = colornames.Yellow
	for _, singleFood := range foods {
		imd.Push(singleFood.vec)
	}
	imd.Circle(6, 0)
	imd.Color = colornames.Blue
	limits := []wall{
		{rect: pixel.R(0, 0, 1024, 720)}, //Border
	}
	for _, l := range limits {
		imd.Push(l.rect.Min, l.rect.Max)
	}
	imd.Rectangle(20)
	imd.Color = colornames.Black
	exits := []wall{
		{rect: pixel.R(0, 360-30, 20, 360+30)},      //Left Exit
		{rect: pixel.R(1000, 360-30, 1024, 360+30)}, //Right Exit
	}
	for _, e := range exits {
		imd.Push(e.rect.Min, e.rect.Max)
	}
	imd.Rectangle(0)
	imd.Color = colornames.Blue
	walls = []wall{
		{rect: pixel.R(70, 70, 300, 300)},                   //SouthWest block
		{rect: pixel.R(1024-300, 70, 1024-70, 300)},         //SouthEast block
		{rect: pixel.R(70, 720-300, 300, 720-70)},           //NorthWest block
		{rect: pixel.R(1024-300, 720-300, 1024-70, 720-70)}, //NorthEast block

		{rect: pixel.R(512-150, 360-100, 512+150, 360+100)}, //Center block

		{rect: pixel.R(512-50, 360+200-50, 512+50, 360+300+50)}, //Center North block
		{rect: pixel.R(512-50, 360-300-50, 512+50, 360-200+50)}, //Center South block
		{rect: pixel.R(512+325-25, 360-25, 512+325+25, 360+25)}, //Center East block
		{rect: pixel.R(512-325-25, 360-25, 512-325+25, 360+25)}, //Center West block
	}
	for _, w := range walls {
		imd.Push(w.rect.Min, w.rect.Max)
		imd.Rectangle(0)
	}

}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pacman",
		Bounds: pixel.R(0, 0, 1024, 720),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	setUpLevel()

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

	pacmanSpriteUp := pixel.NewSprite(pacmanUp, pacmanUp.Bounds())
	pacmanSpriteDown := pixel.NewSprite(pacmanDown, pacmanDown.Bounds())
	pacmanSpriteLeft := pixel.NewSprite(pacmanLeft, pacmanLeft.Bounds())
	pacmanSpriteRight := pixel.NewSprite(pacmanRight, pacmanRight.Bounds())

	currentPos := pixel.V(350, 45) //Starting position

	shouldMoveLeft := true
	shouldMoveRight := false
	shouldMoveUp := false
	shouldMoveDown := false

	hasWallLeft := false
	hasWallRight := false
	hasWallUp := false
	hasWallDown := false

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Black)
		mat := pixel.IM

		for i, singleFood := range foods { //Draw and check collision with food
			if checkCollision(currentPos, singleFood.vec) {
				imd.Color = colornames.Black
				imd.Push(singleFood.vec)
				imd.Circle(6, 0)
				delete(foods, i)
				score += 100
				fmt.Printf("Score: %d Remaining: %d\n", score, len(foods))
				if len(foods) == 0 {
					log.Println("YOU WIN!!")
					os.Exit(0)
				}
			}
		}
		mat = mat.ScaledXY(pixel.ZV, pixel.V(0.55, 0.55))
		mat = mat.Rotated(pixel.ZV, (float64(1) * math.Round(dt*100) / 100))
		imd.Draw(win)

		if shouldMoveLeft {
			if !pause && currentPos.X > 25 && !hasWallLeft {
				currentPos.X -= speed
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteLeft.Draw(win, mat)
			if checkCollisionWall(currentPos, pixel.R(0, 360-30, 20, 360+30)) {
				currentPos.X = 1024
			}
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallRight && !hasWallUp && !hasWallDown {
					hasWallLeft = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallUp = false
				hasWallRight = false
				hasWallDown = false
			}
		} else if shouldMoveRight {
			if !pause && currentPos.X < 1000 && !hasWallRight {
				currentPos.X += speed
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteRight.Draw(win, mat)
			if checkCollisionWall(currentPos, pixel.R(1000, 360-30, 1024, 360+30)) {
				currentPos.X = 0
			}
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallLeft && !hasWallUp && !hasWallDown {
					hasWallRight = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallUp = false
				hasWallDown = false
			}
		} else if shouldMoveUp {
			if !pause && currentPos.Y < 695 && !hasWallUp {
				currentPos.Y += speed
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteUp.Draw(win, mat)
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallRight && !hasWallLeft && !hasWallDown {
					hasWallUp = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallRight = false
				hasWallDown = false
			}
		} else if shouldMoveDown {
			if !pause && currentPos.Y > 25 && !hasWallDown {
				currentPos.Y -= speed
			}
			mat = mat.Moved(currentPos)
			pacmanSpriteDown.Draw(win, mat)
			for _, singeWall := range walls {
				if checkCollisionWall(currentPos, singeWall.rect) && !hasWallRight && !hasWallUp && !hasWallLeft {
					hasWallDown = true
				}
			}
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallRight = false
				hasWallUp = false
			}
		}

		if win.Pressed(pixelgl.KeyLeft) {
			if pause {
				currentPos.X -= speed
			}
			pause = false
			shouldMoveLeft = true
			shouldMoveRight = false
			shouldMoveUp = false
			shouldMoveDown = false
			hasWallRight = false
		} else if win.Pressed(pixelgl.KeyRight) {
			if pause {
				currentPos.X += speed
			}
			pause = false
			shouldMoveLeft = false
			shouldMoveRight = true
			shouldMoveUp = false
			shouldMoveDown = false
			hasWallLeft = false
		} else if win.Pressed(pixelgl.KeyDown) {
			if pause {
				currentPos.Y -= speed
			}
			pause = false
			shouldMoveLeft = false
			shouldMoveRight = false
			shouldMoveUp = false
			shouldMoveDown = true
			hasWallUp = false
		} else if win.Pressed(pixelgl.KeyUp) {
			if pause {
				currentPos.Y += speed
			}
			pause = false
			shouldMoveLeft = false
			shouldMoveRight = false
			shouldMoveUp = true
			shouldMoveDown = false
			hasWallDown = false
			if !checkCollisionAllWalls(currentPos) {
				hasWallLeft = false
				hasWallRight = false
			}

		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

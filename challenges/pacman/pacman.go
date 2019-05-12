package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type node struct {
	vec pixel.Vec
}

type wall struct {
	rect pixel.Rect
}

const constSpeed = 3.5

var imd = imdraw.New(nil)
var speed = constSpeed
var pause = false
var score = 0
var numGhosts int
var currentPos pixel.Vec
var nodes = make(map[int]node)
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
	if vec1.X > vec2.X-20 && vec1.X < vec2.X+20 {
		if vec1.Y > vec2.Y-20 && vec1.Y < vec2.Y+20 {
			return true
		}
	}
	return false
}

func checkCollisionWall(vec pixel.Vec, rect pixel.Rect) bool { //Check if vec1 is colliding with the rectangle using bounding box algorithm
	newRect := pixel.R(vec.X-15, vec.Y-15, vec.X+15, vec.Y+15)
	if newRect.Min.X < rect.Min.X+rect.W() && newRect.Min.X+newRect.W() > rect.Min.X && newRect.Min.Y < rect.Min.Y+rect.H() && newRect.Min.Y+newRect.H() > rect.Min.Y {
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

func checkCollisionGhost(vec pixel.Vec) bool {
	if checkCollision(vec, currentPos) {
		log.Printf("YOU LOSE!!")
		os.Exit(0)
	}
	for _, singeWall := range walls {
		if checkCollisionWall(vec, singeWall.rect) {
			return true
		}
	}
	return false
}

func setUpLevel() {
	{
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
			{rect: pixel.R(-10, 360, 0, 360+30)},        //Left Exit
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

		//Declare node
		{
			{ //row 1
				nodes[0] = node{vec: pixel.V(float64(40), float64(40))}
				nodes[1] = node{vec: pixel.V(float64(80), float64(40))}
				nodes[2] = node{vec: pixel.V(float64(120), float64(40))}
				nodes[3] = node{vec: pixel.V(float64(160), float64(40))}
				nodes[4] = node{vec: pixel.V(float64(200), float64(40))}
				nodes[5] = node{vec: pixel.V(float64(240), float64(40))}
				nodes[6] = node{vec: pixel.V(float64(280), float64(40))}
				nodes[7] = node{vec: pixel.V(float64(320), float64(40))}
				nodes[8] = node{vec: pixel.V(float64(360), float64(40))}
				nodes[9] = node{vec: pixel.V(float64(400), float64(40))}
				nodes[10] = node{vec: pixel.V(float64(440), float64(40))}
				nodes[11] = node{vec: pixel.V(float64(580), float64(40))}
				nodes[12] = node{vec: pixel.V(float64(620), float64(40))}
				nodes[13] = node{vec: pixel.V(float64(660), float64(40))}
				nodes[14] = node{vec: pixel.V(float64(700), float64(40))}
				nodes[15] = node{vec: pixel.V(float64(740), float64(40))}
				nodes[16] = node{vec: pixel.V(float64(780), float64(40))}
				nodes[17] = node{vec: pixel.V(float64(820), float64(40))}
				nodes[18] = node{vec: pixel.V(float64(860), float64(40))}
				nodes[19] = node{vec: pixel.V(float64(900), float64(40))}
				nodes[20] = node{vec: pixel.V(float64(940), float64(40))}
				nodes[21] = node{vec: pixel.V(float64(980), float64(40))}
			}
			{ //row2
				nodes[22] = node{vec: pixel.V(float64(40), float64(80))}
				nodes[23] = node{vec: pixel.V(float64(320), float64(80))}
				nodes[24] = node{vec: pixel.V(float64(360), float64(80))}
				nodes[25] = node{vec: pixel.V(float64(400), float64(80))}
				nodes[26] = node{vec: pixel.V(float64(440), float64(80))}
				nodes[27] = node{vec: pixel.V(float64(580), float64(80))}
				nodes[28] = node{vec: pixel.V(float64(620), float64(80))}
				nodes[29] = node{vec: pixel.V(float64(660), float64(80))}
				nodes[30] = node{vec: pixel.V(float64(700), float64(80))}
				nodes[31] = node{vec: pixel.V(float64(980), float64(80))}
			}
			{ //row3
				nodes[32] = node{vec: pixel.V(float64(40), float64(120))}
				nodes[33] = node{vec: pixel.V(float64(320), float64(120))}
				nodes[34] = node{vec: pixel.V(float64(360), float64(120))}
				nodes[35] = node{vec: pixel.V(float64(400), float64(120))}
				nodes[36] = node{vec: pixel.V(float64(440), float64(120))}
				nodes[37] = node{vec: pixel.V(float64(580), float64(120))}
				nodes[38] = node{vec: pixel.V(float64(620), float64(120))}
				nodes[39] = node{vec: pixel.V(float64(660), float64(120))}
				nodes[40] = node{vec: pixel.V(float64(700), float64(120))}
				nodes[41] = node{vec: pixel.V(float64(980), float64(120))}
			}
			{ //row4
				nodes[42] = node{vec: pixel.V(float64(40), float64(160))}
				nodes[43] = node{vec: pixel.V(float64(320), float64(160))}
				nodes[44] = node{vec: pixel.V(float64(360), float64(160))}
				nodes[45] = node{vec: pixel.V(float64(400), float64(160))}
				nodes[46] = node{vec: pixel.V(float64(440), float64(160))}
				nodes[47] = node{vec: pixel.V(float64(580), float64(160))}
				nodes[48] = node{vec: pixel.V(float64(620), float64(160))}
				nodes[49] = node{vec: pixel.V(float64(660), float64(160))}
				nodes[50] = node{vec: pixel.V(float64(700), float64(160))}
				nodes[51] = node{vec: pixel.V(float64(980), float64(160))}
			}
			{ //row5
				nodes[52] = node{vec: pixel.V(float64(40), float64(200))}
				nodes[53] = node{vec: pixel.V(float64(320), float64(200))}
				nodes[54] = node{vec: pixel.V(float64(360), float64(200))}
				nodes[55] = node{vec: pixel.V(float64(400), float64(200))}
				nodes[56] = node{vec: pixel.V(float64(440), float64(200))}
				nodes[57] = node{vec: pixel.V(float64(580), float64(200))}
				nodes[58] = node{vec: pixel.V(float64(620), float64(200))}
				nodes[59] = node{vec: pixel.V(float64(660), float64(200))}
				nodes[60] = node{vec: pixel.V(float64(700), float64(200))}
				nodes[61] = node{vec: pixel.V(float64(980), float64(200))}
			}
			{ //row6
				nodes[62] = node{vec: pixel.V(float64(40), float64(240))}
				nodes[63] = node{vec: pixel.V(float64(320), float64(240))}
				nodes[64] = node{vec: pixel.V(float64(360), float64(240))}
				nodes[65] = node{vec: pixel.V(float64(400), float64(240))}
				nodes[66] = node{vec: pixel.V(float64(440), float64(240))}
				nodes[72] = node{vec: pixel.V(float64(485), float64(240))}
				nodes[73] = node{vec: pixel.V(float64(535), float64(240))}
				nodes[67] = node{vec: pixel.V(float64(580), float64(240))}
				nodes[68] = node{vec: pixel.V(float64(620), float64(240))}
				nodes[69] = node{vec: pixel.V(float64(660), float64(240))}
				nodes[70] = node{vec: pixel.V(float64(700), float64(240))}
				nodes[71] = node{vec: pixel.V(float64(980), float64(240))}
			}
			{ //row7
				nodes[74] = node{vec: pixel.V(float64(40), float64(280))}
				nodes[75] = node{vec: pixel.V(float64(320), float64(280))}
				nodes[76] = node{vec: pixel.V(float64(700), float64(280))}
				nodes[77] = node{vec: pixel.V(float64(980), float64(280))}
			}
			{ //row8
				nodes[78] = node{vec: pixel.V(float64(40), float64(320))}
				nodes[79] = node{vec: pixel.V(float64(80), float64(320))}
				nodes[80] = node{vec: pixel.V(float64(120), float64(320))}
				nodes[81] = node{vec: pixel.V(float64(160), float64(320))}
				nodes[82] = node{vec: pixel.V(float64(200), float64(320))}
				nodes[83] = node{vec: pixel.V(float64(240), float64(320))}
				nodes[84] = node{vec: pixel.V(float64(280), float64(320))}
				nodes[86] = node{vec: pixel.V(float64(320), float64(320))}
				nodes[87] = node{vec: pixel.V(float64(700), float64(320))}
				nodes[88] = node{vec: pixel.V(float64(740), float64(320))}
				nodes[89] = node{vec: pixel.V(float64(780), float64(320))}
				nodes[90] = node{vec: pixel.V(float64(820), float64(320))}
				nodes[91] = node{vec: pixel.V(float64(860), float64(320))}
				nodes[92] = node{vec: pixel.V(float64(900), float64(320))}
				nodes[93] = node{vec: pixel.V(float64(940), float64(320))}
				nodes[94] = node{vec: pixel.V(float64(980), float64(320))}
			}
			{ //row9
				nodes[95] = node{vec: pixel.V(float64(40), float64(360))}
				nodes[96] = node{vec: pixel.V(float64(80), float64(360))}
				nodes[97] = node{vec: pixel.V(float64(120), float64(360))}
				nodes[98] = node{vec: pixel.V(float64(240), float64(360))}
				nodes[99] = node{vec: pixel.V(float64(280), float64(360))}
				nodes[194] = node{vec: pixel.V(float64(320), float64(360))}
				nodes[195] = node{vec: pixel.V(float64(700), float64(360))}
				nodes[196] = node{vec: pixel.V(float64(740), float64(360))}
				nodes[197] = node{vec: pixel.V(float64(780), float64(360))}
				nodes[198] = node{vec: pixel.V(float64(900), float64(360))}
				nodes[199] = node{vec: pixel.V(float64(940), float64(360))}
				nodes[200] = node{vec: pixel.V(float64(980), float64(360))}
			}
			{ //row10
				nodes[100] = node{vec: pixel.V(float64(40), float64(400))}
				nodes[101] = node{vec: pixel.V(float64(80), float64(400))}
				nodes[102] = node{vec: pixel.V(float64(120), float64(400))}
				nodes[103] = node{vec: pixel.V(float64(160), float64(400))}
				nodes[104] = node{vec: pixel.V(float64(200), float64(400))}
				nodes[105] = node{vec: pixel.V(float64(240), float64(400))}
				nodes[106] = node{vec: pixel.V(float64(280), float64(400))}
				nodes[107] = node{vec: pixel.V(float64(320), float64(400))}
				nodes[108] = node{vec: pixel.V(float64(700), float64(400))}
				nodes[109] = node{vec: pixel.V(float64(740), float64(400))}
				nodes[110] = node{vec: pixel.V(float64(780), float64(400))}
				nodes[111] = node{vec: pixel.V(float64(820), float64(400))}
				nodes[112] = node{vec: pixel.V(float64(860), float64(400))}
				nodes[113] = node{vec: pixel.V(float64(900), float64(400))}
				nodes[114] = node{vec: pixel.V(float64(940), float64(400))}
				nodes[115] = node{vec: pixel.V(float64(980), float64(400))}
			}
			{ //row11
				nodes[116] = node{vec: pixel.V(float64(40), float64(440))}
				nodes[117] = node{vec: pixel.V(float64(320), float64(440))}
				nodes[118] = node{vec: pixel.V(float64(700), float64(440))}
				nodes[119] = node{vec: pixel.V(float64(980), float64(440))}
			}
			{ //row12
				nodes[120] = node{vec: pixel.V(float64(40), float64(480))}
				nodes[121] = node{vec: pixel.V(float64(320), float64(480))}
				nodes[122] = node{vec: pixel.V(float64(360), float64(480))}
				nodes[123] = node{vec: pixel.V(float64(400), float64(480))}
				nodes[124] = node{vec: pixel.V(float64(440), float64(480))}
				nodes[125] = node{vec: pixel.V(float64(485), float64(480))}
				nodes[126] = node{vec: pixel.V(float64(535), float64(480))}
				nodes[127] = node{vec: pixel.V(float64(580), float64(480))}
				nodes[128] = node{vec: pixel.V(float64(620), float64(480))}
				nodes[129] = node{vec: pixel.V(float64(660), float64(480))}
				nodes[130] = node{vec: pixel.V(float64(700), float64(480))}
				nodes[131] = node{vec: pixel.V(float64(980), float64(480))}
			}
			{ //row13
				nodes[132] = node{vec: pixel.V(float64(40), float64(520))}
				nodes[133] = node{vec: pixel.V(float64(320), float64(520))}
				nodes[134] = node{vec: pixel.V(float64(360), float64(520))}
				nodes[135] = node{vec: pixel.V(float64(400), float64(520))}
				nodes[136] = node{vec: pixel.V(float64(440), float64(520))}
				nodes[137] = node{vec: pixel.V(float64(580), float64(520))}
				nodes[138] = node{vec: pixel.V(float64(620), float64(520))}
				nodes[139] = node{vec: pixel.V(float64(660), float64(520))}
				nodes[140] = node{vec: pixel.V(float64(700), float64(520))}
				nodes[141] = node{vec: pixel.V(float64(980), float64(520))}
			}
			{ //row14
				nodes[142] = node{vec: pixel.V(float64(40), float64(560))}
				nodes[143] = node{vec: pixel.V(float64(320), float64(560))}
				nodes[144] = node{vec: pixel.V(float64(360), float64(560))}
				nodes[145] = node{vec: pixel.V(float64(400), float64(560))}
				nodes[146] = node{vec: pixel.V(float64(440), float64(560))}
				nodes[147] = node{vec: pixel.V(float64(580), float64(560))}
				nodes[148] = node{vec: pixel.V(float64(620), float64(560))}
				nodes[149] = node{vec: pixel.V(float64(660), float64(560))}
				nodes[150] = node{vec: pixel.V(float64(700), float64(560))}
				nodes[151] = node{vec: pixel.V(float64(980), float64(560))}
			}
			{ //row15
				nodes[152] = node{vec: pixel.V(float64(40), float64(600))}
				nodes[153] = node{vec: pixel.V(float64(320), float64(600))}
				nodes[154] = node{vec: pixel.V(float64(360), float64(600))}
				nodes[155] = node{vec: pixel.V(float64(400), float64(600))}
				nodes[156] = node{vec: pixel.V(float64(440), float64(600))}
				nodes[157] = node{vec: pixel.V(float64(580), float64(600))}
				nodes[158] = node{vec: pixel.V(float64(620), float64(600))}
				nodes[159] = node{vec: pixel.V(float64(660), float64(600))}
				nodes[160] = node{vec: pixel.V(float64(700), float64(600))}
				nodes[161] = node{vec: pixel.V(float64(980), float64(600))}
			}
			{ //row16
				nodes[162] = node{vec: pixel.V(float64(40), float64(640))}
				nodes[163] = node{vec: pixel.V(float64(320), float64(640))}
				nodes[164] = node{vec: pixel.V(float64(360), float64(640))}
				nodes[165] = node{vec: pixel.V(float64(400), float64(640))}
				nodes[166] = node{vec: pixel.V(float64(440), float64(640))}
				nodes[167] = node{vec: pixel.V(float64(580), float64(640))}
				nodes[168] = node{vec: pixel.V(float64(620), float64(640))}
				nodes[169] = node{vec: pixel.V(float64(660), float64(640))}
				nodes[170] = node{vec: pixel.V(float64(700), float64(640))}
				nodes[171] = node{vec: pixel.V(float64(980), float64(640))}
			}
			{ //row 17
				nodes[172] = node{vec: pixel.V(float64(40), float64(680))}
				nodes[173] = node{vec: pixel.V(float64(80), float64(680))}
				nodes[174] = node{vec: pixel.V(float64(120), float64(680))}
				nodes[175] = node{vec: pixel.V(float64(160), float64(680))}
				nodes[176] = node{vec: pixel.V(float64(200), float64(680))}
				nodes[177] = node{vec: pixel.V(float64(240), float64(680))}
				nodes[178] = node{vec: pixel.V(float64(280), float64(680))}
				nodes[179] = node{vec: pixel.V(float64(320), float64(680))}
				nodes[180] = node{vec: pixel.V(float64(360), float64(680))}
				nodes[181] = node{vec: pixel.V(float64(400), float64(680))}
				nodes[182] = node{vec: pixel.V(float64(440), float64(680))}
				nodes[183] = node{vec: pixel.V(float64(580), float64(680))}
				nodes[184] = node{vec: pixel.V(float64(620), float64(680))}
				nodes[185] = node{vec: pixel.V(float64(660), float64(680))}
				nodes[186] = node{vec: pixel.V(float64(700), float64(680))}
				nodes[187] = node{vec: pixel.V(float64(740), float64(680))}
				nodes[188] = node{vec: pixel.V(float64(780), float64(680))}
				nodes[189] = node{vec: pixel.V(float64(820), float64(680))}
				nodes[190] = node{vec: pixel.V(float64(860), float64(680))}
				nodes[191] = node{vec: pixel.V(float64(900), float64(680))}
				nodes[192] = node{vec: pixel.V(float64(940), float64(680))}
				nodes[193] = node{vec: pixel.V(float64(980), float64(680))}
			}
		}
		imd.Color = colornames.Yellow
		for _, singlenode := range nodes {
			imd.Push(singlenode.vec)
		}
		imd.Circle(8, 0)

	}
}

func squareRun(pos pixel.Vec, win *pixelgl.Window) {
	imd.Color = colornames.Red

	timeToInstruction := 2.0
	ghostSpeed := 10.5
	last := time.Now()
	sumTimePassed := timeToInstruction
	var seed int64
	var current pixel.Vec
	for _, c := range "ghost" {
		seed += int64(c)
	}

	left := pixel.V(-.1*ghostSpeed, 0)
	right := pixel.V(.1*ghostSpeed, 0)
	up := pixel.V(0, .1*ghostSpeed)
	down := pixel.V(0, -.1*ghostSpeed)

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		sumTimePassed += dt
		if sumTimePassed >= timeToInstruction {
			random := rand.Int31n(4)
			switch random {
			case 0:
				current = left
			case 1:
				current = right
			case 2:
				current = up
			case 3:
				current = down
			}
			sumTimePassed = 0.0
		}
		if checkCollisionGhost(pos) || pos.X < 10 || pos.X > 1020 || pos.Y < 10 || pos.Y > 710 { //Check collision with all walls and the borders of the map
			fmt.Printf("%s is collisioning\n", "Test")
			current.X = -current.X
			current.Y = -current.Y
			pos.X += current.X
			pos.Y += current.Y
		}
		pos.X += current.X
		pos.Y += current.Y
		square := imdraw.New(nil)
		square.Color = colornames.Red
		rect := pixel.R(pos.X, pos.Y, pos.X+20, pos.Y+20)
		square.Push(rect.Min, rect.Max)
		square.Rectangle(0)
		square.Draw(win)
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

	currentPos = pixel.V(350, 40) //Starting position

	shouldMoveLeft := true
	shouldMoveRight := false
	shouldMoveUp := false
	shouldMoveDown := false

	hasWallLeft := false
	hasWallRight := false
	hasWallUp := false
	hasWallDown := false

	for i := 0; i < numGhosts; i++ {
		//go ghostRun("blinky", blinkySprite, pixel.V(float64(35*i)+25.0, 35), win)
		go squareRun(pixel.V(float64(35*i)+25.0, 35), win)
	}

	lastnodeEaten := 0
	last := time.Now()
	dtCheck := 0.0
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		dtCheck += dt
		win.Clear(colornames.Black)
		if dtCheck >= .0 {

			mat := pixel.IM

			for i, singlenode := range nodes { //Draw and check collision with node
				if checkCollision(currentPos, singlenode.vec) {
					lastnodeEaten = i
					imd.Color = colornames.Black
					imd.Push(singlenode.vec)
					imd.Circle(8, 0)
					delete(nodes, i)
					score += 100
					//fmt.Printf("Score: %d \n", score)
					fmt.Printf("Last node eaten: %d\n", lastnodeEaten)
					if len(nodes) == 0 {
						log.Println("YOU WIN!!")
						os.Exit(0)
					}
				}
			}
			mat = mat.ScaledXY(pixel.ZV, pixel.V(0.80, 0.80))
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
			dtCheck = 0
		}
	}
}

func main() {
	//usage: go run pacman.go -g <numGhosts>
	numGhosts, _ = strconv.Atoi(os.Args[2])
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}

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

type food struct {
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

func checkCollisionGhost(ghost string, vec pixel.Vec) bool {
	if checkCollision(vec, currentPos) {
		log.Printf("YOU LOSE!! Killed by: %s", ghost)
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

		//Declare food
		{
			{ //row 1
				foods[0] = food{vec: pixel.V(float64(40), float64(40))}
				foods[1] = food{vec: pixel.V(float64(80), float64(40))}
				foods[2] = food{vec: pixel.V(float64(120), float64(40))}
				foods[3] = food{vec: pixel.V(float64(160), float64(40))}
				foods[4] = food{vec: pixel.V(float64(200), float64(40))}
				foods[5] = food{vec: pixel.V(float64(240), float64(40))}
				foods[6] = food{vec: pixel.V(float64(280), float64(40))}
				foods[7] = food{vec: pixel.V(float64(320), float64(40))}
				foods[8] = food{vec: pixel.V(float64(360), float64(40))}
				foods[9] = food{vec: pixel.V(float64(400), float64(40))}
				foods[10] = food{vec: pixel.V(float64(440), float64(40))}
				foods[11] = food{vec: pixel.V(float64(580), float64(40))}
				foods[12] = food{vec: pixel.V(float64(620), float64(40))}
				foods[13] = food{vec: pixel.V(float64(660), float64(40))}
				foods[14] = food{vec: pixel.V(float64(700), float64(40))}
				foods[15] = food{vec: pixel.V(float64(740), float64(40))}
				foods[16] = food{vec: pixel.V(float64(780), float64(40))}
				foods[17] = food{vec: pixel.V(float64(820), float64(40))}
				foods[18] = food{vec: pixel.V(float64(860), float64(40))}
				foods[19] = food{vec: pixel.V(float64(900), float64(40))}
				foods[20] = food{vec: pixel.V(float64(940), float64(40))}
				foods[21] = food{vec: pixel.V(float64(980), float64(40))}
			}
			{ //row2
				foods[22] = food{vec: pixel.V(float64(40), float64(80))}
				foods[23] = food{vec: pixel.V(float64(320), float64(80))}
				foods[24] = food{vec: pixel.V(float64(360), float64(80))}
				foods[25] = food{vec: pixel.V(float64(400), float64(80))}
				foods[26] = food{vec: pixel.V(float64(440), float64(80))}
				foods[27] = food{vec: pixel.V(float64(580), float64(80))}
				foods[28] = food{vec: pixel.V(float64(620), float64(80))}
				foods[29] = food{vec: pixel.V(float64(660), float64(80))}
				foods[30] = food{vec: pixel.V(float64(700), float64(80))}
				foods[31] = food{vec: pixel.V(float64(980), float64(80))}
			}
			{ //row3
				foods[32] = food{vec: pixel.V(float64(40), float64(120))}
				foods[33] = food{vec: pixel.V(float64(320), float64(120))}
				foods[34] = food{vec: pixel.V(float64(360), float64(120))}
				foods[35] = food{vec: pixel.V(float64(400), float64(120))}
				foods[36] = food{vec: pixel.V(float64(440), float64(120))}
				foods[37] = food{vec: pixel.V(float64(580), float64(120))}
				foods[38] = food{vec: pixel.V(float64(620), float64(120))}
				foods[39] = food{vec: pixel.V(float64(660), float64(120))}
				foods[40] = food{vec: pixel.V(float64(700), float64(120))}
				foods[41] = food{vec: pixel.V(float64(980), float64(120))}
			}
			{ //row4
				foods[42] = food{vec: pixel.V(float64(40), float64(160))}
				foods[43] = food{vec: pixel.V(float64(320), float64(160))}
				foods[44] = food{vec: pixel.V(float64(360), float64(160))}
				foods[45] = food{vec: pixel.V(float64(400), float64(160))}
				foods[46] = food{vec: pixel.V(float64(440), float64(160))}
				foods[47] = food{vec: pixel.V(float64(580), float64(160))}
				foods[48] = food{vec: pixel.V(float64(620), float64(160))}
				foods[49] = food{vec: pixel.V(float64(660), float64(160))}
				foods[50] = food{vec: pixel.V(float64(700), float64(160))}
				foods[51] = food{vec: pixel.V(float64(980), float64(160))}
			}
			{ //row5
				foods[52] = food{vec: pixel.V(float64(40), float64(200))}
				foods[53] = food{vec: pixel.V(float64(320), float64(200))}
				foods[54] = food{vec: pixel.V(float64(360), float64(200))}
				foods[55] = food{vec: pixel.V(float64(400), float64(200))}
				foods[56] = food{vec: pixel.V(float64(440), float64(200))}
				foods[57] = food{vec: pixel.V(float64(580), float64(200))}
				foods[58] = food{vec: pixel.V(float64(620), float64(200))}
				foods[59] = food{vec: pixel.V(float64(660), float64(200))}
				foods[60] = food{vec: pixel.V(float64(700), float64(200))}
				foods[61] = food{vec: pixel.V(float64(980), float64(200))}
			}
			{ //row6
				foods[62] = food{vec: pixel.V(float64(40), float64(240))}
				foods[63] = food{vec: pixel.V(float64(320), float64(240))}
				foods[64] = food{vec: pixel.V(float64(360), float64(240))}
				foods[65] = food{vec: pixel.V(float64(400), float64(240))}
				foods[66] = food{vec: pixel.V(float64(440), float64(240))}
				foods[72] = food{vec: pixel.V(float64(485), float64(240))}
				foods[73] = food{vec: pixel.V(float64(535), float64(240))}
				foods[67] = food{vec: pixel.V(float64(580), float64(240))}
				foods[68] = food{vec: pixel.V(float64(620), float64(240))}
				foods[69] = food{vec: pixel.V(float64(660), float64(240))}
				foods[70] = food{vec: pixel.V(float64(700), float64(240))}
				foods[71] = food{vec: pixel.V(float64(980), float64(240))}
			}
			{ //row7
				foods[74] = food{vec: pixel.V(float64(40), float64(280))}
				foods[75] = food{vec: pixel.V(float64(320), float64(280))}
				foods[76] = food{vec: pixel.V(float64(700), float64(280))}
				foods[77] = food{vec: pixel.V(float64(980), float64(280))}
			}
			{ //row8
				foods[78] = food{vec: pixel.V(float64(40), float64(320))}
				foods[79] = food{vec: pixel.V(float64(80), float64(320))}
				foods[80] = food{vec: pixel.V(float64(120), float64(320))}
				foods[81] = food{vec: pixel.V(float64(160), float64(320))}
				foods[82] = food{vec: pixel.V(float64(200), float64(320))}
				foods[83] = food{vec: pixel.V(float64(240), float64(320))}
				foods[84] = food{vec: pixel.V(float64(280), float64(320))}
				foods[86] = food{vec: pixel.V(float64(320), float64(320))}
				foods[87] = food{vec: pixel.V(float64(700), float64(320))}
				foods[88] = food{vec: pixel.V(float64(740), float64(320))}
				foods[89] = food{vec: pixel.V(float64(780), float64(320))}
				foods[90] = food{vec: pixel.V(float64(820), float64(320))}
				foods[91] = food{vec: pixel.V(float64(860), float64(320))}
				foods[92] = food{vec: pixel.V(float64(900), float64(320))}
				foods[93] = food{vec: pixel.V(float64(940), float64(320))}
				foods[94] = food{vec: pixel.V(float64(980), float64(320))}
			}
			{ //row9
				foods[95] = food{vec: pixel.V(float64(40), float64(360))}
				foods[96] = food{vec: pixel.V(float64(80), float64(360))}
				foods[97] = food{vec: pixel.V(float64(120), float64(360))}
				foods[98] = food{vec: pixel.V(float64(240), float64(360))}
				foods[99] = food{vec: pixel.V(float64(280), float64(360))}
				foods[194] = food{vec: pixel.V(float64(320), float64(360))}
				foods[195] = food{vec: pixel.V(float64(700), float64(360))}
				foods[196] = food{vec: pixel.V(float64(740), float64(360))}
				foods[197] = food{vec: pixel.V(float64(780), float64(360))}
				foods[198] = food{vec: pixel.V(float64(900), float64(360))}
				foods[199] = food{vec: pixel.V(float64(940), float64(360))}
				foods[200] = food{vec: pixel.V(float64(980), float64(360))}
			}
			{ //row10
				foods[100] = food{vec: pixel.V(float64(40), float64(400))}
				foods[101] = food{vec: pixel.V(float64(80), float64(400))}
				foods[102] = food{vec: pixel.V(float64(120), float64(400))}
				foods[103] = food{vec: pixel.V(float64(160), float64(400))}
				foods[104] = food{vec: pixel.V(float64(200), float64(400))}
				foods[105] = food{vec: pixel.V(float64(240), float64(400))}
				foods[106] = food{vec: pixel.V(float64(280), float64(400))}
				foods[107] = food{vec: pixel.V(float64(320), float64(400))}
				foods[108] = food{vec: pixel.V(float64(700), float64(400))}
				foods[109] = food{vec: pixel.V(float64(740), float64(400))}
				foods[110] = food{vec: pixel.V(float64(780), float64(400))}
				foods[111] = food{vec: pixel.V(float64(820), float64(400))}
				foods[112] = food{vec: pixel.V(float64(860), float64(400))}
				foods[113] = food{vec: pixel.V(float64(900), float64(400))}
				foods[114] = food{vec: pixel.V(float64(940), float64(400))}
				foods[115] = food{vec: pixel.V(float64(980), float64(400))}
			}
			{ //row11
				foods[116] = food{vec: pixel.V(float64(40), float64(440))}
				foods[117] = food{vec: pixel.V(float64(320), float64(440))}
				foods[118] = food{vec: pixel.V(float64(700), float64(440))}
				foods[119] = food{vec: pixel.V(float64(980), float64(440))}
			}
			{ //row12
				foods[120] = food{vec: pixel.V(float64(40), float64(480))}
				foods[121] = food{vec: pixel.V(float64(320), float64(480))}
				foods[122] = food{vec: pixel.V(float64(360), float64(480))}
				foods[123] = food{vec: pixel.V(float64(400), float64(480))}
				foods[124] = food{vec: pixel.V(float64(440), float64(480))}
				foods[125] = food{vec: pixel.V(float64(485), float64(480))}
				foods[126] = food{vec: pixel.V(float64(535), float64(480))}
				foods[127] = food{vec: pixel.V(float64(580), float64(480))}
				foods[128] = food{vec: pixel.V(float64(620), float64(480))}
				foods[129] = food{vec: pixel.V(float64(660), float64(480))}
				foods[130] = food{vec: pixel.V(float64(700), float64(480))}
				foods[131] = food{vec: pixel.V(float64(980), float64(480))}
			}
			{ //row13
				foods[132] = food{vec: pixel.V(float64(40), float64(520))}
				foods[133] = food{vec: pixel.V(float64(320), float64(520))}
				foods[134] = food{vec: pixel.V(float64(360), float64(520))}
				foods[135] = food{vec: pixel.V(float64(400), float64(520))}
				foods[136] = food{vec: pixel.V(float64(440), float64(520))}
				foods[137] = food{vec: pixel.V(float64(580), float64(520))}
				foods[138] = food{vec: pixel.V(float64(620), float64(520))}
				foods[139] = food{vec: pixel.V(float64(660), float64(520))}
				foods[140] = food{vec: pixel.V(float64(700), float64(520))}
				foods[141] = food{vec: pixel.V(float64(980), float64(520))}
			}
			{ //row14
				foods[142] = food{vec: pixel.V(float64(40), float64(560))}
				foods[143] = food{vec: pixel.V(float64(320), float64(560))}
				foods[144] = food{vec: pixel.V(float64(360), float64(560))}
				foods[145] = food{vec: pixel.V(float64(400), float64(560))}
				foods[146] = food{vec: pixel.V(float64(440), float64(560))}
				foods[147] = food{vec: pixel.V(float64(580), float64(560))}
				foods[148] = food{vec: pixel.V(float64(620), float64(560))}
				foods[149] = food{vec: pixel.V(float64(660), float64(560))}
				foods[150] = food{vec: pixel.V(float64(700), float64(560))}
				foods[151] = food{vec: pixel.V(float64(980), float64(560))}
			}
			{ //row15
				foods[152] = food{vec: pixel.V(float64(40), float64(600))}
				foods[153] = food{vec: pixel.V(float64(320), float64(600))}
				foods[154] = food{vec: pixel.V(float64(360), float64(600))}
				foods[155] = food{vec: pixel.V(float64(400), float64(600))}
				foods[156] = food{vec: pixel.V(float64(440), float64(600))}
				foods[157] = food{vec: pixel.V(float64(580), float64(600))}
				foods[158] = food{vec: pixel.V(float64(620), float64(600))}
				foods[159] = food{vec: pixel.V(float64(660), float64(600))}
				foods[160] = food{vec: pixel.V(float64(700), float64(600))}
				foods[161] = food{vec: pixel.V(float64(980), float64(600))}
			}
			{ //row16
				foods[162] = food{vec: pixel.V(float64(40), float64(640))}
				foods[163] = food{vec: pixel.V(float64(320), float64(640))}
				foods[164] = food{vec: pixel.V(float64(360), float64(640))}
				foods[165] = food{vec: pixel.V(float64(400), float64(640))}
				foods[166] = food{vec: pixel.V(float64(440), float64(640))}
				foods[167] = food{vec: pixel.V(float64(580), float64(640))}
				foods[168] = food{vec: pixel.V(float64(620), float64(640))}
				foods[169] = food{vec: pixel.V(float64(660), float64(640))}
				foods[170] = food{vec: pixel.V(float64(700), float64(640))}
				foods[171] = food{vec: pixel.V(float64(980), float64(640))}
			}
			{ //row 17
				foods[172] = food{vec: pixel.V(float64(40), float64(680))}
				foods[173] = food{vec: pixel.V(float64(80), float64(680))}
				foods[174] = food{vec: pixel.V(float64(120), float64(680))}
				foods[175] = food{vec: pixel.V(float64(160), float64(680))}
				foods[176] = food{vec: pixel.V(float64(200), float64(680))}
				foods[177] = food{vec: pixel.V(float64(240), float64(680))}
				foods[178] = food{vec: pixel.V(float64(280), float64(680))}
				foods[179] = food{vec: pixel.V(float64(320), float64(680))}
				foods[180] = food{vec: pixel.V(float64(360), float64(680))}
				foods[181] = food{vec: pixel.V(float64(400), float64(680))}
				foods[182] = food{vec: pixel.V(float64(440), float64(680))}
				foods[183] = food{vec: pixel.V(float64(580), float64(680))}
				foods[184] = food{vec: pixel.V(float64(620), float64(680))}
				foods[185] = food{vec: pixel.V(float64(660), float64(680))}
				foods[186] = food{vec: pixel.V(float64(700), float64(680))}
				foods[187] = food{vec: pixel.V(float64(740), float64(680))}
				foods[188] = food{vec: pixel.V(float64(780), float64(680))}
				foods[189] = food{vec: pixel.V(float64(820), float64(680))}
				foods[190] = food{vec: pixel.V(float64(860), float64(680))}
				foods[191] = food{vec: pixel.V(float64(900), float64(680))}
				foods[192] = food{vec: pixel.V(float64(940), float64(680))}
				foods[193] = food{vec: pixel.V(float64(980), float64(680))}
			}
		}
		imd.Color = colornames.Yellow
		for _, singleFood := range foods {
			imd.Push(singleFood.vec)
		}
		imd.Circle(8, 0)

	}
}

func ghostRun(ghost string, sprite *pixel.Sprite, pos pixel.Vec, win *pixelgl.Window) {
	timeToInstruction := 4.0
	ghostSpeed := 1.5
	last := time.Now()
	sumTimePassed := 0.0
	var seed int64
	var current pixel.Vec
	for _, c := range ghost {
		seed += int64(c)
	}

	left := pixel.V(-.1*ghostSpeed, 0)
	right := pixel.V(.1*ghostSpeed, 0)
	up := pixel.V(0, .1*ghostSpeed)
	down := pixel.V(0, -.1*ghostSpeed)

	for !win.Closed() {
		//fmt.Println(r.Intn(4))
		dt := time.Since(last).Seconds()
		last = time.Now()
		sumTimePassed += dt
		mat := pixel.IM
		mat = mat.Moved(pos)
		if ghost == "inky" {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(0.25, 0.25))
		} else {
			mat = mat.ScaledXY(pixel.ZV, pixel.V(0.13, 0.13))
		}
		if sumTimePassed >= timeToInstruction {
			random := rand.Int31n(3)
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
		if checkCollisionGhost(ghost, pos) || pos.X < 10 || pos.X > 1020 || pos.Y < 10 || pos.Y > 710 { //Check collision with all walls and the borders of the map
			fmt.Printf("%s is collisioning\n", ghost)
			current.X = -current.X
			current.Y = -current.Y
			pos.X += current.X
			pos.Y += current.Y
		}
		pos.X += current.X
		pos.Y += current.Y
		mat = mat.Moved(pos)
		sprite.Draw(win, mat)
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

	blinky, err := loadPicture("blinky.png")
	if err != nil {
		panic(err)
	}
	pacmanSpriteUp := pixel.NewSprite(pacmanUp, pacmanUp.Bounds())
	pacmanSpriteDown := pixel.NewSprite(pacmanDown, pacmanDown.Bounds())
	pacmanSpriteLeft := pixel.NewSprite(pacmanLeft, pacmanLeft.Bounds())
	pacmanSpriteRight := pixel.NewSprite(pacmanRight, pacmanRight.Bounds())

	blinkySprite := pixel.NewSprite(blinky, blinky.Bounds())

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
		go ghostRun("blinky", blinkySprite, pixel.V(float64(35*i)+25.0, 35), win)
	}

	for !win.Closed() {
		win.Clear(colornames.Black)
		mat := pixel.IM

		for i, singleFood := range foods { //Draw and check collision with food
			if checkCollision(currentPos, singleFood.vec) {
				imd.Color = colornames.Black
				imd.Push(singleFood.vec)
				imd.Circle(8, 0)
				delete(foods, i)
				score += 100
				fmt.Printf("Score: %d \n", score)
				if len(foods) == 0 {
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
	}
}

func main() {
	//usage: go run pacman.go -g <numGhosts>
	numGhosts, _ = strconv.Atoi(os.Args[2])
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}

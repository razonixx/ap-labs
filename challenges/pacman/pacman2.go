package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Node struct {
	id        int
	neighbors []*Node
	history   []*Node
	vec       pixel.Vec
	isEaten   bool
}

type wall struct {
	rect pixel.Rect
}

const constSpeed = 3.5

var imd = imdraw.New(nil)
var speed = constSpeed
var pause = true
var score = 0
var numGhosts int
var currentPos pixel.Vec
var nodes = make(map[int]*Node, 250)
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

func printNode(n Node) {
	fmt.Printf("Current ID: %d\n", n.id)
	fmt.Print("Neighbors: \n")
	for i := range n.neighbors {
		fmt.Print("ID: ")
		fmt.Println(n.neighbors[i].id)
	}

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

		//Declare nodes
		{
			{ //row 1
				nodes[0] = &Node{0, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(40)), false}
				nodes[1] = &Node{1, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(40)), false}
				nodes[2] = &Node{2, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(40)), false}
				nodes[3] = &Node{3, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(40)), false}
				nodes[4] = &Node{4, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(40)), false}
				nodes[5] = &Node{5, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(40)), false}
				nodes[6] = &Node{6, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(40)), false}
				nodes[7] = &Node{7, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(40)), false}
				nodes[8] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(40)), false}
				nodes[9] = &Node{9, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(40)), false}
				nodes[10] = &Node{10, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(40)), false}
				nodes[11] = &Node{11, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(40)), false}
				nodes[12] = &Node{12, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(40)), false}
				nodes[13] = &Node{13, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(40)), false}
				nodes[14] = &Node{14, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(40)), false}
				nodes[15] = &Node{15, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(40)), false}
				nodes[16] = &Node{16, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(40)), false}
				nodes[17] = &Node{17, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(40)), false}
				nodes[18] = &Node{18, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(40)), false}
				nodes[19] = &Node{19, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(40)), false}
				nodes[20] = &Node{20, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(40)), false}
				nodes[21] = &Node{21, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(40)), false}

				/*for i := 0; i < 11; i++ {
					nodes[i] = &Node{0, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(0+40*(i+1)), float64(40)), false}
				, false}
				for i := 12; i < 22; i++ {
					nodes[i] = &Node{0, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580+(40*(i-12))), float64(40)), false}
				, false}*/
			}
			{ //row2
				nodes[22] = &Node{22, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(80)), false}
				nodes[23] = &Node{23, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(80)), false}
				nodes[24] = &Node{24, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(80)), false}
				nodes[25] = &Node{25, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(80)), false}
				nodes[26] = &Node{26, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(80)), false}
				nodes[27] = &Node{27, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(80)), false}
				nodes[28] = &Node{28, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(80)), false}
				nodes[29] = &Node{29, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(80)), false}
				nodes[30] = &Node{30, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(80)), false}
				nodes[31] = &Node{31, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(80)), false}
			}
			{ //row3
				nodes[32] = &Node{32, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(120)), false}
				nodes[33] = &Node{33, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(120)), false}
				nodes[34] = &Node{34, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(120)), false}
				nodes[35] = &Node{35, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(120)), false}
				nodes[36] = &Node{36, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(120)), false}
				nodes[37] = &Node{37, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(120)), false}
				nodes[38] = &Node{38, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(120)), false}
				nodes[39] = &Node{39, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(120)), false}
				nodes[40] = &Node{40, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(120)), false}
				nodes[41] = &Node{41, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(120)), false}
			}
			{ //row4
				nodes[42] = &Node{42, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(160)), false}
				nodes[43] = &Node{43, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(160)), false}
				nodes[44] = &Node{44, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(160)), false}
				nodes[45] = &Node{45, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(160)), false}
				nodes[46] = &Node{46, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(160)), false}
				nodes[47] = &Node{47, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(160)), false}
				nodes[48] = &Node{48, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(160)), false}
				nodes[49] = &Node{49, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(160)), false}
				nodes[50] = &Node{50, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(160)), false}
				nodes[51] = &Node{51, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(160)), false}
			}
			{ //row5
				nodes[52] = &Node{52, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(200)), false}
				nodes[53] = &Node{53, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(200)), false}
				nodes[54] = &Node{54, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(200)), false}
				nodes[55] = &Node{55, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(200)), false}
				nodes[56] = &Node{56, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(200)), false}
				nodes[57] = &Node{57, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(200)), false}
				nodes[58] = &Node{58, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(200)), false}
				nodes[59] = &Node{59, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(200)), false}
				nodes[60] = &Node{60, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(200)), false}
				nodes[61] = &Node{61, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(200)), false}
			}
			{ //row6
				nodes[62] = &Node{62, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(240)), false}
				nodes[63] = &Node{63, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(240)), false}
				nodes[64] = &Node{64, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(240)), false}
				nodes[65] = &Node{65, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(240)), false}
				nodes[66] = &Node{66, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(240)), false}
				nodes[72] = &Node{72, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(240)), false}
				nodes[73] = &Node{73, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(240)), false}
				nodes[67] = &Node{67, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(240)), false}
				nodes[68] = &Node{68, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(240)), false}
				nodes[69] = &Node{69, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(240)), false}
				nodes[70] = &Node{70, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(240)), false}
				nodes[71] = &Node{71, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(240)), false}
			}
			{ //row7
				nodes[74] = &Node{1, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(280)), false}
				nodes[75] = &Node{1, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(280)), false}
				nodes[76] = &Node{1, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(280)), false}
				nodes[77] = &Node{1, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(280)), false}
			}
			{ //row8
				nodes[78] = &Node{7, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(320)), false}
				nodes[79] = &Node{7, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(320)), false}
				nodes[80] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(320)), false}
				nodes[81] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(320)), false}
				nodes[82] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(320)), false}
				nodes[83] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(320)), false}
				nodes[84] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(320)), false}
				nodes[86] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(320)), false}
				nodes[87] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(320)), false}
				nodes[88] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(320)), false}
				nodes[89] = &Node{8, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(320)), false}
				nodes[90] = &Node{9, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(320)), false}
				nodes[91] = &Node{9, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(320)), false}
				nodes[92] = &Node{9, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(320)), false}
				nodes[93] = &Node{9, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(320)), false}
				nodes[94] = &Node{9, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(320)), false}
			}
			{ //row9
				nodes[95] = &Node{95, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(360)), false}
				nodes[96] = &Node{96, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(360)), false}
				nodes[97] = &Node{97, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(360)), false}
				nodes[98] = &Node{98, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(360)), false}
				nodes[99] = &Node{99, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(360)), false}
				nodes[194] = &Node{194, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(360)), false}
				nodes[195] = &Node{195, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(360)), false}
				nodes[196] = &Node{196, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(360)), false}
				nodes[197] = &Node{197, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(360)), false}
				nodes[198] = &Node{198, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(360)), false}
				nodes[199] = &Node{199, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(360)), false}
				nodes[200] = &Node{200, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(360)), false}
			}
			{ //row10
				nodes[100] = &Node{100, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(400)), false}
				nodes[101] = &Node{101, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(400)), false}
				nodes[102] = &Node{102, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(400)), false}
				nodes[103] = &Node{103, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(400)), false}
				nodes[104] = &Node{104, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(400)), false}
				nodes[105] = &Node{105, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(400)), false}
				nodes[106] = &Node{106, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(400)), false}
				nodes[107] = &Node{107, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(400)), false}
				nodes[108] = &Node{108, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(400)), false}
				nodes[109] = &Node{109, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(400)), false}
				nodes[110] = &Node{110, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(400)), false}
				nodes[111] = &Node{111, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(400)), false}
				nodes[112] = &Node{112, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(400)), false}
				nodes[113] = &Node{113, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(400)), false}
				nodes[114] = &Node{114, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(400)), false}
				nodes[115] = &Node{115, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(400)), false}
			}
			{ //row11
				nodes[116] = &Node{116, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(440)), false}
				nodes[117] = &Node{117, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(440)), false}
				nodes[118] = &Node{118, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(440)), false}
				nodes[119] = &Node{119, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(440)), false}
			}
			{ //row12
				nodes[120] = &Node{120, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(480)), false}
				nodes[121] = &Node{121, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(480)), false}
				nodes[122] = &Node{122, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(480)), false}
				nodes[123] = &Node{123, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(480)), false}
				nodes[124] = &Node{124, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(480)), false}
				nodes[125] = &Node{125, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(485), float64(480)), false}
				nodes[126] = &Node{126, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(535), float64(480)), false}
				nodes[127] = &Node{127, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(480)), false}
				nodes[128] = &Node{128, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(480)), false}
				nodes[129] = &Node{129, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(480)), false}
				nodes[130] = &Node{130, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(480)), false}
				nodes[131] = &Node{131, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(480)), false}
			}
			{ //row13
				nodes[132] = &Node{132, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(520)), false}
				nodes[133] = &Node{133, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(520)), false}
				nodes[134] = &Node{134, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(520)), false}
				nodes[135] = &Node{135, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(520)), false}
				nodes[136] = &Node{136, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(520)), false}
				nodes[137] = &Node{137, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(520)), false}
				nodes[138] = &Node{138, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(520)), false}
				nodes[139] = &Node{139, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(520)), false}
				nodes[140] = &Node{140, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(520)), false}
				nodes[141] = &Node{141, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(520)), false}
			}
			{ //row14
				nodes[142] = &Node{142, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(560)), false}
				nodes[143] = &Node{143, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(560)), false}
				nodes[144] = &Node{144, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(560)), false}
				nodes[145] = &Node{145, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(560)), false}
				nodes[146] = &Node{146, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(560)), false}
				nodes[147] = &Node{147, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(560)), false}
				nodes[148] = &Node{148, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(560)), false}
				nodes[149] = &Node{149, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(560)), false}
				nodes[150] = &Node{150, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(560)), false}
				nodes[151] = &Node{151, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(560)), false}
			}
			{ //row15
				nodes[152] = &Node{152, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(600)), false}
				nodes[153] = &Node{153, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(600)), false}
				nodes[154] = &Node{154, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(600)), false}
				nodes[155] = &Node{155, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(600)), false}
				nodes[156] = &Node{156, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(600)), false}
				nodes[157] = &Node{157, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(600)), false}
				nodes[158] = &Node{158, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(600)), false}
				nodes[159] = &Node{159, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(600)), false}
				nodes[160] = &Node{160, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(600)), false}
				nodes[161] = &Node{161, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(600)), false}
			}
			{ //row16
				nodes[162] = &Node{162, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(640)), false}
				nodes[163] = &Node{163, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(640)), false}
				nodes[164] = &Node{164, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(640)), false}
				nodes[165] = &Node{165, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(640)), false}
				nodes[166] = &Node{166, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(640)), false}
				nodes[167] = &Node{167, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(640)), false}
				nodes[168] = &Node{168, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(640)), false}
				nodes[169] = &Node{169, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(640)), false}
				nodes[170] = &Node{170, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(640)), false}
				nodes[171] = &Node{171, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(640)), false}
			}
			{ //row 17
				nodes[172] = &Node{172, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(40), float64(680)), false}
				nodes[173] = &Node{173, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(80), float64(680)), false}
				nodes[174] = &Node{174, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(120), float64(680)), false}
				nodes[175] = &Node{175, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(160), float64(680)), false}
				nodes[176] = &Node{176, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(200), float64(680)), false}
				nodes[177] = &Node{177, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(240), float64(680)), false}
				nodes[178] = &Node{178, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(280), float64(680)), false}
				nodes[179] = &Node{179, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(320), float64(680)), false}
				nodes[180] = &Node{180, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(360), float64(680)), false}
				nodes[181] = &Node{181, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(400), float64(680)), false}
				nodes[182] = &Node{182, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(440), float64(680)), false}
				nodes[183] = &Node{183, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(580), float64(680)), false}
				nodes[184] = &Node{184, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(620), float64(680)), false}
				nodes[185] = &Node{185, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(660), float64(680)), false}
				nodes[186] = &Node{186, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(700), float64(680)), false}
				nodes[187] = &Node{187, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(740), float64(680)), false}
				nodes[188] = &Node{188, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(780), float64(680)), false}
				nodes[189] = &Node{189, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(820), float64(680)), false}
				nodes[190] = &Node{190, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(860), float64(680)), false}
				nodes[191] = &Node{191, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(900), float64(680)), false}
				nodes[192] = &Node{192, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(940), float64(680)), false}
				nodes[193] = &Node{193, make([]*Node, 0, 10), make([]*Node, 0, 300), pixel.V(float64(980), float64(680)), false}
			}
		}

		//Declare neighbors
		{
			nodes[0].neighbors = append(nodes[0].neighbors, nodes[1], nodes[22])
			nodes[1].neighbors = append(nodes[1].neighbors, nodes[0], nodes[2])
			nodes[2].neighbors = append(nodes[2].neighbors, nodes[1], nodes[3])
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

	nodesEaten := 0
	last := time.Now()
	dtCheck := 0.0
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		dtCheck += dt
		win.Clear(colornames.Black)
		if dtCheck >= .0 {

			mat := pixel.IM

			for _, singlenode := range nodes { //Draw and check collision with node
				if checkCollision(currentPos, singlenode.vec) {
					//lastnodeEaten = i
					if !singlenode.isEaten {
						imd.Color = colornames.Black
						imd.Push(singlenode.vec)
						imd.Circle(8, 0)
						//delete(nodes, i)
						score += 100
						nodesEaten++
						fmt.Printf("Score: %d \n", score)
						//fmt.Printf("Last node eaten: %d\n", lastnodeEaten)
						if nodesEaten == len(nodes) {
							log.Println("YOU WIN!!")
							os.Exit(0)
						}
					}
					singlenode.isEaten = true
					printNode(*singlenode)
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

func calculateH(current, destination pixel.Vec) int32 {
	return int32(math.Abs(current.X-destination.X) + math.Abs(current.Y-destination.Y))
}

func main() {
	//usage: go run pacman.go -g <numGhosts>
	numGhosts, _ = strconv.Atoi(os.Args[2])
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}

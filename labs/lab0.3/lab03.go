package main

import (
	"fmt"
	"math"
	"strings"
)

func Pic(dx, dy int) [][]uint8 {
	mat := make([][]uint8, dy)
	for y := range mat {
		mat[y] = make([]uint8, dx)
		for x := range mat[y] {
			mat[y][x] = uint8((x * y))
		}
	}
	return mat
}

func wordCount(s string) map[string]int {
	var arr = strings.Fields(s)
	var m = make(map[string]int)
	for index := 0; index < len(arr); index++ {
		x := m[arr[index]]
		x++
		m[arr[index]] = x
	}
	return m
}

type Point struct{ x, y float64 }

func (p Point) X() float64 {
	return p.x
}

func (p Point) Y() float64 {
	return p.y
}

// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}

type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

func printPoint(p Point) {
	fmt.Printf("%f, %f\n", p.X(), p.Y())
}

func (p Point) printPoint() {
	fmt.Printf("%f, %f\n", p.X(), p.Y())
}

func main() {
	var m = wordCount("I would very very much like like to work at intel .")
	fmt.Println(m)

	var p = Point{4, 12}

	printPoint(p)
	p.printPoint()
}

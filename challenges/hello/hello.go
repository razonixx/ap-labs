package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const pi = math.Pi
const byt uint8 = 0

//Declaration with no capital letter is private
type vector struct {
	magnitude int
	direction int
}

//Vertex Declaration with capital letter is public
type Vertex struct {
	X float64
	Y float64
}

func cosa() {
	//defer evaluates when called, but executes until its parent function returns
	defer fmt.Println("Goodbye world")
	fmt.Println("Hello world")
	fmt.Printf("%d\n", 10)
	fmt.Println(time.Now())
	rand.Seed(10)
	fmt.Printf("%d\n", rand.Intn(10))
	fmt.Printf("%f\n", math.Pi)
	fmt.Printf("%d\n", add(2, 2))

	a, b, c := swap("Hello", "world", 5)

	fmt.Printf("%d %s %s\n", a, b, c)

	fmt.Printf("%d\n", square(4))

	var sq float64 = 8596
	fmt.Printf("My result: %f\n", sqrt(sq))
	fmt.Printf("math.Sqrt result: %f\n", math.Sqrt(sq))
}

func add(x int, y int) int {
	return x + y
}

func swap(x string, y string, z int) (int, string, string) {
	return z, y, x
}

func square(num int) (res int) {
	res = num * num
	return
}

func arr() {
	var x, y, z int = 1, 2, 3
	if y == 1 {
		//var a, b, c = "js", "like", true
		//d := "Implicit type like python, but only inside functinos"
	}
	w := 0
	for i := 0; i < 10; i++ {
		w += (x + y + z) * i
	}

	for w > 0 {
		x++
	}
}

func pow(x, n, lim float64) float64 {
	// This if statement works like a for, it declares a variable (v) that is only visible inside the if itself and an else, if it had one
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

func sqrt(x float64) float64 {
	var z float64 = 1
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func pointers() {
	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j
}

func structs() {
	var v1 = vector{2, 4}

	fmt.Printf("%d\n", v1.magnitude)
	fmt.Printf("%f\n", Vertex{2.1, 4.3})
}

func arrs() {
	//var list [10]int
	var arr = [10]int{1, 2, 3, 4, 5}
	hello := "hello"
	fmt.Printf("%s\n", hello)
	//Slice: arr[0:6]
	//Includes lower bound but excludes upper bound
	fmt.Printf("%d\n", arr[0:1])
}

func fMake() {
	a := make([]int, 4)
	a[3] = 1
	printSlice("a", a)
}

func printSlice(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}

func mat() {
	// Create a tic-tac-toe board.
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	// The players take turns.
	//board[0][0] = "X"
	//board[2][2] = "O"
	board[1][2] = "X"
	//board[1][0] = "O"
	//board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
}

func arrAppend() {
	var arr []int
	arr = append(arr, 1, 2, 3)
	fmt.Println(arr)
}

func forEach() {
	num := []int{1, 2, 4, 8, 16, 32, 64, 128}

	//foreach with variable of the value of the array at index i
	for i, val := range num {
		fmt.Printf("2^%d = %d\n", i, val)
	}

	//foreach with no additional varaible
	for i := range num {
		fmt.Printf("2^%d = 0\n", i)
	}

	//foreach without using index, only values in collection
	for _, val := range num {
		fmt.Printf("%d\n", val)
	}
}

func pic(dx, dy int) [][]uint8 {
	pic := make([][]uint8, dy)

	for y := range pic {

		pic[y] = make([]uint8, dx)

		for x := range pic[y] {
			pic[y][x] = uint8((x + y) / 2)
		}
	}
	return pic
}

func main() {
	//cosa()
	//pointers()
	//structs()
	//arrs()
	//fMake()
	//mat()
	//arrAppend()
	//forEach()

	fmt.Println(pic(5, 5))
}

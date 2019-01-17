package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const pi = math.Pi
const byt uint8 = 0

//Type with no capital letter is private
type vector struct {
	magnitude int
	direction int
}

//Vertex Type with capital letter is public
type Vertex struct {
	X float64
	Y float64
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
	hello := "hello"
	fmt.Printf("%s\n", hello)
}

func main() {
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

	pointers()
	structs()
	arrs()
}

package main

import (
	"fmt"
	"math"
	"math/rand"
)

func add(x int, y int) int{
	return x + y
}

func swap(x string, y string, z int) (int, string, string){
	return z, y, x
}

func square(num int) (res int){
	res = num * num
	return
}

func main(){
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
}
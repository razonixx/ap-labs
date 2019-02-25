package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	//var port string
	/*conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Fatal(err)
	}*/
	done := make(chan int)

	/*var string1 string
	string1 = os.Args[1]
	fmt.Println(strings.Split(string1, "=")[0])*/

	y := len(os.Args)
	//fmt.Printf("Recieved %d parameters from command line.\n", y)

	for i := 1; i < y; i++ {
		city := strings.Split(os.Args[i], "=")[0]
		port := strings.Split(os.Args[i], "=")[1]
		//fmt.Println("City: " + city)
		//fmt.Println("Poft: " + port)
		conn, err := net.Dial("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for true {
				fmt.Printf(city + ": ")
				time.Sleep((1 * time.Second) + (7 * time.Nanosecond))
				//done <- 0
			}
		}()
		go func() {
			_, err := io.Copy(os.Stdout, conn)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("done")
			done <- 2
		}()
	}

	x := 1
	x = <-done // wait for background goroutine to finish
	log.Println("Channel Closed with value: ", x)
	close(done)
}

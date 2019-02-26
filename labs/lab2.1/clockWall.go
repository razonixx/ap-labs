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

func clockWall(conn net.Conn, channel chan int, city string) {
	var err error
	for true {
		time.Sleep((1 * time.Second))
		fmt.Printf(city + ": ")
		_, err = io.CopyN(os.Stdout, conn, 9)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("\nConnection with " + city + " clock closed.")
	channel <- 2
}

func main() {
	done := make(chan int)
	y := len(os.Args)

	for i := 1; i < y; i++ {
		city := strings.Split(os.Args[i], "=")[0]
		port := strings.Split(os.Args[i], "=")[1]
		conn, err := net.Dial("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		go clockWall(conn, done, city)
	}

	_ = <-done // wait for background goroutine to finish
	close(done)
}

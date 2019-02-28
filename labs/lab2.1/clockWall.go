package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func clockWall(conn net.Conn, city string, wg *sync.WaitGroup) {
	for true {
		_, err := io.Copy(os.Stdout, conn)
		if err == nil {
			break
		}
	}
	log.Println("Connection with " + city + " clock closed.")
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	y := len(os.Args)

	for i := 1; i < y; i++ {
		city := strings.Split(os.Args[i], "=")[0]
		port := strings.Split(os.Args[i], "=")[1]
		conn, err := net.Dial("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go clockWall(conn, city, &wg)
	}
	wg.Wait()
	log.Println("All clocks are closed. Exiting...")
}

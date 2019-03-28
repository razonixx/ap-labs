// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	var user string
	var server string

	if os.Args[1] == "-user" {
		user = os.Args[2]
	}
	if os.Args[3] == "-server" {
		server = os.Args[4]
	}
	fmt.Printf("User: %s Server: %s\n", user, server)

	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(conn, user)

	done := make(chan struct{})
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err == nil {
			log.Fatal(err)
		}
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

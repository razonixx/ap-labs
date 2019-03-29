// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	var user string
	var server string

	if len(os.Args) < 5 {
		log.Fatalf("Usage: ./client -user <user> -server <server>\n")
	}

	if os.Args[1] == "-user" {
		user = os.Args[2]
	}
	if os.Args[3] == "-server" {
		server = os.Args[4]
	}

	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.WriteString(conn, user+"\n")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		mustCopy(os.Stdout, conn)
		log.Println("Server closed")
		done <- struct{}{} // signal the main goroutine
	}()

	go func() {
		mustCopy(conn, os.Stdin)
	}()
	<-done // wait for background goroutine to finish
	conn.Close()
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}

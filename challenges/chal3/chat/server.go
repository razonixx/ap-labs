// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type user struct {
	userName string
	ip       string
	channel  chan string
}

var users = make([]user, 1000)
var userCount = 0

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all except sender
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	defer conn.Close()
	ch := make(chan string)       // outgoing client messages
	serverCh := make(chan string) // outgoing client messages
	go clientWriter(conn, serverCh)
	input := bufio.NewScanner(conn)
	input.Scan() //Go to username
	who := input.Text()
	users[userCount].userName = who
	users[userCount].ip = conn.RemoteAddr().String()
	users[userCount].channel = ch
	userCount++ //Increase number of users
	go clientWriterNoNewLine(conn, ch, who)
	serverCh <- "irc-server> Welcome to the IRC server!"
	serverCh <- "irc-server> Your user, " + who + ", has succesfully logged in!"
	messages <- who + " has arrived"
	entering <- ch

	log.Printf("New connected user: %s\n", who)
	for input.Scan() {
		splitted := strings.Split(input.Text(), " ")

		if splitted[0] == "/users" {
			ch <- "irc-server> Users: "
			for _, user := range users {
				if len(user.userName) > 0 {
					ch <- "irc-server> " + user.userName
				}
			}

		} else if splitted[0] == "/msg" {
			for i, user := range users {
				if user.userName == splitted[1] {
					tempString := "Direct message from " + who + ": "
					for _, token := range splitted[2:] {
						tempString += (token + " ")
					}
					user.channel <- tempString
					break
				}
				if i == userCount {
					ch <- "irc-server> Username not found"
				}
			}

		} else if splitted[0] == "/time" {
			ch <- "irc-server> Time: " + time.Now().Format("15:04:05")

		} else if splitted[0] == "/user" {
			for i, user := range users {
				if user.userName == splitted[1] {
					ch <- "irc-server> Username: " + user.userName + " IP Address: " + user.ip
					break
				}
				if i == userCount {
					ch <- "irc-server> Username not found"
				}
			}

		} else {
			messages <- who + "> " + input.Text()
		}
	}

	leaving <- ch
	messages <- who + " has left"
	log.Printf("%s disconnected\n", who)
	userCount-- //Decrease number of users
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func clientWriterNoNewLine(conn net.Conn, ch <-chan string, user string) {
	for msg := range ch {
		_, err := fmt.Fprintf(conn, "%s\n", msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//!-handleConn

//!+main
func main() {
	var host string
	var port string

	if os.Args[1] == "-host" {
		host = os.Args[2]
	}
	if os.Args[3] == "-port" {
		port = os.Args[4]
	}
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting server at %s:%s\n", host, port)

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main

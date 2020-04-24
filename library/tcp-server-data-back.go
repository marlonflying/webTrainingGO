package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func handleConn(conn net.Conn) {
	message, err := bufio.NewReader(conn).ReadString('?') //? EOF
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Print("Message Recieved from the client: ", string(message))
	conn.Write([]byte(message + "\n"))
	conn.Close()
}

func main() {
	listener, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		log.Fatal("Error starting tcp server : ", err)
	}
	defer listener.Close()
	log.Println("Listening on " + connHost + ":" + connPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err.Error())
		}
		go handleConn(conn)
	}
}

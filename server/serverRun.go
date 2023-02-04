package server

import (
	"fmt"
	"log"
	"net"
)

type Chat struct {
	clients []Client
	channel chan Message
}

type Client struct {
	name string
	conn net.Conn
}

type Message struct {
	user Client
	msg  string
}

func StartServer(serverPort string) {
	addr := fmt.Sprintf("localhost:%s", serverPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	var ServerChat Chat
	ServerChat.channel = make(chan Message)
	ServerChat.clients = []Client{}

	defer listener.Close()
	go ServerChat.ProcessMessages()
	log.Printf("Listening on the port :%s", serverPort)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go ServerChat.ProcessClient(conn)
		}
	}
}

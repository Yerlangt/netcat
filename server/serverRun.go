package server

import (
	"fmt"
	"log"
	"net"
)

func StartServer(serverPort string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", serverPort))
	if err != nil {
		log.Fatal("ERROR on starting server")
	}
	log.Printf("Listening on the port :%s", serverPort)

	ServerChat := Chat{[]Client{}, make(chan Message), []Message{}}

	defer listener.Close()
	go ServerChat.ProcessMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else if len(ServerChat.clients) == 10 {
			fmt.Fprintln(conn, "chat is full, try later...")
		} else {
			go ServerChat.ProcessClient(conn)
		}
	}
}

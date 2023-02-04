package server

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/helpers"
	"time"
)

// TODO validations of name and messages
// History
func (ServerChat *Chat) ProcessMessages() {
	var newMessage Message

	for {
		newMessage = <-ServerChat.channel
		for _, client := range ServerChat.clients {
			if newMessage.user != client {
				fmt.Fprintf(client.conn, "\n")
				fmt.Fprint(client.conn, newMessage.msg)
				fmt.Fprint(client.conn, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]:", client.name))
			}
		}
	}
}

func (ServerChat *Chat) ProcessClient(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	fmt.Fprint(conn, helpers.StarterMessage())

	var name string
	for scanner.Scan() {
		name = scanner.Text()
		break
	}
	newClient := ServerChat.newClientAdd(conn, name)

	var newMessage Message

	for scanner.Scan() {
		newMessage.user = newClient
		newMessage.msg = scanner.Text()
		fmt.Fprint(conn, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]:", name))
		ServerChat.channel <- FormatTextMessage(newMessage)
	}
	ServerChat.channel <- Message{newClient, fmt.Sprintf("%s has left our chat...\n", newClient.name)}
}

func (ServerChat *Chat) newClientAdd(conn net.Conn, name string) Client {
	var newClient Client
	newClient.conn = conn
	newClient.name = name
	ServerChat.clients = append(ServerChat.clients, newClient)

	var newMessage Message
	newMessage.user = newClient
	newMessage.msg = fmt.Sprintf("%s has joined our chat...\n", newClient.name)
	ServerChat.channel <- newMessage
	fmt.Fprint(conn, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]:", name))
	return newClient
}

func FormatTextMessage(msg Message) Message {
	return Message{msg.user, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]: %s\n", msg.user.name, msg.msg)}
}

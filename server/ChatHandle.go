package server

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/helpers"
	"time"
)

type Chat struct {
	clients []Client
	channel chan Message
	history []Message
}

type Client struct {
	name string
	conn net.Conn
}

type Message struct {
	user Client
	msg  string
}

// TODO validations of name
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

	helpers.PrintStarterMessage(conn)

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
		if len(newMessage.msg) == 0 {
			PrintBaseMessage(newClient)
			continue
		}
		PrintBaseMessage(newClient)
		ServerChat.channel <- Message{newClient, GetRegularTextMessage(newMessage)}
	}
	ServerChat.channel <- Message{newClient, GetClientDeleteMessage(newClient)}
	ServerChat.DeleteClient(newClient)
}

func (ServerChat *Chat) newClientAdd(conn net.Conn, name string) Client {
	var newClient Client
	newClient.conn = conn
	newClient.name = name

	ServerChat.clients = append(ServerChat.clients, newClient)

	ServerChat.channel <- Message{newClient, GetClientAddMessage(newClient)}
	PrintBaseMessage(newClient)

	return newClient
}

func (ServerChat *Chat) DeleteClient(client Client) {
	for i, v := range ServerChat.clients {
		if v == client {
			ServerChat.clients = append(ServerChat.clients[:i], ServerChat.clients[i+1:]...)
		}
	}
}

func GetRegularTextMessage(msg Message) string {
	return fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]: %s\n", msg.user.name, msg.msg)
}

func PrintBaseMessage(client Client) {
	fmt.Fprint(client.conn, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]:", client.name))
}

func GetClientAddMessage(client Client) string {
	return fmt.Sprintf("%s has joined our chat...\n", client.name)
}

func GetClientDeleteMessage(client Client) string {
	return fmt.Sprintf("%s has left our chat...\n", client.name)
}

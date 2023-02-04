package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net-cat/helpers"
	"time"
)

func (ServerChat *Chat) ProcessMessages() {
	var newMessage Message

	for {
		newMessage = <-ServerChat.channel
		for _, client := range ServerChat.clients {
			if newMessage.user != client {
				fmt.Fprint(client.conn, newMessage.msg)
				fmt.Printf("\n")
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
		ServerChat.channel <- FormatTextMessage(newMessage)
	}
	// start messages
	// propmt the name
	// upload history
	// add to chat.Clients
	// wait and get the messages
}

func (ServerChat *Chat) NameValidator(name string) error {
	errEnd := " try again: "
	if len(name) > 20 {
		return errors.New("length of name must contain max 20 letters." + errEnd)
	}
	for _, v := range name {
		if !((v >= 'A' && v <= 'Z') || (v >= 'a' && v <= 'z')) {
			return errors.New("name has not-valid characters. allowed[a-zA-Z]." + errEnd)
		}
	}
	for _, v := range ServerChat.clients {
		if v.name == name {
			return errors.New("name you entered already exists." + errEnd)
		}
	}
	return nil
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

	return newClient
}

func FormatTextMessage(msg Message) Message {
	return Message{msg.user, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]: %s\n", msg.user.name, msg.msg)}
}

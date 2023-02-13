package server

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/helpers"
	"sync"
	"time"
)

type Chat struct {
	clients []Client
	channel chan Message
	history []Message
	mutex   sync.Mutex
}

type Client struct {
	name string
	conn net.Conn
}

type Message struct {
	user Client
	msg  string
}

func (ServerChat *Chat) ProcessMessages() {
	var newMessage Message

	for {
		newMessage = <-ServerChat.channel
		ServerChat.mutex.Lock()
		for _, client := range ServerChat.clients {
			if newMessage.user != client {
				fmt.Fprint(client.conn, "\n"+newMessage.msg)
				ServerChat.history = append(ServerChat.history, Message{client, ("\n" + newMessage.msg)})
				ServerChat.PrintBaseMessage(client)
			}
		}
		ServerChat.mutex.Unlock()
	}
}

func (ServerChat *Chat) ProcessClient(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	helpers.PrintStarterMessage(conn)

	var name string
	for scanner.Scan() {
		name = scanner.Text()
		ok, err := ServerChat.AddNameValidate(name)
		if !(ok) {
			fmt.Fprint(conn, err)
			fmt.Fprint(conn, "[ENTER YOUR NAME]:")
		} else {
			break
		}

	}
	ServerChat.mutex.Lock()
	ServerChat.restoreHistory(conn)
	ServerChat.mutex.Unlock()

	newClient := ServerChat.newClientAdd(conn, name)

	var newMessage Message
	for scanner.Scan() {
		newMessage.user = newClient
		newMessage.msg = scanner.Text()
		if len(newMessage.msg) == 0 {
			ServerChat.PrintBaseMessage(newClient)
			continue
		}
		ServerChat.PrintBaseMessage(newClient)

		ServerChat.mutex.Lock()
		ServerChat.history = append(ServerChat.history, Message{newClient, GetRegularTextMessage(newMessage)})
		ServerChat.mutex.Unlock()

		ServerChat.channel <- Message{newClient, GetRegularTextMessage(newMessage)}
	}

	ServerChat.mutex.Lock()
	ServerChat.history = append(ServerChat.history, Message{newClient, GetClientDeleteMessage(newClient)})
	ServerChat.mutex.Unlock()

	ServerChat.channel <- Message{newClient, GetClientDeleteMessage(newClient)}

	ServerChat.DeleteClient(newClient)
}

func (ServerChat *Chat) AddNameValidate(name string) (bool, string) {
	if len(name) < 3 {
		return false, "Your name should consist at least 3 symbol. Try again\n"
	}
	for _, v := range ServerChat.clients {
		if v.name == name {
			return false, "This name is already exits. Try another\n"
		}
	}
	return true, ""
}

func (ServerChat *Chat) newClientAdd(conn net.Conn, name string) Client {
	var newClient Client
	newClient.conn = conn
	newClient.name = name

	ServerChat.mutex.Lock()
	ServerChat.clients = append(ServerChat.clients, newClient)
	ServerChat.history = append(ServerChat.history, Message{newClient, GetClientAddMessage(newClient)})
	ServerChat.mutex.Unlock()

	ServerChat.channel <- Message{newClient, GetClientAddMessage(newClient)}

	ServerChat.PrintBaseMessage(newClient)

	return newClient
}

func (ServerChat *Chat) DeleteClient(client Client) {
	ServerChat.mutex.Lock()
	for i, v := range ServerChat.clients {
		if v == client {
			ServerChat.clients = append(ServerChat.clients[:i], ServerChat.clients[i+1:]...)
		}
	}
	ServerChat.mutex.Unlock()
}

func GetRegularTextMessage(msg Message) string {
	return fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]: %s\n", msg.user.name, msg.msg)
}

func (ServerChat *Chat) PrintBaseMessage(client Client) {
	fmt.Fprint(client.conn, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]:", client.name))
}

func GetClientAddMessage(client Client) string {
	return fmt.Sprintf("%s has joined our chat...\n", client.name)
}

func GetClientDeleteMessage(client Client) string {
	return fmt.Sprintf("%s has left our chat...\n", client.name)
}

func (ServerChat *Chat) restoreHistory(conn net.Conn) {
	history := ServerChat.history

	for _, v := range history {
		fmt.Fprint(conn, v.msg)
	}
}

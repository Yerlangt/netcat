package main

import (
	"net-cat/helpers"
	"net-cat/server"
	"os"
)

func main() {
	serverPort := helpers.CheckArgs(os.Args)
	if len(serverPort) != 0 {
		server.StartServer(serverPort)
	}
}

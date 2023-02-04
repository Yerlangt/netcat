package helpers

import (
	"fmt"
	"os"
)

func CheckArgs(args []string) string {
	serverPort := ""
	if len(args) == 2 {
		serverPort = os.Args[1]
	} else if len(args) == 1 {
		serverPort = "8989"
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return ""
	}
	return serverPort
}

func StarterMessage() string {
	msg := "Welcome to TCP-Chat!\n"
	msg += "         _nnnn_\n"
	msg += "	dGGGGMMb\n"
	msg += "       @p~qp~~qMb\n"
	msg += "       M|@||@) M|\n"
	msg += "       @,----.JM|\n"
	msg += "      JS^\\__/  qKL\n"
	msg += "     dZP        qKRb\n"
	msg += "    dZP          qKKb\n"
	msg += "   fZP            SMMb\n"
	msg += "   HZM            MMMM\n"
	msg += "   FqM            MMMM\n"
	msg += " __| \".        |\\dS\"qML\n"
	msg += " |    `.       | `' \\Zq\n"
	msg += "_)      \\.___.,|     .'\n"
	msg += "\\____   )MMMMMP|   .'\n"
	msg += "     `-'       `--'\n"
	msg += "[ENTER YOUR NAME]:"
	return msg
}

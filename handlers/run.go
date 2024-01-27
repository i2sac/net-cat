package handlers

import "fmt"

func Exec(args []string) {
	if lenArgs := len(args); lenArgs <= 3 {
		port := "8989"
		if lenArgs == 1 || (lenArgs == 2 && IsPort(args[1])) {
			if lenArgs == 2 {
				port = args[1]
			}
			CreateServer(port)
		} else if lenArgs == 3 && IsIP(args[1]) && IsPort(args[2]) {
			ip, port := args[1], args[2]
			ConnectToServer(ip, port)
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}

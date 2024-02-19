package handlers

import (
	"fmt"
	"log"
)

func Exec(args []string) {
	if lenArgs := len(args); lenArgs <= 3 {
		port := "8989"
		if lenArgs == 1 || (lenArgs == 2 && IsPort(args[1])) { // If there is only port given : Create server
			if lenArgs == 2 {
				port = args[1]
			}
			NetCatServer = *NewServer("localhost:"+port)
			err := NetCatServer.Start()
			if err != nil {
				log.Fatal(err)				
			}
		} else if lenArgs == 3 && IsIP(args[1]) && IsPort(args[2]) { // If port and address given : Connect to server
			ip, port := args[1], args[2]
			NetCatServer.ConnectNewUser(ip, port)
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}

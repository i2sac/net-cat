package handlers

import (
	"fmt"
	"log"
)

func Exec(args []string) {
	if lenArgs := len(args); lenArgs == 1 || (lenArgs == 2 && IsPort(args[1])) {
		port := "8989"

		if lenArgs == 2 {
			port = args[1]
		}
		NetCatServer = *NewServer(Host + port)
		err := NetCatServer.Start()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}

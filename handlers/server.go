package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func CreateServer(port string) {
	fmt.Println("Listening on the port :" + port)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
	}
}

func ConnectToServer(ip, port string) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var name string
	welcomeText, err := os.ReadFile("welcome-text.txt")
	if err != nil {
		fmt.Println("Don't delete or rename \033[31mwelcome-text.txt\033[00m file")
		return
	}
	fmt.Print(string(welcomeText))
	fmt.Scanln(&name)
	if len(name) == 0 {
		fmt.Println("Enter a valid name.")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + name + "]:")
		reader.ReadString('\n')
	}
}

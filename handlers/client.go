package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var ClientName string

type Msg struct {
	Type   string `json:"Type"`
	Author string `json:"Author"`
	Text   string `json:"Text"`
	Date   string `json:"Date"`
}

func (s *Server) ConnectNewUser(ip, port string) {
	// Connect to server
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Unable to connect to " + ip + ":" + port)
		return
	}
	defer conn.Close()

	// Read message from server
	fmt.Print(ReadConnMsg(conn))

	// User login
	_, err = fmt.Scanln(&ClientName)
	LogError(err)

	conn.Write([]byte(ClientName))

	// Send message to server
	s.SendMsg(conn)
}

func (s *Server) UserMessages(conn net.Conn) {
	for {
		// Read message from server
		if newMsg := ReadConnMsg(conn); len(newMsg) > 0 {
			var txt Msg

			err := json.Unmarshal([]byte(newMsg), &txt)
			LogError(err)

			if txt.Type == "notif" {
				fmt.Print("\n" + txt.Text)
			} else {
				fmt.Print("\n" + UserMsgDate(txt.Author, txt.Date) + txt.Text)
			}
		} else {
			continue
		}
	}
}

func (s *Server) SendMsg(conn net.Conn) {
	// Listen to new message
	go s.UserMessages(conn)

	for {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(UserMsgDate(ClientName, timeStr))
		msg, err := reader.ReadString('\n')
		LogError(err)

		req, err := json.Marshal(Msg{"msg", ClientName, strings.ReplaceAll(msg, "\n", ""), timeStr})
		LogError(err)
		conn.Write(req)
	}
}

func UserMsgDate(name, timeStr string) string {
	return "[" + timeStr + "][" + name + "]:"
}

func ReadConnMsg(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err == io.EOF {
		fmt.Println("\nServer stopped")
		os.Exit(0)
	}
	LogError(err)
	return string(buffer[:n])
}

func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

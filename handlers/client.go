package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

var ClientName string

type Msg struct {
	Author string `json:"Author"`
	Text   string `json:"Text"`
	Date   string `json:Date`
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
	if _, err := fmt.Scanln(&ClientName); err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte(ClientName))

	// Send message to server
	s.SendMsg(conn)

	// Listen to new message
	s.UserMessages(conn)
}

func (s *Server) UserMessages(conn net.Conn) {
	for {
		// Read message from server
		if newMsg := ReadConnMsg(conn); len(newMsg) > 0 {
			var txt Msg
			if err := json.Unmarshal([]byte(newMsg), &txt); err != nil {
				log.Fatal(err)
			}
			fmt.Println(UserMsgDate(txt.Author, txt.Date) + txt.Text)
		} else {
			return
		}
	}
}

func (s *Server) SendMsg(conn net.Conn) {
	for {
		timeStr:=time.Now().Format("2006-01-02 15:04:05")
		fmt.Print(UserMsgDate(ClientName, timeStr))
		var msg string
		fmt.Scanln(&msg)
		req, err := json.Marshal(Msg{ClientName, msg, timeStr})
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(req)
	}
}

func UserMsgDate(name, timeStr string) string {
	return "[" + timeStr + "][" + name + "]:"
}

func ReadConnMsg(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	return string(buffer[:n])
}

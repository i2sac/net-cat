package handlers

import (
	"net"
	"os"
	"time"
)

type Msg struct {
	Type   string `json:"Type"`
	Author string `json:"Author"`
	Text   string `json:"Text"`
	Date   string `json:"Date"`
}

var MsgLog []Msg

func (s *Server) MsgLoop(username string, conn net.Conn) {
	for {
		msg := s.ShowMsgField(username, conn)

		if len(msg) > 0 {
			s.BroadcastMsg(msg, "msg", username)
			LogMsg("msg", username, msg)
		}
	}
}

func (s *Server) ShowMsgField(username string, conn net.Conn) string {
	if s.ToClient(UserMsgDate(username), conn) == "OK" {
		return s.ReadMsgResponse(conn)
	}
	return ""
}

func (s *Server) ToClient(msg string, conn net.Conn) string {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return "STOP"
	}
	return "OK"
}

func LogMsg(typeMsg, author, msg string) {
	var formatedMsg string
	if typeMsg == "msg" {
		formatedMsg = UserMsgDate(author) + msg
	}

	file, err := os.OpenFile("msglogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	LogError(err)

	defer file.Close()

	_, err = file.WriteString(formatedMsg + "\n")
	LogError(err)

	// Add message in Msg slice
	Colorize(&msg, typeMsg)
	MsgLog = append(MsgLog, Msg{typeMsg, author, msg, time.Now().Format("2006-01-02 15:04:05")})
}

func UserMsgDate(name string) string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + name + "]:"
}

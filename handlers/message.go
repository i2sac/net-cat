package handlers

import (
	"net"
	"os"
	"strings"
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
		res := s.ReadMsgResponse(conn)
		if len(res) > LimitChars {
			res = res[:LimitChars]
		}
		res = strings.ReplaceAll(res, "\n", "")
		res = strings.ReplaceAll(res, "\t", "")
		res = strings.ReplaceAll(res, "\r", "")
		res = strings.ReplaceAll(res, "\v", "")
		return res
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
		formatedMsg += UserMsgDate(author)
	}
	formatedMsg += msg

	file, err := os.OpenFile("msglogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	LogError(err)

	defer file.Close()

	_, err = file.WriteString(formatedMsg + "\n")
	LogError(err)

	// Add message in Msg slice
	Colorize(&msg, typeMsg)
	MsgLog = append(MsgLog, Msg{typeMsg, author, msg, time.Now().Format("2006-01-02 15:04:05")})
}

func (s *Server) BroadcastMsg(msg string, msgType, excluded string) {
	for conn, usr := range s.clients {
		if usr != excluded {
			conn.Write([]byte(FormatInsert(msg, msgType, excluded)))
		}
	}
}

func MsgLogsToText(logs []Msg) string {
	var txt string
	for _, msg := range logs {
		if msg.Type == "msg" {
			txt += Blue + UserMsgDate(msg.Author) + ColorAnsiEnd
		}
		txt += msg.Text
		if msg.Type == "error" || msg.Type == "notif" {
			txt += "\n"
		}
	}
	return txt
}

func UserMsgDate(name string) string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + name + "]:"
}

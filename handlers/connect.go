package handlers

import (
	"net"
	"os"
	"strings"
	"time"
)

func (s *Server) ConnectUser(conn net.Conn) {
	// Display welcome text
	welcomeText, err := os.ReadFile("welcome-text.txt")
	if err != nil {
		s.ToClient("don't delete or rename \033[31mwelcome-text.txt\033[00m file", conn)
	}
	_, err = conn.Write(welcomeText)
	LogError(err)

	var loginError, username string
	var loginSuccess bool
	for {
		state, res := s.ShowLoginField(loginError, conn)
		if state == "OK" || state == "STOP" {
			loginSuccess = state == "OK"
			username = res
			break
		} else {
			loginError = state
		}
	}

	if loginSuccess {
		s.AddClient(conn, username)
	}
}

func (s *Server) ShowLoginField(err string, conn net.Conn) (string, string) {
	if s.ToClient(err+"[ENTER YOUR NAME]:", conn) == "OK" {
		return s.ReadLoginResponse(conn)
	}
	return "STOP", ""
}

func (s *Server) ToClient(msg string, conn net.Conn) string {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return "STOP"
	}
	return "OK"
}

func (s *Server) ReadLoginResponse(conn net.Conn) (string, string) {
	res := strings.ReplaceAll(s.ReadFromConn(conn), "\n", "")

	if len(res) == 0 {
		return "\rEmpty username !\n", ""
	} else if !IsAlphaNumeric(res) {
		return "\rThe username should be alphanumeric !\nEx: AlphaZero345\n", ""
	}
	return "OK", res
}

func (s *Server) ReadFromConn(conn net.Conn) string {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return ""
	}

	return string(buf[:n])
}

func (s *Server) AddClient(conn net.Conn, name string) {
	if !ExistingUsers[name] && len(s.clients) < maxUsers {
		s.clients[conn] = name     // Save client
		ExistingUsers[name] = true // Mark client as existing

		// Send message to client
		s.MsgToClient("notif", name+" has joined our chat...", time.Now().Format("2006-01-02 15:04:05"), conn)

		// Send logs
		if len(MsgLog) > 0 {
			s.MsgToClient("logs", "Read the log file", time.Now().Format("2006-01-02 15:04:05"), conn)
		}
	} else if ExistingUsers[name] {
		s.MsgToClient("error", "That username already exists.", time.Now().Format("2006-01-02 15:04:05"), conn)
	} else if len(s.clients) == 8 {
		s.MsgToClient("error", "Max number of users reached.", time.Now().Format("2006-01-02 15:04:05"), conn)
	}
}

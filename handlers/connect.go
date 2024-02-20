package handlers

import (
	"fmt"
	"net"
	"os"
)

func (s *Server) ConnectUser(conn net.Conn) {
	// Display welcome text
	welcomeText, err := os.ReadFile("welcome-text.txt")
	if err != nil {
		s.ToClient("don't delete or rename \033[31mwelcome-text.txt\033[00m file", conn)
	}
	_, err = conn.Write(welcomeText)
	LogError(err)

	var loginRes, username string

	s.LoginLoop(&username, &loginRes, conn)

	s.ShowMsgField(username, conn)
}

func (s *Server) ShowLoginField(err string, conn net.Conn) (string, string) {
	if s.ToClient(err+"[ENTER YOUR NAME]:", conn) == "OK" {
		return s.ReadLoginResponse(conn)
	}
	return "STOP", ""
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

func (s *Server) AddClient(conn net.Conn, name string) (bool, string) {
	if !ExistingUsers[name] && len(s.clients) < maxUsers {
		s.clients[conn] = name     // Save client
		ExistingUsers[name] = true // Mark client as existing

		res := name + " has joined our chat..."
		Colorize(&res, "notif")
		return true, res + "\n"

		/*
			// Send message to client
			s.MsgToClient("notif", name+" has joined our chat...", time.Now().Format("2006-01-02 15:04:05"), conn)

			// Send logs
			if len(MsgLog) > 0 {
				s.MsgToClient("logs", "Read the log file", time.Now().Format("2006-01-02 15:04:05"), conn)
			}
		*/
	} else if ExistingUsers[name] {
		// s.MsgToClient("error", "That username already exists.", time.Now().Format("2006-01-02 15:04:05"), conn)
		res := "That username already exists."
		Colorize(&res, "error")
		return false, res + "\n"
	} else if len(s.clients) == 8 {
		// s.MsgToClient("error", "Max number of users reached.", time.Now().Format("2006-01-02 15:04:05"), conn)
		res := "Max number of users reached."
		Colorize(&res, "error")
		return false, res + "\n"
	}
	return false, "Username error"
}

func (s *Server) LoginLoop(username, loginRes *string, conn net.Conn) {
	for {
		state, res := s.ShowLoginField(*loginRes, conn)
		if state == "OK" || state == "STOP" {
			*username = res

			success, response := s.AddClient(conn, *username)
			*loginRes = response

			if success {
				break
			}
		} else {
			*loginRes = state
		}
	}

	if len(*username) > 0 {
		fmt.Print(*loginRes)
		s.BroadcastMsg(*loginRes, "notif", *username)
	}
}

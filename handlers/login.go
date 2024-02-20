package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

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
		logs, err := json.Marshal(MsgLog)
		LogError(err)
		err = os.WriteFile("msglogs.json", logs, 0755)
		LogError(err)
		
		Colorize(loginRes, "notif")
		fmt.Print(*loginRes)
		s.BroadcastMsg(*loginRes, "notif", *username)

		if len(MsgLog) > 0 {
			s.ToClient(FormatInsert("", "logs", *username), conn)
		}

		LogMsg("notif", *username, strings.ReplaceAll(*loginRes, "\n", ""))
	}
}

func (s *Server) ShowLoginField(err string, conn net.Conn) (string, string) {
	if s.ToClient(err+"[ENTER YOUR NAME]:", conn) == "OK" {
		return s.ReadLoginResponse(conn)
	}
	return "STOP", ""
}

func (s *Server) AddClient(conn net.Conn, name string) (bool, string) {
	if !ExistingUsers[name] && len(s.clients) < maxUsers {
		s.clients[conn] = name     // Save client
		ExistingUsers[name] = true // Mark client as existing

		res := name + " has joined our chat..."
		return true, res + "\n"
	} else if ExistingUsers[name] {
		res := "That username already exists."
		Colorize(&res, "error")
		return false, res + "\n"
	} else if len(s.clients) == 8 {
		res := "Max number of users reached."
		Colorize(&res, "error")
		return false, res + "\n"
	}
	return false, "Username error"
}

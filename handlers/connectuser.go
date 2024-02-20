package handlers

import (
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

	go s.MsgLoop(username, conn)
}


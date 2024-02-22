package handlers

import (
	"net"
	"strings"
)

func (s *Server) ReadLoginResponse(conn net.Conn) (string, string) {
	res := strings.ReplaceAll(s.ReadFromConn(conn), "\n", "")

	if len(res) == 0 {
		txt := "Empty username !"
		Colorize(&txt, "error")
		return "\r" + txt + "\n", ""
	} else if !IsAlphaNumeric(res) {
		txt := "The username should be alphanumeric !\nEx: AlphaZero345"
		Colorize(&txt, "error")
		return "\r" + txt + "\n", ""
	} else if len(res) > MaxUsernameLength {
		txt := "Username too long!"
		Colorize(&txt, "error")
		return "\r" + txt + "\n", ""
	}
	return "OK", res
}

func (s *Server) ReadMsgResponse(conn net.Conn) string {
	res := strings.ReplaceAll(s.ReadFromConn(conn), "\n", "")

	if len(res) == 0 || !IsReadable(res) {
		return ""
	}
	return res
}

func (s *Server) ReadFromConn(conn net.Conn) string {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		s.CloseConnection(conn, s.clients[conn])
		return ""
	}

	return string(buf[:n])
}

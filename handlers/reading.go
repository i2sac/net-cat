package handlers

import (
	"net"
	"strings"
)

func (s *Server) ReadLoginResponse(conn net.Conn) (string, string) {
	res := strings.ReplaceAll(s.ReadFromConn(conn), "\n", "")

	if len(res) == 0 {
		return "\rEmpty username !\n", ""
	} else if !IsAlphaNumeric(res) {
		return "\rThe username should be alphanumeric !\nEx: AlphaZero345\n", ""
	} else if len(res) > MaxUsernameLength {
		return "\rUsername too long!\n", ""
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

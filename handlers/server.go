package handlers

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan []byte
	clients    map[net.Conn]string
}

var NetCatServer Server
var ExistingUsers = make(map[string]bool)
var MaxUsers = 10
var LimitChars = 100
var MaxUsernameLength = 30

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan []byte, 10),
		clients:    make(map[net.Conn]string),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	fmt.Println("Listening on the port :", s.listenAddr[len("localhost:"):])

	// Create log file
	os.WriteFile("msglogs.log", []byte(""), 0755)

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go s.ConnectUser(conn)

	}
}

func (s *Server) CloseConnection(conn net.Conn, client string) {
	delete(s.clients, conn)
	ExistingUsers[client] = false
	conn.Close()

	if len(client) > 0 {
		res := client + " has left our chat..."

		LogMsg("notif", client, res)
		fmt.Println(Orange + res + ColorAnsiEnd)

		s.BroadcastMsg(res, "notif", client)
	}
}

var Orange = ColorAnsiStart(255, 94, 0)
var Red = ColorAnsiStart(255, 0, 0)
var Blue = ColorAnsiStart(0, 60, 255)

func Colorize(msg *string, typeMsg string) {
	switch typeMsg {
	case "notif":
		*msg = Orange + *msg + ColorAnsiEnd
	case "error":
		*msg = Red + *msg + ColorAnsiEnd
	case "msg":
		*msg = Blue + *msg + ColorAnsiEnd + "\n"
	}
}

// Function that creates the escape color string for the given RGB color
func ColorAnsiStart(R, G, B int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", R, G, B)
}

// Color string to reset string color
var ColorAnsiEnd = "\033[0m"

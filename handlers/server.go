package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	host, port string
	listener   net.Listener
	clients    map[string]net.Conn
}

var NetCatServer = Server{clients: make(map[string]net.Conn)}
var ExistingUsers = make(map[string]bool)

const MaxUsers = 8

func (s *Server) CreateServer(port string) {
	s.port = port
	fmt.Println("Listening on the port :" + port)
	
	ln, err := net.Listen("tcp", s.host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	s.listener = ln
	
}

func (s *Server) AcceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.HandleConnection(conn)
		// Find new sent messages
		go s.RefreshMsg()
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	// Display welcome text
	welcomeText, err := os.ReadFile("welcome-text.txt")
	if err != nil {
		fmt.Println("Don't delete or rename \033[31mwelcome-text.txt\033[00m file")
		return
	}
	_, err = conn.Write(welcomeText)
	if err != nil {
		log.Fatal(err)
	}

	// Read client username
	reader := bufio.NewReader(conn)
	usr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Add new client
	s.AddClient(conn, usr)
}

func (s *Server) AddClient(conn net.Conn, name string) {
	if !ExistingUsers[name] && len(s.clients) < 8 {
		s.clients[name] = conn
		s.BroadcastMsg(name+" has joined our chat...", name)
	} else {
		conn.Close()
	}
}

func (s *Server) BroadcastMsg(msg string, excluded string) {
	for usr, conn := range s.clients {
		if usr != excluded {
			conn.Write([]byte(msg))
		}
	}
}

func (s *Server) RefreshMsg() {
	// Broadcast new messages
	for {
		for usr, conn := range s.clients {
			s.BroadcastMsg(ReadConnMsg(conn), usr)
		}
	}
}

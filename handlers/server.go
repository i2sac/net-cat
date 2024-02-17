package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
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
var maxUsers = 10

var MsgLog []Msg

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

		go s.ShowLogin(conn)

		go s.readLoop(conn)

		go s.printLoop(conn)

	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	msgCount := 0
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				name := s.clients[conn]

				newMsg := Msg{"notif", name, name + " has left our chat...\n", time.Now().Format("2006-01-02 15:04:05")}
				req, err := json.Marshal(newMsg)
				LogError(err)
				if len(newMsg.Author) > 0 {
					s.msgch <- req
					s.closeConnection(conn, newMsg.Author)
				}
			} else {
				fmt.Println("read error:", err)
			}
			break
		}

		msg := buf[:n]
		if msgCount == 0 { // First Message = Username
			msgTxt := string(msg)
			s.AddClient(conn, msgTxt)
		} else {
			s.msgch <- msg
		}
		msgCount++
	}
}

func (s *Server) closeConnection(conn net.Conn, client string) {
	delete(s.clients, conn)
	ExistingUsers[client] = false
	conn.Close()
}

func (s *Server) printLoop(conn net.Conn) {
	for msg := range s.msgch {
		newMSG := Msg{}
		err := json.Unmarshal(msg, &newMSG)
		LogError(err)
		MsgLog = append(MsgLog, newMSG)

		if newMSG.Type == "msg" && len(newMSG.Text) > 0 {
			s.BroadcastMsg(msg, newMSG.Author)
		} else {
			fmt.Print(newMSG.Text)
			if newMSG.Text == "error" {
				conn.Write([]byte(msg))
				conn.Close()
			}
		}
	}
}

func (s *Server) ShowLogin(conn net.Conn) error {
	// Display welcome text
	welcomeText, err := os.ReadFile("welcome-text.txt")
	if err != nil {
		return errors.New("don't delete or rename \033[31mwelcome-text.txt\033[00m file")
	}
	_, err = conn.Write(welcomeText)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *Server) AddClient(conn net.Conn, name string) {
	if !ExistingUsers[name] && len(s.clients) < maxUsers {
		s.clients[conn] = name     // Save client
		ExistingUsers[name] = true // Mark client as existing

		// Send message to client
		s.MsgToClient("notif", name+" has joined our chat...\n", time.Now().Format("2006-01-02 15:04:05"), conn)
	} else if ExistingUsers[name] {
		s.MsgToClient("error", "That username already exists.\n", time.Now().Format("2006-01-02 15:04:05"), conn)
	} else if len(s.clients) == 8 {
		s.MsgToClient("error", "Max number of users reached.\n", time.Now().Format("2006-01-02 15:04:05"), conn)
	}
}

func (s *Server) BroadcastMsg(msg []byte, excluded string) {
	for conn, usr := range s.clients {
		if usr != excluded {
			conn.Write([]byte(msg))
		}
	}
}

func MsgLogToText() string {
	var txt string
	for _, msg := range MsgLog {
		txt += UserMsgDate(msg.Author, msg.Date) + msg.Text + "\n"
	}
	return txt
}

func (s *Server) MsgToClient(typeMsg, txt, t string, conn net.Conn) {
	name := s.clients[conn]
	newMsg := Msg{typeMsg, name, txt, time.Now().Format("2006-01-02 15:04:05")}
	req, err := json.Marshal(newMsg)
	LogError(err)

	if typeMsg == "error" {
		conn.Write(req)
	} else {
		s.msgch <- req
	}
}

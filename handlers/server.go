package handlers

import (
	"fmt"
	"net"
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

		go s.ConnectUser(conn)

		// go s.readLoop(conn)

		// go s.printLoop(conn)

	}
}

/*
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	msgCount := 0
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				name := s.clients[conn]

				newMsg := Msg{"notif", name, name + " has left our chat...", time.Now().Format("2006-01-02 15:04:05")}
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
*/

func (s *Server) CloseConnection(conn net.Conn, client string) {
	delete(s.clients, conn)
	ExistingUsers[client] = false
	conn.Close()

	res := client + " has left our chat..."
	fmt.Println(Orange + res + ColorAnsiEnd)

	s.BroadcastMsg(res, "notif", client)
}

/*
func (s *Server) printLoop(conn net.Conn) {
	for msg := range s.msgch {
		newMSG := Msg{}
		err := json.Unmarshal(msg, &newMSG)
		LogError(err)

		Colorize(&newMSG)
		MsgLog = append(MsgLog, newMSG)

		if len(newMSG.Text) > 0 {
			if newMSG.Type != "msg" {
				fmt.Println(newMSG.Text)
			}
			s.BroadcastMsg(EncodeMsg(newMSG), newMSG.Author)
		}
	}
}
*/

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

/*
func (s *Server) MsgToClient(typeMsg, txt, t string, conn net.Conn) {
	name := s.clients[conn]
	newMsg := Msg{typeMsg, name, txt, time.Now().Format("2006-01-02 15:04:05")}

	Colorize(&newMsg)

	req := EncodeMsg(newMsg)

	if typeMsg == "error" {
		fmt.Println(newMsg.Text)
		state := s.ToClient(txt, conn)
		fmt.Println(state)
	} else if typeMsg == "logs" {
		logs, err := json.Marshal(MsgLog)
		LogError(err)
		err = os.WriteFile("msglogs.json", logs, 0755)
		LogError(err)
		conn.Write(req)
	} else {
		s.msgch <- req
	}
}
*/

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

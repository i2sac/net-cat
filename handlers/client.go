package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var ClientName string
var MsgLineLength int

type Msg struct {
	Type   string `json:"Type"`
	Author string `json:"Author"`
	Text   string `json:"Text"`
	Date   string `json:"Date"`
}

func (s *Server) ConnectNewUser(ip, port string) {
	// Connect to server
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Unable to connect to " + ip + ":" + port)
		return
	}
	defer conn.Close()

	// Read message from server
	fmt.Println(ReadConnMsg(conn))

	AskUserName()

	conn.Write([]byte(ClientName))

	// Send message to server
	s.SendMsg(conn)

	go func() {
		for {
			if res := ReadConnMsg(conn); len(res) > 0 {
				fmt.Print("\r\033[K" + DecodeMsg(res).Text)
				os.Exit(0)
			} else {
				break
			}
		}
	}()
}

func (s *Server) UserMessages(conn net.Conn) {
	for {
		// Read message from server
		if newMsg := ReadConnMsg(conn); len(newMsg) > 0 {
			txt := DecodeMsg(newMsg)

			fmt.Print("\r\033[1L") // Insertion à la ligne précédente
			if txt.Type == "notif" {
				fmt.Print(txt.Text)
			} else if txt.Type == "error" {
				fmt.Print(txt.Text)
				os.Exit(0)
			} else {
				fmt.Print(UserMsgDate(txt.Author, txt.Date) + txt.Text + "\033[1B\r")
			}
			fmt.Printf("\033[%dC", MsgLineLength)
		} else {
			continue
		}
	}
}

func (s *Server) SendMsg(conn net.Conn) {
	// Listen to new message
	go s.UserMessages(conn)

	for {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		reader := bufio.NewReader(os.Stdin)

		msgLine := UserMsgDate(ClientName, timeStr)
		MsgLineLength = len(msgLine)

		fmt.Print(msgLine)
		msg, err := reader.ReadString('\n')
		LogError(err)

		req, err := json.Marshal(Msg{"msg", ClientName, strings.ReplaceAll(msg, "\n", ""), timeStr})
		LogError(err)
		conn.Write(req)
	}
}

func UserMsgDate(name, timeStr string) string {
	return "[" + timeStr + "][" + name + "]:"
}

func ReadConnMsg(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err == io.EOF {
		fmt.Println("\nServer stopped")
		os.Exit(0)
	}
	LogError(err)
	return string(buffer[:n])
}

func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if r == ' ' || !((r >= 'a' && r <= 'z') /*miniscules*/ || (r >= 'A' && r <= 'Z') /*majuscules*/ || (r >= '0' && r <= '9') /*chiffres*/) {
			return false
		}
	}
	return true
}

func DecodeMsg(newMsg string) Msg {
	var txt Msg
	err := json.Unmarshal([]byte(newMsg), &txt)
	LogError(err)
	return txt
}

func AskUserName() {
	// User login
	var errTxt string
	for len(ClientName) == 0 {
		fmt.Print(errTxt + "[ENTER YOUR NAME]:")

		reader := bufio.NewReader(os.Stdin)
		ClientName, _ = reader.ReadString('\n')
		ClientName = strings.ReplaceAll(ClientName, "\n", "")

		if len(ClientName) == 0 {
			errTxt = "\rEmpty username !\n"
		} else if !IsAlphaNumeric(ClientName) {
			errTxt = "\rThe username should be alphanumeric !\nEx: AlphaZero345\n"
			ClientName = ""
		}
	}
}

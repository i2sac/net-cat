package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

	// Lecture de l'erreur ou des logs
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

			if txt.Type != "logs" {
				// Enregistrez la position du curseur actuelle
				fmt.Print("\033[s")   // Sauvegarde la position du curseur
				fmt.Println("\033[u") // Rajoute une nouvelle ligne puis déplace le curseur à la colonne initiale sur la même ligne

				// Insérez une nouvelle ligne et imprimez le message
				fmt.Print("\r\033[A\033[1L") // Déplace le curseur au début de la ligne, remonte d'une ligne et insère une ligne vide
				if txt.Type == "notif" {
					fmt.Print(txt.Text)
				} else if txt.Type == "error" {
					fmt.Println(txt.Text)
					os.Exit(0)
				} else if txt.Type == "msg" {
					fmt.Print(Blue + UserMsgDate(txt.Author) + ColorAnsiEnd + txt.Text)
				}

				// Restaure la position du curseur
				fmt.Print("\033[1B\r")
				fmt.Print("\033[u\033[1B") // Restaure la position du curseur
			} else {
				fmt.Print("\033[s")
				fmt.Println("\033[A")

				logsRaw, err := os.ReadFile("msglogs.json")
				LogError(err)

				var logs []Msg
				err = json.Unmarshal(logsRaw, &logs)
				LogError(err)

				logsText := MsgLogsToText(logs)

				nbLines := strings.Count(logsText, "\n")
				fmt.Print("\033[u")
				fmt.Print("\033[s")
				fmt.Printf("\r\033[A\n\033[%dL", nbLines)
				fmt.Print(logsText)
				fmt.Printf("\033[u\033[%dB", nbLines)

				msgLine := UserMsgDate(ClientName)
				fmt.Print("\r" + msgLine)
			}
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

		msgLine := UserMsgDate(ClientName)
		MsgLineLength = len(msgLine)

		fmt.Print(msgLine)
		msg, err := reader.ReadString('\n')
		LogError(err)
		msg = strings.ReplaceAll(msg, "\n", "")
		if len(msg) > 0 && IsReadable(msg) {
			conn.Write(EncodeMsg(Msg{"msg", ClientName, msg, timeStr}))
		}
	}
}

func UserMsgDate(name string) string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + name + "]:"
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

func AskUserName() {
	// User login
	var errTxt string
	for len(ClientName) == 0 {
		fmt.Print(errTxt + "[ENTER YOUR NAME]:")

		reader := bufio.NewReader(os.Stdin)
		var err error
		ClientName, err = reader.ReadString('\n')
		LogError(err)
		ClientName = strings.ReplaceAll(ClientName, "\n", "")

		if len(ClientName) == 0 {
			errTxt = "\rEmpty username !\n"
		} else if !IsAlphaNumeric(ClientName) {
			errTxt = "\rThe username should be alphanumeric !\nEx: AlphaZero345\n"
			ClientName = ""
		}
	}
}

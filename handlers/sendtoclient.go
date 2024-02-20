package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func FormatInsert(msg, msgType, username string) string {
	var res string
	if msgType != "logs" {
		// Enregistrez la position du curseur actuelle
		res = "\033[s"    // Sauvegarde la position du curseur
		res += "\033[u\n" // Rajoute une nouvelle ligne puis déplace le curseur à la colonne initiale sur la même ligne

		// Insérez une nouvelle ligne et imprimez le message
		res += "\r\033[A\033[1L" // Déplace le curseur au début de la ligne, remonte d'une ligne et insère une ligne vide
		if msgType == "msg" {
			msg = UserMsgDate(username) + msg
		}

		// Ajouter le message
		Colorize(&msg, msgType)
		res += msg

		res += "\033[1B\r"     // Restaure la position du curseur
		res += "\033[u\033[1B" // Restaure la position du curseur
	} else {
		res = "\033[s"
		res += "\033[A\n"

		logsRaw, err := os.ReadFile("msglogs.json")
		LogError(err)

		var logs []Msg
		err = json.Unmarshal(logsRaw, &logs)
		LogError(err)

		logsText := MsgLogsToText(logs)

		nbLines := strings.Count(logsText, "\n")
		res += "\033[u"
		res += "\033[s"
		res += fmt.Sprintf("\r\033[A\n\033[%dL", nbLines)
		res += logsText
		res += fmt.Sprintf("\033[u\033[%dB", nbLines)

		// res += "\r" + UserMsgDate(username)
	}
	return res
}

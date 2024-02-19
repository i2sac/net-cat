package handlers

import (
	"encoding/json"
	"log"
)

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

func IsReadable(s string) bool {
	for _, r := range s {
		if r >= 0 && r < ' ' || r == 127 {
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

func EncodeMsg(msg Msg) []byte {
	res, err := json.Marshal(msg)
	LogError(err)
	return res
}

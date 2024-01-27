package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) <= 3 {
		port := "8989"
		if len(os.Args) == 2 && IsPort(os.Args[1]) {
			port = os.Args[1]
			fmt.Println("Listening on the port :" + port)
		} else if len(os.Args) == 3 && IsPort(os.Args[1]) && IsIP(os.Args[2]) {
			var name string
			welcomeText, err := os.ReadFile("welcome-text.txt")
			if err != nil {
				fmt.Println("Don't delete or rename \033[31mwelcome-text.txt\033[00m file")
				return
			}
			fmt.Print(string(welcomeText))
			fmt.Scanln(&name)
			if len(name) == 0 {
				fmt.Println("Enter a valid name.")
				return
			}
			fmt.Println(name)
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}

func IsPort(s string) bool {
	return len(s) > 0 && !regexp.MustCompile(`\D`).MatchString(s)
}

func IsIP(s string) bool {
	oct := `([1-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])`
	return len(s) > 0 && (regexp.MustCompile(`^`+oct+`.`+oct+`.`+oct+`.`+oct+`$`).MatchString(s) || s == "localhost")
}

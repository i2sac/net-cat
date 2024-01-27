package main

import (
	"net-cat/handlers"
	"os"
)

func main() {
	handlers.Exec(os.Args)
}

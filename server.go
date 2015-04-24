package main

import (
	"chatroom/bootstrap"
	_ "chatroom/routes"
	"log"
	"os"
)

func main() {
	port := "5000"
	arglen := len(os.Args)

	if arglen > 1 {
		port = os.Args[2]
	}

	log.Println("running server on port", port)

	bootstrap.Start(port, func() {})
}

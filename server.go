package main

import (
	"log"
	"os"
	"chatroom/bootstrap"
	_ "chatroom/routes"
)

func main() {
	port := "5000"
	argLen := len(os.Args)

	if argLen > 1 {
		port = os.Args[2]
	}

	log.Println("Running server on port", port)

	bootstrap.Start(port, func() {})
}

package webSocket

import (
	"github.com/go-martini/martini"
	"chatroom/service/webSocket"
)

func Socket(params martini.Params, receiver <-chan *webSocket.Message, sender chan<- *webSocket.Message, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
	return webSocket.Listen(params, receiver, sender, done, disconnect, err)
}

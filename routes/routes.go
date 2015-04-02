package routes

import (
	sockets "github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"chatroom/config"
	"chatroom/controller/home"
	ctrlWebSocket "chatroom/controller/webSocket"
	"chatroom/service/webSocket"
)

type Routes struct{}

func init() {
	config.AppendValue(config.Controller, &Routes{})
}

func (ctn *Routes) SetRouter(m *martini.ClassicMartini) {

	m.Get("/", home.Home)
	m.Get("/socket", sockets.JSON(webSocket.Message{}), ctrlWebSocket.Socket)

}

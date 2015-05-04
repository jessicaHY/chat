package routes

import (
	"chatroom/config"
	"chatroom/controller/ajax"
	"chatroom/controller/home"
	ctrlWebSocket "chatroom/controller/webSocket"
	"chatroom/service/webSocket"
	sockets "github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
)

type Routes struct{}

func init() {
	config.AppendValue(config.Controller, &Routes{})
}

func (ctn *Routes) SetRouter(m *martini.ClassicMartini) {

	m.Get("/", home.Home)
	m.Get("/socket/:roomId", ctrlWebSocket.PreCheck, sockets.JSON(webSocket.ChatMsg{}), ctrlWebSocket.HandlerSocket)
	m.Get("/room/info/:roomId", ajax.RoomInfo)
	m.Post("/room/add", ajax.AddRoom)
	m.Post("/room/edit/:roomId", ajax.EditRoom)
	m.Get("/room/close/:roomId", ajax.CloseRoom)
	m.Get("/room/list/:bookId", ajax.QueryRoom)
	m.Get("/room/buy/:roomId", ajax.BuyRoom)

	m.Post("/room/shutup/add", ajax.AddShutup)
	m.Post("/room/shutup/del", ajax.DelShutup)
}

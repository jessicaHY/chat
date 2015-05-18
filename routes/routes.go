package routes

import (
	"chatroom/config"
	"chatroom/controller/ajax"
	"chatroom/controller/home"
	ctrlWebSocket "chatroom/controller/webSocket"
	"chatroom/service/webSocket"
	sockets "github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"chatroom/helper"
)

type Routes struct{}

func init() {
	config.AppendValue(config.Controller, &Routes{})
}

func (ctn *Routes) SetRouter(m *martini.ClassicMartini) {

	m.Get("/room/", home.Home)
	m.Get("/room/book/:bookId", ajax.GetRoomByBookId)
	m.Get("/room/socket/:roomId", ctrlWebSocket.PreCheck, sockets.JSON(webSocket.ChatMsg{}), ctrlWebSocket.HandlerSocket)
	m.Get("/room/info/:roomId", ajax.RoomInfo)
	m.Post("/room/add", ajax.AddRoom)
	m.Post("/room/edit/:roomId", ajax.EditRoom)
	m.Get("/room/close/:roomId", ajax.CloseRoom)
	m.Get("/room/end/:roomId", ajax.EndRoom)
	m.Get("/room/list/:bookId", ajax.QueryRoom)
	m.Get("/room/buy/:roomId", ajax.BuyRoom)

	m.Get("/room/check/:roomId", func(params martini.Params, rend render.Render, req *http.Request){
		if _, _, err, _ := ctrlWebSocket.UserCheck(params, req); err != helper.NoError {
			rend.JSON(200, helper.Error(err))
		} else {
			rend.JSON(200, helper.Success())
		}

	})

	m.Get("/room/listuser/:roomId", ctrlWebSocket.GetUserList)
	m.Post("/room/shutup/add", ajax.AddShutup)
	m.Post("/room/shutup/del", ajax.DelShutup)

	m.Get("/room/gift/list", ajax.ListGift)
	m.Get("/room/donate/:roomId", ajax.ListDonateByRoom)
	m.Post("/room/donate/:roomId", ajax.DonateRoom)
}

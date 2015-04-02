package webSocket

import (
	"github.com/go-martini/martini"
	"chatroom/service/webSocket"
	"chatroom/helper"
	"chatroom/utils/JSON"
)

func Socket(params martini.Params, receiver <-chan *webSocket.Message, sender chan<- *webSocket.Message, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
	return webSocket.Listen(params, receiver, sender, done, disconnect, err)
}

type Message struct {
	Content string
}

type UserMessage struct {
	Content string
	UserPic string
	Name string
}

func init() {

	webSocket.OnAppend(func(clientLength int) {
		//do sth..
	})

	webSocket.OnEmit("authorSend", func(msg *webSocket.Message) JSON.Type {
		content := &Message{}
		if err := JSON.ParseToStruct(msg.Params, content); err == nil {
			webSocket.BroadCastAll(&webSocket.Message{
			Method: "authorMessage",
			Params: content,
		})
			return helper.Success(JSON.Type{})
		}

		return helper.Error(helper.ParamsError)
	})

	webSocket.OnEmit("userSend", func(msg *webSocket.Message) JSON.Type {

		content := &UserMessage{}
		if err := JSON.ParseToStruct(msg.Params, content); err == nil {
			content.UserPic = "http://img.heiyanimg.com/people/1993937s7.jpg"
			content.Name = "抖S是我"

			webSocket.BroadCastAll(&webSocket.Message{
				Method: "userMessage",
				Params: content,
			})
			return helper.Success(JSON.Type{})
		}

		return helper.Error(helper.ParamsError)
	})
}

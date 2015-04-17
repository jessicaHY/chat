package webSocket

import (
	"chatroom/helper"
	"chatroom/service/models"
	"chatroom/service/redis"
	"chatroom/service/webSocket"
	"chatroom/utils/JSON"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"time"
)

func init() {

	webSocket.OnAppend(func(client *webSocket.SocketClient, r *webSocket.Room) {
		fmt.Println("OnAppend....")
		//发三条信息
		ucount, acount, _ := redis.ZCard(r.RoomId)
		client.UserMsgIndex = ucount
		if acount > FIRST_CONTENT_SIZE { //确定该client对应的作者信息的起始和结束
			client.AuthorStartIndex = acount - FIRST_CONTENT_SIZE
			client.AuthorEndIndex = acount - FIRST_CONTENT_SIZE
		}
		fmt.Println(client)

		msgs, err := GetWebSocketChatMsg(redis.AuthorMessage, r.RoomId, client.AuthorEndIndex, -1)
		if err != nil {
			fmt.Println(err)
		}
		r.SendSelf(client, &webSocket.ChatMsg{Method: "authorMessage", Params: msgs})
		client.AuthorEndIndex += len(msgs)
		fmt.Println(client)
	})

	webSocket.OnRemove(func(userId int) {
		delete(UserInfoMap, userId)
	})
	//作者发信息
	webSocket.OnEmit("authorSend", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {
		fmt.Println("authorSend")
		if client.UserId != r.AuthorId {
			return helper.Error(helper.NoRightError)
		}
		uMsg := &UserMsg{}
		if err := JSON.ParseToStruct(msg.Params, uMsg); err == nil {
			//insert into db
			tMsg, err := models.AddMessage(client.UserId, r.RoomId, models.Chapter, uMsg.Content)
			if err != nil {
				fmt.Println(err)
				return helper.Error(helper.ParamsError)
			}
			uMsg.Id = tMsg.Id
			uMsg.CreateTime = tMsg.CreateTime
			uMsg.Info = GetUserInfo(tMsg.UserId)
			fmt.Println(uMsg)

			//save to redis
			b, err := json.Marshal(uMsg)
			if err != nil {
				fmt.Println(err)
				return helper.Error(helper.ParamsError)
			}
			_, err = redis.ZAddAuthorMsg(r.RoomId, tMsg.Id, string(b))
			if err != nil {
				fmt.Println(err)
				return helper.Error(helper.ParamsError)
			}
			//tell thread to tell everyclient
			fmt.Println(r.ThreadChannel)
			select {
			case r.ThreadChannel <- 1:
				fmt.Println("to ...runMsgTask")
				break
			default:
				fmt.Println("back msg")
				break
			}
			return helper.Success(JSON.Type{})
		}
		return helper.Error(helper.ParamsError)
	})

	//用户发消息
	webSocket.OnEmit("userSend", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {
		fmt.Println("userSend...")
		if client.UserId <= 0 {
			return helper.Error(helper.NoLoginError)
		}
		uMsg := &UserMsg{}
		if err := JSON.ParseToStruct(msg.Params, uMsg); err == nil {
			uMsg.Id = 0
			uMsg.CreateTime = time.Now()
			uMsg.Info = GetUserInfo(client.UserId)
			fmt.Println(uMsg)

			//save to redis
			b, err := json.Marshal(uMsg)
			if err != nil {
				return helper.Error(helper.ParamsError)
			}
			_, err = redis.ZAddUserMsg(r.RoomId, string(b))
			if err != nil {
				return helper.Error(helper.ParamsError)
			}
			//tell thread to tell everyclient
			fmt.Println(r.ThreadChannel)
			select {
			case r.ThreadChannel <- 1:
				glog.Infoln("to ...runMsgTask")
				break
			default:
				glog.Infoln("back msg")
				break
			}
			return helper.Success(JSON.Type{})
		}

		return helper.Error(helper.ParamsError)
	})

	//用户点击获取更多
	webSocket.OnEmit("getMessage", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {
		param := &Param{}
		if err := JSON.ParseToStruct(msg.Params, param); err == nil {
			if param.Next { //获取往后的数据
				//needs to send
				msgs, err := GetWebSocketChatMsg(redis.UserMessage, r.RoomId, client.UserMsgIndex, -1)
				if err != nil {
					return helper.Error(helper.ParamsError)
				}
				r.SendSelf(client, &webSocket.ChatMsg{Method: "userMessage", Params: msgs})
				client.UserMsgIndex += len(msgs)

				msgs, err = GetWebSocketChatMsg(redis.AuthorMessage, r.RoomId, client.AuthorEndIndex, -1)
				if err != nil {
					return helper.Error(helper.ParamsError)
				}
				r.SendSelf(client, &webSocket.ChatMsg{Method: "authorMessage", Params: msgs})
				client.AuthorEndIndex += len(msgs)
			} else { //获取之前的作者数据
				begin := client.AuthorStartIndex - param.Size
				if begin < 0 {
					begin = 0
				}
				msgs, err := GetWebSocketChatMsg(redis.AuthorMessage, r.RoomId, client.AuthorStartIndex, client.AuthorStartIndex)
				if err != nil {
					return helper.Error(helper.ParamsError)
				}
				r.SendSelf(client, &webSocket.ChatMsg{Method: "authorMessage", Params: msgs})
				client.AuthorEndIndex += len(msgs)
			}
			return helper.Success(JSON.Type{})
		}
		return helper.Error(helper.ParamsError)
	})

	//publish chapter
	webSocket.OnEmit("publish", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {
		content := models.PushCompleteRoom(r.RoomId)
		r.SendSelf(client, &webSocket.ChatMsg{Method: "content", Params: content})
		return helper.Success(JSON.Type{})
	})

	redis.OnEmpity(func(roomId int64) bool {
		msgts := models.ListMessage(roomId)
		size := len(msgts)
		args := make(map[int64]string)
		for _, m := range msgts {
			uMsg := &UserMsg{}
			uMsg.Id = m.Id
			uMsg.CreateTime = m.CreateTime
			uMsg.Content = m.Content
			uMsg.Info = GetUserInfo(m.UserId)

			b, err := json.Marshal(uMsg)
			if err != nil {
				fmt.Println(err)
				break
			}
			args[m.Id] = string(b)
		}
		count, err := redis.ZAddAuthorMsgs(roomId, args)
		if err != nil {
			fmt.Println(err)
		}
		return size == count
	})
}

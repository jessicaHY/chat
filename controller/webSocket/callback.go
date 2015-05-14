package webSocket

import (
	"chatroom/helper"
	"chatroom/service/models"
	"chatroom/service/redis"
	"chatroom/service/webSocket"
	"chatroom/utils/JSON"
	"encoding/json"
	"time"
	"log"
	"chatroom/utils/Constants"
)

func NotifyAllClients(r *webSocket.Room) {
	log.Println("NotifyAllClients...")
	select {
	case r.ThreadChannel <- true:
		log.Println("to ...runMsgTask")
	break
	default:
		log.Println("back msg")
	break
	}
}

func init() {

	//该用户第一次连接时调用
	webSocket.BeforeAppend(func(client *webSocket.SocketClient, r *webSocket.Room) {
		if client.UserId <= 0 {
			return
		}
		var userType = Constants.User
		if client.UserId == r.AuthorId {
			userType = Constants.Writer
		}
		uMsg := &UserMsg{0, "", time.Now(), userType, models.Reply, Constants.IsIn, GetUserInfo(client.UserId)}
		log.Println(uMsg)
		//save to redis
		b, err := json.Marshal(uMsg)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = redis.ZAddUserMsg(r.RoomId, string(b))
		NotifyAllClients(r)
	})

	webSocket.OnAppend(func(client *webSocket.SocketClient, r *webSocket.Room) {
		log.Println("OnAppend....")
		//发三条信息
		uCount, aCount, _ := redis.ZCard(r.RoomId)
		client.UserMsgIndex = uCount
		if aCount > Constants.FIRST_CONTENT_SIZE { //确定该client对应的作者信息的起始和结束
			client.AuthorStartIndex = aCount - Constants.FIRST_CONTENT_SIZE
			client.AuthorEndIndex = aCount - Constants.FIRST_CONTENT_SIZE
		}
		log.Println(client)

		msg, err := GetWebSocketChatMsg(redis.AuthorMessage, r.RoomId, client.AuthorEndIndex, -1)
		if err != nil {
			log.Println(err)
		}
		r.SendSelf(client, &webSocket.ChatMsg{Method: "authorMessage", Params: msg})
		client.AuthorEndIndex += len(msg)
		log.Println(client)
	})
	//该用户无连接后调用
	webSocket.OnRemove(func(userId int, r *webSocket.Room) {
		if userId <= 0 {
			return
		}
		var userType = Constants.User
		if userId == r.AuthorId {
			userType = Constants.Writer
		}
		uMsg := &UserMsg{0, "", time.Now(), userType, models.Reply, Constants.IsOut, GetUserInfo(userId)}
		log.Println(uMsg)
		//save to redis
		b, err := json.Marshal(uMsg)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = redis.ZAddUserMsg(r.RoomId, string(b))
		delete(UserInfoMap, userId)
		NotifyAllClients(r)
	})

	webSocket.OnEmit("sendMessage", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {

			if client.UserId <= 0 {//未登录
				return helper.Error(helper.NoLoginError)
			}
			if _, ok := r.ShutUpUserIds[client.UserId]; ok {//被封禁
				return helper.Error(helper.NoRightError)
			}

			uMsg := &UserMsg{}
			if err := JSON.ParseToStruct(msg.Params, uMsg); err == nil {

				if client.UserId == r.AuthorId {
					//insert into db
					tMsg, err := models.AddMessage(client.UserId, r.RoomId, uMsg.ContentType, uMsg.Content)
					if err != nil {
						log.Println(err)
						return helper.Error(helper.ParamsError)
					}
					uMsg.Id = tMsg.Id
					uMsg.CreateTime = tMsg.CreateTime
					uMsg.Info = GetUserInfo(tMsg.UserId)
					log.Println(uMsg)

					//save to redis
					b, err := json.Marshal(uMsg)
					if err != nil {
						log.Println(err)
						return helper.Error(helper.ParamsError)
					}
					_, err = redis.ZAddAuthorMsg(r.RoomId, tMsg.Id, string(b))
					if err != nil {
						log.Println(err)
						return helper.Error(helper.ParamsError)
					}

				} else {
					uMsg.Id = 0
					uMsg.CreateTime = time.Now()
					uMsg.Info = GetUserInfo(client.UserId)
					log.Println(uMsg)

					//save to redis
					b, err := json.Marshal(uMsg)
					if err != nil {
						return helper.Error(helper.ParamsError)
					}
					_, err = redis.ZAddUserMsg(r.RoomId, string(b))
					if err != nil {
						return helper.Error(helper.ParamsError)
					}
				}

				//tell thread to tell everyclient
				NotifyAllClients(r)
				return helper.Success()
			}
			return helper.Error(helper.ParamsError)
	})

	//用户点击获取更多
	webSocket.OnEmit("getMessage", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {
		param := &Param{}
		if err := JSON.ParseToStruct(msg.Params, param); err == nil {
			log.Println(param)
			if !param.Pre { //获取往后的数据
				//needs to send
				msg, err := GetWebSocketChatMsg(redis.UserMessage, r.RoomId, client.UserMsgIndex, -1)
				if err != nil {
					return helper.Error(helper.ParamsError)
				}
				r.SendSelf(client, &webSocket.ChatMsg{Method: "userMessage", Params: msg})
				client.UserMsgIndex += len(msg)

				msg, err = GetWebSocketChatMsg(redis.AuthorMessage, r.RoomId, client.AuthorEndIndex, -1)
				if err != nil {
					return helper.Error(helper.ParamsError)
				}
				r.SendSelf(client, &webSocket.ChatMsg{Method: "authorMessage", Params: msg})
				client.AuthorEndIndex += len(msg)
			} else { //获取之前的作者数据
				begin := client.AuthorStartIndex - param.Size
				if begin < 0 {
					begin = 0
				}
				msg, err := GetWebSocketChatMsg(redis.AuthorMessage, r.RoomId, client.AuthorStartIndex, client.AuthorStartIndex)
				if err != nil {
					return helper.Error(helper.ParamsError)
				}
				r.SendSelf(client, &webSocket.ChatMsg{Method: "authorMessage", Params: msg, Pre: true})
				client.AuthorEndIndex += len(msg)
			}
			r.SendSelf(client, &webSocket.ChatMsg{Method: "userCount", Params: r.GetUserCount()})
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

	redis.OnEmpty(func(roomId int64) bool {
		msg := models.ListMessage(roomId)
		size := len(msg)
		args := make(map[int64]string)
		for _, m := range msg {
			uMsg := &UserMsg{}
			uMsg.Id = m.Id
			uMsg.CreateTime = m.CreateTime
			uMsg.Content = m.Content
			uMsg.Info = GetUserInfo(m.UserId)

			b, err := json.Marshal(uMsg)
			if err != nil {
				log.Println(err)
				break
			}
			args[m.Id] = string(b)
		}
		count, err := redis.ZAddAuthorMsgs(roomId, args)
		if err != nil {
			log.Println(err)
		}
		return size == count
	})
}

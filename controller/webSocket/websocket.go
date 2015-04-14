package webSocket

import (
	"github.com/go-martini/martini"
	"chatroom/service/webSocket"
	"chatroom/helper"
	"chatroom/utils/JSON"
	"strconv"
	"chatroom/service/models"
	"chatroom/service/httpGet"
	"chatroom/service/redis"
	"github.com/golang/glog"
	"time"
	"fmt"
	"net/http"
	"encoding/json"
)

var UserInfoMap map[int]*httpGet.UserInfo = make(map[int]*httpGet.UserInfo)	//缓存用户信息

type UserMsg struct {
	Id			int64		`json:"id"`
	Content 	string		`json:"content"`
	CreateTime	time.Time	`json:"createTime"`
	Info		*httpGet.UserInfo	`json:"userInfo"`
}

type Param struct {
	Next 	bool		`json:"next"`
	Size 	int		`json:"size`
}
const (
	User httpGet.UserType = iota
	Writer
	Staff
)
const (
	FIRST_CONTENT_SIZE = 3	//进入聊天室时默认发送几条消息
)

func HandlerSocket(params martini.Params, req *http.Request, receiver <-chan *webSocket.ChatMsg, sender chan<- *webSocket.ChatMsg, done <-chan bool, disconnect chan<- int, err <-chan error) (int,string) {
	roomId := 0;
	if rId, err := strconv.Atoi(params["roomId"]); err == nil {
		roomId = rId
	}
	fmt.Println(roomId)
	info := CheckAuthorRight(req, roomId)
	if info == nil {
		return http.StatusForbidden, "no login"
	}
	return webSocket.AppendClient(info.Id, roomId, receiver, sender, done, disconnect, err)
}

func CheckAuthorRight(req *http.Request, roomId int) *httpGet.UserInfo {
	info := &httpGet.UserInfo{}
	info.CheckAuthorRight(req.Cookies(), roomId)
	fmt.Println(info)
	UserInfoMap[info.Id] = info
	return info
}

func GetUserInfo(userId int) *httpGet.UserInfo {
	info, ok := UserInfoMap[userId]
	if ok {
		return info
	}
	info = &httpGet.UserInfo{}
	info.GetUserInfo(userId)
	UserInfoMap[userId] = info
	return info
}

func GetWebSocketChatMsg(messageType redis.MessageType, roomId int, start int, stop int) ([]UserMsg, error) {
	replys, err := redis.ZRange(messageType, roomId, start, stop)
	if err != nil {
		return nil, err
	}
	msgs := []UserMsg{}
	for _,v := range replys {
		msg := UserMsg{}
		err = json.Unmarshal([]byte(v), &msg)
		fmt.Println(msg)
		if err == nil {
			msgs = append(msgs, msg)
		} else {
			fmt.Println(err)
		}
	}
	return msgs, nil
}

func init() {

	webSocket.OnAppend(func(client *webSocket.SocketClient, r *webSocket.Room) {
		//发三条信息
		ucount, acount, _ := redis.ZCard(r.RoomId)
		client.UserMsgIndex = ucount
		if acount > FIRST_CONTENT_SIZE {//确定该client对应的作者信息的起始和结束
			client.AuthorStartIndex = acount - FIRST_CONTENT_SIZE
			client.AuthorEndIndex = acount - FIRST_CONTENT_SIZE
		}

		msgs,err := GetWebSocketChatMsg(redis.AuthorMessage, r.RoomId, client.AuthorEndIndex, -1)
		if err != nil {
			return
		}
		r.SendSelf(client, &webSocket.ChatMsg{Method: "authorMessage", Params: msgs})
		client.AuthorEndIndex += len(msgs)
	})

	webSocket.OnRemove(func(userId int){
		delete(UserInfoMap, userId)
	})
	//作者发信息
	webSocket.OnEmit("authorSend", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {
		fmt.Println("authorSend")
		uMsg := &UserMsg{}
		if err := JSON.ParseToStruct(msg.Params, uMsg); err == nil {
			//insert into db
			tMsg := &models.MsgTable{}
			tMsg.UserId = client.UserId
			tMsg.RoomId = r.RoomId
			tMsg.Type = 1
			tMsg.Content = uMsg.Content
			err := tMsg.Save();
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
			if param.Next {//获取往后的数据
				//needs to send
				msgs,err := GetWebSocketChatMsg(redis.UserMessage, r.RoomId, client.UserMsgIndex, -1)
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
			} else {//获取之前的作者数据
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

	webSocket.OnEmit("publish", func(msg *webSocket.ChatMsg, client *webSocket.SocketClient, r *webSocket.Room) JSON.Type {
		content := models.PushCompleteRoom(r.RoomId)
		r.SendSelf(client, &webSocket.ChatMsg{Method: "content", Params: content})
		return helper.Success(JSON.Type{})
	})
}

package webSocket

import (
	"chatroom/helper"
	"chatroom/service/httpGet"
	"chatroom/service/models"
	"chatroom/service/redis"
	"chatroom/service/webSocket"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"reflect"
	"time"
)

var UserInfoMap map[int]*httpGet.UserInfo = make(map[int]*httpGet.UserInfo) //缓存用户信息

type UserMsg struct {
	Id         int64             `json:"id"`
	Content    string            `json:"content"`
	CreateTime time.Time         `json:"createTime"`
	Info       *httpGet.UserInfo `json:"userInfo"`
}

type Param struct {
	Next bool `json:"next"`
	Size int  `json:"size`
}

const (
	User httpGet.UserType = iota
	Writer
	Staff
)
const (
	FIRST_CONTENT_SIZE = 3 //进入聊天室时默认发送几条消息
)

func PreCheck(params martini.Params, req *http.Request, context martini.Context, w http.ResponseWriter) {
	fmt.Println("PreCheck....")
	roomId := helper.Int64(params["roomId"])
	fmt.Println(roomId)
	r, err := models.GetRoom(roomId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		fmt.Fprintf(w, "")
	}
	if r.Status == models.Closed {
		w.WriteHeader(403)
		fmt.Fprintf(w, "")
	}
	info, _ := httpGet.GetLoginUserInfo(req.Cookies(), roomId)
	if r.Price > 0 { //免费的未登陆可以进入
		if info.Code != httpGet.SUCCESS || info.Data.Id <= 0 {
			w.WriteHeader(403)
			fmt.Fprintf(w, "no login")
		}
		if info.Data.Id != r.UserId && !info.Data.Subscribed {
			w.WriteHeader(403)
			fmt.Fprintf(w, "no pay")
		}
	}
	context.Set(reflect.TypeOf(info.Data.Id), reflect.ValueOf(info.Data.Id))
	context.Set(reflect.TypeOf(roomId), reflect.ValueOf(roomId))
	//	context.Next()
}

func HandlerSocket(context martini.Context, receiver <-chan *webSocket.ChatMsg, sender chan<- *webSocket.ChatMsg, done <-chan bool, disconnect chan<- int, errc <-chan error) (int, string) {
	var rId int64 = 0
	var uId int = 0
	roomId := context.Get(reflect.TypeOf(rId)).Int()
	userId := context.Get(reflect.TypeOf(uId)).Int()
	uId = int(userId)
	fmt.Println("to handler socket...", roomId, userId)
	return webSocket.AppendClient(uId, roomId, receiver, sender, done, disconnect, errc)
}

func GetUserInfo(userId int) *httpGet.UserInfo {
	info, ok := UserInfoMap[userId]
	if ok {
		return info
	}

	result, err := httpGet.GetUserInfo(userId)
	if err != nil || result.Code != httpGet.SUCCESS {
		return nil
	}
	UserInfoMap[userId] = result.Data
	return result.Data
}

func GetWebSocketChatMsg(messageType redis.MessageType, roomId int64, start int, stop int) ([]UserMsg, error) {
	replys, err := redis.ZRange(messageType, roomId, start, stop)
	if err != nil {
		return nil, err
	}
	msgs := []UserMsg{}
	for _, v := range replys {
		msg := UserMsg{}
		err = json.Unmarshal([]byte(v), &msg)
		if err == nil {
			msgs = append(msgs, msg)
		} else {
			fmt.Println(err)
		}
	}
	return msgs, nil
}

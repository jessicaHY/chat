package webSocket

import (
	"chatroom/helper"
	"chatroom/service/httpGet"
	"chatroom/service/models"
	"chatroom/service/redis"
	"chatroom/service/webSocket"
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"reflect"
	"time"
	"log"
)

var UserInfoMap map[int]*httpGet.UserInfo = make(map[int]*httpGet.UserInfo) //缓存用户信息

type UserMsg struct {
	Id         int64             `json:"id"`
	Content    string            `json:"content"`
	CreateTime time.Time         `json:"createTime"`
	Type	   int				 `json:"type"`
	Info       *httpGet.UserInfo `json:"userInfo"`
}

type Param struct {
	Pre bool `json:"pre"`
	Size int  `json:"size"`
	UserId int	`json:"userId"`
}

const (
	User httpGet.UserType = iota
	Writer
	Staff
)
const (
	ContentType = iota
	InType
	OutType
)
const (
	FIRST_CONTENT_SIZE = 3 //进入聊天室时默认发送几条消息
)

func PreCheck(params martini.Params, rend render.Render, req *http.Request, context martini.Context) {
	log.Println("PreCheck....")
	var suc, roomId, info = UserCheck(params, rend, req)
	if !suc {
		return
	}
	if info.Data != nil {
		context.Set(reflect.TypeOf(info.Data.Id), reflect.ValueOf(info.Data.Id))
	} else {// no login
		var uId int = 0
		context.Set(reflect.TypeOf(uId), reflect.ValueOf(uId))
	}
	context.Set(reflect.TypeOf(roomId), reflect.ValueOf(roomId))
	//	context.Next()
}

func UserCheck(params martini.Params, rend render.Render, req *http.Request) (bool, int64, *httpGet.UserResult) {
	roomId := helper.Int64(params["roomId"])
	log.Println(roomId)
	r, err := models.GetRoom(roomId)
	if err != nil || r == nil {
		log.Println(err)
		rend.JSON(200, helper.Error(helper.EmptyError))
		return false, roomId, nil
	}
	if r.Status == models.Closed {
		rend.JSON(200, helper.Error(helper.ClosedError))
		return false, roomId, nil
	}
	info, _ := httpGet.GetLoginUserInfo(req.Cookies(), roomId)
	if r.Price > 0 { //免费的未登陆可以进入
		if info.Code != httpGet.SUCCESS || info.Data.Id <= 0 {
			rend.JSON(200, helper.Error(helper.NoLoginError))
			return false, roomId, nil
		}
		if info.Data.Id != r.UserId && !info.Data.Subscribed {
			rend.JSON(200, helper.Error(helper.NeedSubscribeError))
			return false, roomId, info
		}
	}
	rend.JSON(200, helper.Success())
	return true, roomId, info
}

func HandlerSocket(context martini.Context, receiver <-chan *webSocket.ChatMsg, sender chan<- *webSocket.ChatMsg, done <-chan bool, disconnect chan<- int, errc <-chan error) (int, string) {
	var rId int64 = 0
	var uId int = 0
	roomId := context.Get(reflect.TypeOf(rId)).Int()
	userId := context.Get(reflect.TypeOf(uId)).Int()
	uId = int(userId)
	log.Println("to handler socket...", roomId, userId)
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
	reply, err := redis.ZRange(messageType, roomId, start, stop)
	if err != nil {
		return nil, err
	}
	msg := []UserMsg{}
	for _, v := range reply {
		m := UserMsg{}
		err = json.Unmarshal([]byte(v), &m)
		if err == nil {
			msg = append(msg, m)
		} else {
			log.Println(err)
		}
	}
	return msg, nil
}

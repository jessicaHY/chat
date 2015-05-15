package webSocket

import (
	"chatroom/helper"
	"chatroom/service/httpGet"
	"chatroom/service/models"
	"chatroom/service/redis"
	"chatroom/service/webSocket"
	"encoding/json"
	"github.com/go-martini/martini"
	"net/http"
	"reflect"
	"log"
	"time"
	"github.com/martini-contrib/render"
)

var UserInfoMap map[int]*httpGet.UserInfo = make(map[int]*httpGet.UserInfo) //缓存用户信息

type UserMsg struct {
	Id         int64            	`json:"id"`
	Content    string           	`json:"content"`
	CreateTime time.Time        	`json:"createTime"`
	UserType	   	int				`json:"user_type"`		//用户类型
	ContentType		models.ContentType				`json:"content_type"`	//内容类型（章节，聊天）
	MessageType		int				`json:"message_type"`	//消息类型(内容，进入，离开)
	Info       *httpGet.UserInfo 	`json:"userInfo"`
}

type Param struct {
	Pre bool `json:"pre"`
	Size int  `json:"size"`
	UserId int	`json:"userId"`
}

func PreCheck(params martini.Params, req *http.Request, context martini.Context) {
	log.Println("PreCheck....")
	var suc, roomId, _, info = UserCheck(params, req)
	if !suc {
		return
	}

	log.Println("check over...")
	if info.Data != nil {
		context.Set(reflect.TypeOf(info.Data.Id), reflect.ValueOf(info.Data.Id))
	} else {// no login
		var uId int = 0
		context.Set(reflect.TypeOf(uId), reflect.ValueOf(uId))
	}
	context.Set(reflect.TypeOf(roomId), reflect.ValueOf(roomId))
	//	context.Next()
	log.Println("check over...2")
}

func UserCheck(params martini.Params, req *http.Request) (bool, int64, helper.ErrorType, *httpGet.UserResult) {
	roomId := helper.Int64(params["roomId"])
	log.Println(roomId)
	r, err := models.GetRoom(roomId)
	if err != nil || r == nil {
		log.Println(err)
		return false, roomId, helper.EmptyError, nil
	}
	if r.Status == models.Closed {
		return false, roomId, helper.ClosedError, nil
	}
	info, _ := httpGet.GetLoginUserInfo(req.Cookies(), roomId)
	if r.Price > 0 { //免费的未登陆可以进入
		if info.Code != httpGet.SUCCESS || info.Data.Id <= 0 {
			return false, roomId, helper.NoLoginError, nil
		}
		if info.Data.Id != r.UserId && !info.Data.Subscribed {
			return false, roomId, helper.NeedSubscribeError, info
		}
	}
	return true, roomId, helper.NoError, info
}

func HandlerSocket(context martini.Context, receiver <-chan *webSocket.ChatMsg, sender chan<- *webSocket.ChatMsg, done <-chan bool, disconnect chan<- int, errc <-chan error) (int, string) {

	log.Println("to HandlerSocket...")
	var rId int64 = 0
	var uId int = 0
	roomId := context.Get(reflect.TypeOf(rId)).Int()
	userId := context.Get(reflect.TypeOf(uId)).Int()
	uId = int(userId)
	log.Println("to handler socket...", roomId, userId)
	return webSocket.AppendClient(uId, roomId, receiver, sender, done, disconnect, errc)
}

func GetUserList(params martini.Params, rend render.Render) {
	roomId := helper.Int64(params["roomId"])
	log.Println(roomId)
	userIds := webSocket.ListUser(roomId)
	userInfos := make([]*httpGet.UserInfo, len(userIds), len(userIds))
	for id := range userIds {
		userInfos = append(userInfos, GetUserInfo(id))
	}
	rend.JSON(200, userInfos)
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

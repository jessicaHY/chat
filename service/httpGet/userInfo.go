package httpGet

import (
	"chatroom/helper"
	"chatroom/utils"
	"chatroom/utils/JSON"
	"io/ioutil"
	"net/http"
	"strconv"
	"log"
)

type UserType int
type UserInfo struct {
	Id         int    `json:"id"` //进入直播的时候到wings获取用户信息
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	IsAuthor   bool   `json:"isAuthor"`
	Subscribed bool   `json:"subscribed"`
}

type Result struct {
	Code    int       `json:"code"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

type UserResult struct {
	Result
	Data    *UserInfo `json:"data"`
}

const (
	SUCCESS = 1
	ERROR   = 0
)

//用于增删改直播，判断用户是否是作者
func CheckAuthorRight(cookies []*http.Cookie, bookId int) (*UserResult, error) {
	info := &UserResult{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.HOST+"/ajax/room/author/check?bookId="+strconv.Itoa(bookId), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return info, err
	}
	result := string(body)
	log.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		log.Println(err)
		return info, err
	}
	log.Println(info)
	return info, nil
}

//get logined user info
func GetLoginUserInfo(cookies []*http.Cookie, roomId int64) (*UserResult, error) {
	info := &UserResult{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.HOST+"/ajax/room/login/info?roomId="+helper.Itoa64(roomId), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return info, err
	}
	result := string(body)
	log.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		log.Println(err)
		return info, err
	}
	log.Println(info)
	return info, nil
}

func GetUserInfo(userId int) (*UserResult, error) {
	info := &UserResult{}
	resp, err := http.Get(utils.HOST + "/ajax/room/user/info?userId=" + strconv.Itoa(userId))
	if err != nil {
		log.Println(err)
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return info, err
	}
	result := string(body)
	log.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		log.Println(err)
		return info, err
	}
	log.Println(info)
	return info, nil
}

func BuyRoom(cookies []*http.Cookie, roomId int64, money int) (*UserResult, helper.ErrorType) {
	info := &UserResult{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.HOST+"/ajax/room/buy?roomId="+helper.Itoa64(roomId)+"&money="+strconv.Itoa(money), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)
	if resp.Body == nil {
		return info, helper.NetworkError
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return info, helper.IOError
	}
	result := string(body)
	log.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		log.Println(err)
		return info, helper.DataFormatError
	}
	if info.Code == ERROR {
		aaa, _ := strconv.Atoi(info.Type)
		return info, helper.GetWingsErrorType(aaa)
	}
	log.Println(info)
	return info, helper.NoError
}

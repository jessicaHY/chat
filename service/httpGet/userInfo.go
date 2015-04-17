package httpGet

import (
	"chatroom/helper"
	"chatroom/utils"
	"chatroom/utils/JSON"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UserType int
type UserInfo struct {
	Id         int    `json:"id"` //进入直播的时候到wings获取用户信息
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	IsAuthor   bool   `json:"isAuthor"`
	Subscribed bool   `json:"subscribed"`
}

type WingsResult struct {
	Code    int       `json:"code"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Data    *UserInfo `json:"data"`
}

const (
	SUCCESS = 1
	ERROR   = 0
)

//用于增删改直播，判断用户是否是作者
func CheckAuthorRight(cookies []*http.Cookie, bookId int) (*WingsResult, error) {
	info := &WingsResult{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.HOST+"/ajax/room/author/check?bookId="+strconv.Itoa(bookId), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	result := string(body)
	fmt.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		fmt.Println(err)
		return info, err
	}
	fmt.Println(info)
	return info, nil
}

//get logined user info
func GetLoginUserInfo(cookies []*http.Cookie, roomId int64) (*WingsResult, error) {
	info := &WingsResult{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.HOST+"/ajax/room/login/info?roomId="+helper.Itoa64(roomId), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)
	if resp.Body == nil {
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	result := string(body)
	fmt.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		fmt.Println(err)
		return info, err
	}
	fmt.Println(info)
	return info, nil
}

func GetUserInfo(userId int) (*WingsResult, error) {
	info := &WingsResult{}
	resp, err := http.Get(utils.HOST + "/ajax/room/user/info?userId=" + strconv.Itoa(userId))
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	result := string(body)
	fmt.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		fmt.Println(err)
		return info, err
	}
	fmt.Println(info)
	return info, nil
}

func BuyRoom(cookies []*http.Cookie, roomId int64, money int) (*WingsResult, error) {
	info := &WingsResult{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", utils.HOST+"/ajax/room/buy?roomId="+helper.Itoa64(roomId)+"&money="+strconv.Itoa(money), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)
	if resp.Body == nil {
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	result := string(body)
	fmt.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		fmt.Println(err)
		return info, err
	}
	if info.Code == ERROR {
		return info, errors.New(info.Message)
	}
	fmt.Println(info)
	return info, nil
}

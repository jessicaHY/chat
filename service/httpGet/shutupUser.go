package httpGet

import (
	"net/http"
	"chatroom/utils/Constants"
	"strconv"
	"log"
	"io/ioutil"
	"chatroom/utils/JSON"
	"chatroom/service/models"
	"strings"
)

type ShutResult struct {
	Result
	Data    []int `json:"data"`
}

func GetShutUpList(bookId int) (map[int]int, error) {
	resp, err := http.Get(Constants.HOST + "/ajax/room/shut/user/list?bookId=" + strconv.Itoa(bookId))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := string(body)
	log.Println(result)
	info := &ShutResult{}
	if err = JSON.ParseToStruct(result, info); err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("------", info)
	m := make(map[int]int)
	for _, v := range info.Data {
		m[v] = v
	}
	return m, nil
}

func AddShutUp(cookies []*http.Cookie, roomId int64, userId int, days int) (*ShutResult, error) {
	info := &ShutResult{}
	r, err := models.GetRoom(roomId)
	if err != nil || r == nil {
		return info, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", Constants.HOST+"/ajax/room/shutup/add",
		strings.NewReader("bookId=" + strconv.Itoa(r.GetHostId()) + "&userId=" + strconv.Itoa(userId) + "&days=" + strconv.Itoa(days)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)
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

func DelShutUp(cookies []*http.Cookie, roomId int64, userId int) (*ShutResult, error) {
	info := &ShutResult{}
	r, err := models.GetRoom(roomId)
	if err != nil || r == nil {
		return info, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", Constants.HOST+"/ajax/room/shutup/del",
		strings.NewReader("bookId=" + strconv.Itoa(r.GetHostId()) + "&userId=" + strconv.Itoa(userId)))
	for _, v := range cookies {
		req.AddCookie(v)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
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

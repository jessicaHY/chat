package httpGet

import (
	"net/http"
	"chatroom/helper"
	"chatroom/utils/Constants"
	"strconv"
	"strings"
	"io/ioutil"
	"log"
	"chatroom/utils/JSON"
)

func Donate(cookies []*http.Cookie, userId int, donateId int64, bookId int, price int) (*UserResult, helper.ErrorType) {
	info := &UserResult{}
	client := &http.Client{}
	param := "userId=" + strconv.Itoa(userId) + "&donateId=" + helper.Itoa64(donateId) + "&money="+strconv.Itoa(price)+"&bookId=" + strconv.Itoa(bookId)
	req, err := http.NewRequest("POST", Constants.HOST+"/system/room/donate", strings.NewReader(param))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
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

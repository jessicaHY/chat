package ajax

import (
	"chatroom/helper"
	"chatroom/service/httpGet"
	"chatroom/service/models"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"time"
)

type AjaxResult struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	ParamError   = "param error"
	NoLogin      = "no login"
	NotExisted   = "not existed"
	NetworkError = "network error"
	NoNeedToBuy  = "no need to buy"
	Failed       = "failed"
)

func Result(ar AjaxResult) string {
	b, err := json.Marshal(ar)
	if err != nil {
		return ""
	}
	return string(b)
}

func Success(message string, data interface{}) string {
	return Result(AjaxResult{"success", message, data})
}

func Error(message string) string {
	return Result(AjaxResult{"error", message, nil})
}

func AddRoom(req *http.Request) string {
	bookId, err := strconv.Atoi(req.FormValue("bookId"))
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	price, err := strconv.Atoi(req.FormValue("price"))
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.FormValue("startTime"), loc)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	fmt.Println(startTime)
	info, _ := httpGet.CheckAuthorRight(req.Cookies(), bookId)
	fmt.Println(info)
	if info.Code != httpGet.SUCCESS || !info.Data.IsAuthor {
		return Error(NoLogin)
	}
	r, err := models.AddRoom(bookId, models.BOOK, info.Data.Id, price, req.FormValue("content"), startTime)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	fmt.Println(r)
	return Success("", r)
}

func EditRoom(params martini.Params, req *http.Request) string {
	roomId := helper.Int64(params["roomId"])
	fmt.Println(roomId)
	if roomId <= 0 {
		return Error(ParamError)
	}
	//	price64, err := strconv.ParseFloat(req.FormValue("price"), 32)
	//	fmt.Println(price64)
	//	if err != nil {
	//		fmt.Println(err)
	//		return Error(ParamError)
	//	}
	//	price := float32(price64)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.FormValue("startTime"), loc)
	fmt.Println(startTime)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	//query room by roomId
	r, err := models.GetRoom(roomId)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	info, _ := httpGet.GetLoginUserInfo(req.Cookies(), 0)
	if info.Code != httpGet.SUCCESS || info.Data.Id != r.UserId {
		return Error(NoLogin)
	}
	err = models.EditRoom(roomId, req.FormValue("content"), startTime)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	return Success("", nil)
}

func CloseRoom(params martini.Params, req *http.Request) string {
	roomId := helper.Int64(params["roomId"])
	fmt.Println(roomId)
	if roomId <= 0 {
		return Error(ParamError)
	}
	//query room by roomId
	r, err := models.GetRoom(roomId)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	info, _ := httpGet.GetLoginUserInfo(req.Cookies(), 0)
	if info.Code != httpGet.SUCCESS || info.Data.Id != r.UserId {
		return Error(NoLogin)
	}
	err = models.UpdateRoomStatus(roomId, models.Closed)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	return Success("", nil)
}

func QueryRoom(params martini.Params) string {
	bookId, _ := strconv.Atoi(params["bookId"])
	fmt.Println(bookId)
	if bookId <= 0 {
		return Error(ParamError)
	}
	//query room by roomId
	r, err := models.ListNormalRoom(bookId, models.BOOK)
	if err != nil {
		fmt.Println(err)
		return Error(ParamError)
	}
	return Success("", r)
}

func BuyRoom(params martini.Params, req *http.Request) string {
	roomId := helper.Int64(params["roomId"])
	fmt.Println(roomId)
	if roomId <= 0 {
		return Error(ParamError)
	}
	r, err := models.GetRoom(roomId)
	if err != nil {
		fmt.Println(err)
		return Error(NotExisted)
	}
	//get logined user info
	info, err := httpGet.GetLoginUserInfo(req.Cookies(), roomId)
	if info.Code != httpGet.SUCCESS {
		fmt.Println(err)
		return Error(NetworkError)
	}
	//author
	if info.Data.Id == r.UserId {
		fmt.Println(err)
		return Error(NoNeedToBuy)
	}
	info, err = httpGet.BuyRoom(req.Cookies(), roomId, r.Price)
	if info.Code == httpGet.SUCCESS && info.Data.Subscribed {
		return Success("", r)
	} else {
		return Error(Failed)
	}
}

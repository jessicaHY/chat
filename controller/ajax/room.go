package ajax

import (
	"chatroom/helper"
	"chatroom/service/httpGet"
	"chatroom/service/models"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"time"
	"github.com/martini-contrib/render"
	"chatroom/utils/JSON"
	"html/template"
	"log"
)

func GetRoomByBookId(params martini.Params, rend render.Render) {
	bookId, err := strconv.Atoi(params["bookId"])
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	r, err1 := models.LastNormalRoom(bookId, models.BOOK)
	if err1 != helper.NoError {
		rend.JSON(404, helper.Error(err1))
		return
	}
	rend.JSON(200, helper.Success(r))
}

func RoomInfo(params martini.Params, rend render.Render) {
	roomId := helper.Int64(params["roomId"])
	log.Println(roomId)
	if roomId <= 0 {
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	r, err := models.GetRoom(roomId)
	if err != nil {
		log.Println(err)
		rend.JSON(404, helper.Error(helper.ParamsError))
		return
	}
	rend.JSON(200, helper.Success(r))
}

func AddRoom(req *http.Request, rend render.Render) {
	bookId, err := strconv.Atoi(req.FormValue("bookId"))
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	price, err := strconv.Atoi(req.FormValue("price"))
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.FormValue("startTime"), loc)
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	log.Println(startTime)
	info, _ := httpGet.CheckAuthorRight(req.Cookies(), bookId)
	log.Println(info)
	if info.Code != httpGet.SUCCESS || !info.Data.IsAuthor {
		rend.JSON(403, helper.Error(helper.NoLoginError))
		return
	}
	content := template.HTMLEscapeString(req.FormValue("content"))
	r, err := models.AddRoom(bookId, models.BOOK, info.Data.Id, price, content, startTime)
	if err != nil {
		log.Println(err)
		rend.JSON(500, helper.Error(helper.DbError))
		return
	}
	log.Println(r)
	rend.JSON(200, helper.Success(r))
}

func EditRoom(params martini.Params, req *http.Request, rend render.Render) {
	roomId := helper.Int64(params["roomId"])
	log.Println(roomId)
	if roomId <= 0 {
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.FormValue("startTime"), loc)
	log.Println(startTime)
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	//query room by roomId
	r, err := models.GetRoom(roomId)
	if err != nil || r == nil {
		log.Println(err)
		rend.JSON(404, helper.Error(helper.ParamsError))
		return
	}
	info, _ := httpGet.GetLoginUserInfo(req.Cookies(), 0)
	if info.Code != httpGet.SUCCESS || info.Data.Id != r.UserId {
		rend.JSON(404, helper.Error(helper.NoLoginError))
		return
	}
	content := template.HTMLEscapeString(req.FormValue("content"))
	err = models.EditRoom(roomId, content, startTime)
	if err != nil {
		log.Println(err)
		rend.JSON(404, helper.Error(helper.DbError))
		return
	}
	rend.JSON(200, helper.Success(JSON.Type{}))
}

func CloseRoom(params martini.Params, req *http.Request, rend render.Render) {
	roomId := helper.Int64(params["roomId"])
	log.Println(roomId)
	if roomId <= 0 {
		rend.JSON(404, helper.Error(helper.ParamsError))
		return
	}
	//query room by roomId
	r, err := models.GetRoom(roomId)
	if err != nil {
		log.Println(err)
		rend.JSON(404, helper.Error(helper.ParamsError))
		return
	}
	info, _ := httpGet.GetLoginUserInfo(req.Cookies(), 0)
	if info.Code != httpGet.SUCCESS || info.Data.Id != r.UserId {
		rend.JSON(404, helper.Error(helper.NoLoginError))
		return
	}
	err = models.UpdateRoomStatus(roomId, models.Closed)
	if err != nil {
		log.Println(err)
		rend.JSON(404, helper.Error(helper.ParamsError))
		return
	}
	rend.JSON(200, helper.Success(JSON.Type{}))
}

func QueryRoom(params martini.Params, rend render.Render) {
	bookId, _ := strconv.Atoi(params["bookId"])
	log.Println(bookId)
	if bookId <= 0 {
		rend.JSON(404, helper.Error(helper.ParamsError))
		return
	}
	//query room by roomId
	r, err := models.ListNormalRoom(bookId, models.BOOK)
	if err != nil {
		log.Println(err)
		rend.JSON(404, helper.Error(helper.ParamsError))
		return
	}
	rend.JSON(200, helper.Success(r))
}

func BuyRoom(params martini.Params, req *http.Request, rend render.Render) {
	roomId := helper.Int64(params["roomId"])
	log.Println(roomId)
	if roomId <= 0 {
		rend.JSON(200, helper.Error(helper.ParamsError))
		return
	}
	r, err := models.GetRoom(roomId)
	if err != nil {
		log.Println(err)
		rend.JSON(200, helper.Error(helper.EmptyError))
		return
	}
	//get logined user info
	info, err := httpGet.GetLoginUserInfo(req.Cookies(), roomId)
	if info.Code != httpGet.SUCCESS {
		log.Println(err)
		rend.JSON(200, helper.Error(helper.NetworkError))
		return
	}
	//author
	if info.Data.Id == r.UserId {
		log.Println(err)
		rend.JSON(200, helper.Error(helper.NoNeedError))
		return
	}
	info, errType := httpGet.BuyRoom(req.Cookies(), roomId, r.Price, r.GetHostId())
	if info.Code == httpGet.SUCCESS && info.Data.Subscribed {
		rend.JSON(200, helper.Success())
	} else {
		rend.JSON(200, helper.Error(errType))
	}
}

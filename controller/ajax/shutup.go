package ajax

import (
	"net/http"
	"log"
	"strconv"
	"github.com/martini-contrib/render"
	"chatroom/helper"
	"chatroom/service/httpGet"
	"chatroom/utils/JSON"
	"chatroom/service/webSocket"
)

func AddShutup(req *http.Request, rend render.Render) {
	roomId := helper.Int64(req.FormValue("roomId"))
	if roomId <= 0 {
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	userId, err := strconv.Atoi(req.FormValue("userId"))
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	days, err := strconv.Atoi(req.FormValue("days"))
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	log.Println("add shutup...", roomId, userId, days)
	result, err := httpGet.AddShutUp(req.Cookies(), roomId, userId, days )
	log.Println(result)
	if result.Code != httpGet.SUCCESS {
		rend.JSON(500, helper.Error(helper.DefaultError))
		return
	}
	webSocket.AddShutUp(roomId, userId)
	rend.JSON(200, helper.Success(JSON.Type{}))
}

func DelShutup(req *http.Request, rend render.Render) {
	roomId := helper.Int64(req.FormValue("roomId"))
	if roomId <= 0 {
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}
	userId, err := strconv.Atoi(req.FormValue("userId"))
	if err != nil {
		log.Println(err)
		rend.JSON(403, helper.Error(helper.ParamsError))
		return
	}

	result, err := httpGet.DelShutUp(req.Cookies(), roomId, userId )
	log.Println(result)
	if result.Code != httpGet.SUCCESS {
		rend.JSON(500, helper.Error(helper.DefaultError))
		return
	}
	webSocket.DelShutUp(roomId, userId)
	rend.JSON(200, helper.Success(JSON.Type{}))
}

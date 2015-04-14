package models

import (
	"time"
	"github.com/golang/glog"
	"strconv"
	"fmt"
	"chatroom/helper"
)


type ContentType int
const (
	Chapter ContentType = iota
	Reply
)

type MsgTable struct {
	Id         int64			`xorm:"id pk autoincr"`
	UserId		int				`xorm:"user_id"`
	RoomId		int 			`xorm:"room_id"`
	Type		ContentType		`xorm:"type"`
	Content 	string			`xorm:"content"`
	CreateTime	time.Time		`xorm:"created"`
}

type RoomTable struct {
	Id		int64		`xorm:"id pk"`
	Status	int8		`xorm:"status"`
}

func (msg *MsgTable) Save() error {
	_, err := engine.Insert(msg)
	if err != nil {
		glog.Fatalln(err)
		return err
	}
	return nil
}

func ListMessage(roomId int) []*MsgTable {
	msgs := []*MsgTable{}
	engine.Where("room_id=?", roomId).Find(&msgs)
	return  msgs
}

func (r *RoomTable) Save() error {
	_, err := engine.Insert(r)
	if err != nil {
		glog.Fatalln(err)
		return err
	}
	return nil
}

func CheckRoom(roomId int) {
	r, err := GetRoom(roomId)
	if err != nil {
		glog.Fatalln(err)
		return
	}
	if r == nil {
		r = &RoomTable{}
		r.Id = helper.Int64(roomId)
		r.Status = 0
		engine.Insert(r)
	}
}

func GetRoom(roomId int) (*RoomTable, error) {
	rId, err := strconv.ParseInt(strconv.Itoa(roomId), 10, 64);
	fmt.Println(rId)
	fmt.Println(roomId)
	if err != nil {
		glog.Fatalln(err)
		return nil, err
	}
	var r RoomTable
	has, err := engine.Id(rId).Get(&r)
	if err != nil {
		glog.Fatalln(err)
		return nil, err
	}
	if has {
		return &r, err
	}
	return nil, err
}

func ListRoomTable() []*RoomTable {
	rooms := []*RoomTable{}
	engine.Where("status=0").Find(&rooms)
	return rooms
}

func PushCompleteRoom(roomId int) string {
	msgs := []*MsgTable{}
	engine.Where("room_id=?", roomId).Find(&msgs)
	content := ""
	for _,m := range msgs {
		content = content + `
		` + m.Content
	}
	fmt.Println(content)
	return content
}

package models

import (
	"fmt"
	"github.com/golang/glog"
	"time"
)

type ContentType int

const (
	Chapter ContentType = iota
	Reply
)

type MsgTable struct {
	Id         int64       `xorm:"id pk autoincr"`
	UserId     int         `xorm:"user_id"`
	RoomId     int64       `xorm:"room_id"`
	Type       ContentType `xorm:"type"`
	Content    string      `xorm:"content"`
	CreateTime time.Time   `xorm:"created"`
}

func AddMessage(userId int, roomId int64, contentType ContentType, content string) (*MsgTable, error) {
	msg := &MsgTable{UserId: userId, RoomId: roomId, Type: contentType, Content: content}
	_, err := engine.Insert(msg)
	if err != nil {
		glog.Fatalln(err)
	}
	return msg, err
}

func ListMessage(roomId int64) []MsgTable {
	msgs := []MsgTable{}
	err := engine.Where("room_id=?", roomId).Find(&msgs)
	if err != nil {
		glog.Fatalln(err)
	}
	return msgs
}

func PushCompleteRoom(roomId int64) string {
	msgs := ListMessage(roomId)
	content := ""
	for _, m := range msgs {
		content = content + `
		` + m.Content
	}
	fmt.Println(content)
	return content
}

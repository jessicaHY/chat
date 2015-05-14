package models

import (
	"chatroom/helper"
	"errors"
	"time"
	"log"
)

type RoomTable struct {
	Id        int64      `xorm:"id pk autoincr" json:"id"`
	HostIt    int64      `xorm:"host_it index" json:"hostIt"`
	UserId    int        `xorm:"user_id"  json:"userId"`
	Price     int        `xorm:"price"  json:"price"`
	Content   string     `xorm:"content"  json:"content"`
	StartTime time.Time  `xorm:"start_time"  json:"startTime"`
	Status    RoomStatus `xorm:"status"  json:"status"`
}

type RoomStatus int8

const (
	Normal RoomStatus = 0
	Ended  RoomStatus = 3
	Closed RoomStatus = -1
)

type HostType int8

const (
	BOOK HostType = 2
)

func AddRoom(hostId int, hostType HostType, userId int, price int, content string, startTime time.Time) (*RoomTable, error) {
	hostIt := helper.Int64(hostId)*10000 + helper.Int64(int8(hostType))
	log.Println(hostIt)
	count := CountNormalRoom(hostId, hostType)
	if count > 0 {
		return nil, errors.New("already exists")
	}
	r := &RoomTable{HostIt: hostIt, UserId: userId, Price: price, Content: content, StartTime: startTime, Status: Normal}
	_, err := engine.InsertOne(r)
	return r, err
}

func EditRoom(roomId int64, content string, startTime time.Time) error {
	r := new(RoomTable)
	r.Content = content
	r.StartTime = startTime
	_, err := engine.Id(roomId).Update(r)
	return err
}

func UpdateRoomStatus(roomId int64, status RoomStatus) error {
	r := new(RoomTable)
	r.Status = status
	_, err := engine.Id(roomId).Cols("status").Update(r)
	return err
}

func ListAllRoom(hostId int, hostType HostType) ([]RoomTable, error) {
	hostIt := helper.Int64(hostId)*10000 + helper.Int64(int8(hostType))
	rs := []RoomTable{}
	err := engine.Where("host_it=?", hostIt).Find(&rs)
	return rs, err
}

func CountNormalRoom(hostId int, hostType HostType) int64 {
	hostIt := helper.Int64(hostId)*10000 + helper.Int64(int8(hostType))
	log.Println(hostIt)
	room := new(RoomTable)
	total, err := engine.Where("host_it=? and status=?", hostIt, Normal).Count(room)
	if err != nil {
		log.Println(err)
	}
	return total
}

func ListNormalRoom(hostId int, hostType HostType) ([]RoomTable, error) {
	hostIt := helper.Int64(hostId)*10000 + helper.Int64(int8(hostType))
	rs := []RoomTable{}
	err := engine.Where("host_it=? and status=?", hostIt, Normal).Find(&rs)
	return rs, err
}

func LastNormalRoom(hostId int, hostType HostType) (RoomTable, helper.ErrorType) {
	hostIt := helper.Int64(hostId)*10000 + helper.Int64(int8(hostType))
	rs := []RoomTable{}
	err := engine.Where("host_it=? and status=?", hostIt, Normal).Desc("id").Limit(1, 0).Find(&rs)
	if err != nil {
		return RoomTable{}, helper.DbError
	}
	if len(rs) == 0 {
		return RoomTable{}, helper.EmptyError
	} else {
		return rs[0], helper.NoError
	}
}

func (r *RoomTable) GetHostId() int {
	return int(r.HostIt / 10000)
}

func GetRoom(roomId int64) (*RoomTable, error) {
	r := &RoomTable{}
	has, err := engine.Id(roomId).Get(r)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	if has {
		return r, err
	}
	return nil, err
}

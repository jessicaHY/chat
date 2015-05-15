package models

import (
	"time"
	"chatroom/helper"
	"log"
	"chatroom/utils/Constants"
)

type DonateTable struct {
	Id		int64  	`xorm:"id pk autoincr"`
	UserId 	int		`xorm:"user_id"`
	RoomId 	int64	`xorm:"room_id"`
	GiftId	int64		`xorm:"gift_id"`
	Count	int		`xorm:"count"`
	Price	int		`xorm:"price"`
	Status	int		`xorm:"status"`
	PayStatus	bool	`xorm:"pay_status"`
	CreateTime	time.Time 	`xorm:"create_time"`
	UpdateTime 	time.Time	`xorm:"update_time"`
}


func ListDonateByRoom(roomId int64) ([]DonateTable, error) {
	rs := []DonateTable{}
	err := engine.Where("room_id=? and status=? and pay_status=?", roomId, Constants.STATUS_NORMAL, true).Find(&rs)
	return rs, err
}

func AddDonate(roomId int64, giftId int64, userId int, count int) (*DonateTable, helper.ErrorType) {
	g, err := GetGift(giftId)
	if err != nil {
		log.Fatalln(err)
		return nil, helper.DbError
	}
	if g == nil {
		return nil, helper.EmptyError
	}

	d, err := addDonate(roomId, giftId, userId, count, g.Price*count)
	if err != nil {
		log.Fatalln(err)
		return nil, helper.DbError
	}
	return d, helper.NoError
}

func addDonate(roomId int64, giftId int64, userId int, count int, price int) (*DonateTable, error) {
	d := &DonateTable{}
	d.UserId = userId
	d.RoomId = roomId
	d.GiftId = giftId
	d.Count = count
	d.Price = price
	d.Status = Constants.STATUS_NORMAL
	d.PayStatus = false
	d.CreateTime = time.Now()

	_, err := engine.InsertOne(d)
	return d, err
}

func UpdateDonate(donateId int64) error {
	d := &DonateTable{}
	d.PayStatus = true
	d.UpdateTime = time.Now()
	_, err := engine.Id(donateId).Update(d)
	log.Println("UpdateDonate", err)
	return err
}

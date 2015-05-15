package models

import (
	"time"
	"chatroom/utils/Constants"
	"log"
)

type GiftTable struct {
	Id 	int64	`xorm:"id pk autoincr" json:"id"`
	Name	string 	`xorm:"name" json:"name"`
	Price	int		`xorm:"price" json:"price"`
	Icon	string	`xorm:"icon" json:"icon"`
	Unit	string	`xorm:"unit" json:"unit"`
	Content string 	`xorm:"content" json:"content"`
	Group	Constants.GroupType	`xorm:"_group" json:"group"`
	Status	int8	`xorm:"status" json:"status"`
	CreateTime	time.Time	`xorm:"create_time" json:"createTime"`
}

func ListGift(group Constants.GroupType) ([]GiftTable, error){
	rs := []GiftTable{}
	err := engine.Where("_group=? and status=?", group, Constants.STATUS_NORMAL).Find(&rs)
	return rs, err
}

func AddGift(name string, price int, icon string, unit string, content string, group Constants.GroupType) (*GiftTable, error) {
	g := &GiftTable{}
	g.Name = name
	g.Price = price
	g.Icon = icon
	g.Unit = unit
	g.Content = content
	g.Group = group
	g.Status = Constants.STATUS_NORMAL
	g.CreateTime = time.Now()

	_, err := engine.InsertOne(g)
	return g, err
}

func EditGift(id int64, name string, price int, icon string, unit string, content string, group Constants.GroupType) error {
	g := new(GiftTable)
	g.Name = name
	g.Price = price
	g.Icon = icon
	g.Unit = unit
	g.Content = content
	g.Group = group
	g.Status = Constants.STATUS_NORMAL
	g.CreateTime = time.Now()
	_, err := engine.Id(id).Update(g)
	return err
}

func GetGift(id int64) (*GiftTable, error) {
	g := &GiftTable{}
	has, err := engine.Id(id).Get(g)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	if has {
		return g, err
	}
	return nil, err
}

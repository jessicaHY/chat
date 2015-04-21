package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	"log"
)

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "rock:heiyan@(192.168.1.111:3306)/chatroom?charset=utf8")
	if err != nil {
		log.Fatalln(err)
	}
	engine.SetMapper(core.SnakeMapper{})
	GlobalInit()
}

func GlobalInit() {
	err := engine.Sync2(new(MsgTable), new(RoomTable))
	if err != nil {
		log.Fatalln(err)
	}
}

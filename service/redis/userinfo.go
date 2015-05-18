package redis

import (
	"chatroom/helper"
	"github.com/garyburd/redigo/redis"
	"log"
)

var onUserInfoEmpty = func(int64, int) bool {return false }

func OnUserInfoEmpty(callback func(int64, int) bool) {
	onUserInfoEmpty = callback
}

const (
	MapPre = "UserInfo:"

)

func HSetUserInfo(roomId int64, userId int, content string) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("HSET", MapPre + helper.Itoa64(roomId), userId, content));
}

func HGetUserInfo(roomId int64, userId int) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("HGET", MapPre + helper.Itoa64(roomId), userId))
}

func HMGetUserInfo(roomId int64, userIds []int) ([]string, error) {
	var key = MapPre + helper.Itoa64(roomId);
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	for _, v := range userIds {
		conn.Send("HGET", key, v)
	}
	replys, err := redis.Strings(conn.Do("EXEC"))
	if err != nil {
		log.Println(err)
	}
	return replys, err
}

func HDelUserInfo(roomId int64, userId int) bool {
	conn := pool.Get()
	defer conn.Close()

	res, err := redis.Int(conn.Do("HDEL", MapPre + helper.Itoa64(roomId), userId))
	if err != nil {
		log.Println(err)
	}
	return res == 1
}


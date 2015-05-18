package redis

import (
	"chatroom/helper"
	"log"
	"github.com/garyburd/redigo/redis"
)


type MessageType int

const (
	UserMessage MessageType = iota
	AuthorMessage
)
const (
	AuthorPre = "author:"
	UserPre   = "user:"
)

var onMessageEmpty = func(int64) bool { return false }

func OnMessageEmpty(callback func(int64) bool) {
	onMessageEmpty = callback
}

func ZAddUserMsg(roomId int64, content string) (int, error) {
	conn := pool.Get()
	defer conn.Close()
	key := UserPre + helper.Itoa64(roomId)
	num, err := conn.Do("ZCARD", key)
	if err != nil {
		return 0, err
	}
	return redis.Int(conn.Do("ZADD", key, num, content))
}

func ZAddAuthorMsg(roomId int64, msgId interface{}, content string) (int, error) {
	checkAuthroMsgExisted(roomId)
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("ZADD", AuthorPre+helper.Itoa64(roomId), msgId, content))
}

func ZAddAuthorMsgs(roomId int64, m map[int64]string) (int, error) {
	//这里不需要check，该方法是批量初始化redis里的数据
	log.Println("ZAddAuthorMsgs....")
	key := AuthorPre + helper.Itoa64(roomId)
	conn := pool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	for i, v := range m {
		conn.Send("ZADD", key, i, v)
	}
	replys, err := redis.Values(conn.Do("EXEC"))
	return len(replys), err
}

func ZRange(mType MessageType, roomId int64, start int, stop int) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	switch mType {
	case UserMessage:
		return redis.Strings(conn.Do("ZRANGE", UserPre+helper.Itoa64(roomId), start, stop))
	case AuthorMessage:
		checkAuthroMsgExisted(roomId)
		return redis.Strings(conn.Do("ZRANGE", AuthorPre+helper.Itoa64(roomId), start, stop))
	}
	return nil, nil
}

func Del(roomId int64) {
	conn := pool.Get()
	defer conn.Close()
	conn.Do("DEL", UserPre+helper.Itoa64(roomId))
}

func ZCard(roomId int64) (int, int, error) {
	conn := pool.Get()
	defer conn.Close()

	err := conn.Send("ZCARD", UserPre+helper.Itoa64(roomId))
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	checkAuthroMsgExisted(roomId)
	err = conn.Send("ZCARD", AuthorPre+helper.Itoa64(roomId))
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	err = conn.Flush()
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	ucount, err := redis.Int(conn.Receive())
	acount, err := redis.Int(conn.Receive())
	return ucount, acount, err
}

func checkAuthroMsgExisted(roomId int64) {
	conn := pool.Get()
	exists, err := redis.Bool(conn.Do("EXISTS", AuthorPre+helper.Itoa64(roomId)))
	conn.Close()

	if err != nil {
		log.Fatalln(err)
		return
	}
	if !exists {
		onMessageEmpty(roomId)
	}
}

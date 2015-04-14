package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"strconv"
)

var (
	pool *redis.Pool
)

type MessageType int
const (
	UserMessage MessageType = iota
	AuthorMessage
)
const (
	AuthorPre = "author:"
	UserPre = "user:"
)
func init() {
	pool = &redis.Pool {
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
//			if _, err := c.Do("AUTH",""); err != nil {
//				c.Close()
//				return nil, err
//			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func ZAddUserMsg(roomId int, content string) (int, error) {
	conn := pool.Get()
	defer conn.Close()
	key := UserPre + strconv.Itoa(roomId)
	num, err := conn.Do("ZCARD", key)
	if err != nil {
		return  0, err
	}
	return redis.Int(conn.Do("ZADD", key, num, content))
}

func ZAddAuthorMsg(roomId int, msgId interface {}, content string) (int, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("ZADD", AuthorPre + strconv.Itoa(roomId), msgId, content))
}

func ZRange(mType MessageType, roomId int, start int, stop int) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	switch mType{
	case UserMessage:
		return  redis.Strings(conn.Do("ZRANGE", UserPre + strconv.Itoa(roomId), start, stop))
	case AuthorMessage:
		return  redis.Strings(conn.Do("ZRANGE", AuthorPre + strconv.Itoa(roomId), start, stop))
	}
	return nil, nil
}

func Del(roomId int) {
	conn := pool.Get()
	defer conn.Close()
	conn.Do("DEL", UserPre + strconv.Itoa(roomId))
}

func ZCard(roomId int) (int, int, error) {
	conn := pool.Get()
	defer conn.Close()

	err := conn.Send("ZCARD", UserPre + strconv.Itoa(roomId))
	if err != nil {
		return 0, 0, err
	}
	err = conn.Send("ZCARD", AuthorPre + strconv.Itoa(roomId))
	if err != nil {
		return 0, 0, err
	}
	err = conn.Flush()
	if err != nil {
		return 0, 0, err
	}
	ucount, err := redis.Int(conn.Receive())
	acount, err := redis.Int(conn.Receive())
	return ucount, acount, err
}

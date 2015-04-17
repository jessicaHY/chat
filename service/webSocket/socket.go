package webSocket

import (
	"chatroom/helper"
	"chatroom/service/models"
	"chatroom/utils/JSON"
	"fmt"
	"github.com/golang/glog"
	"sync"
)

var syncLock = sync.Mutex{} //LiveMap
var RoomMap map[int64]*Room = make(map[int64]*Room)
var onEmitCallback = map[string]func(*ChatMsg, *SocketClient, *Room) JSON.Type{}
var onAppend = func(*SocketClient, *Room) {}
var onRemove = func(int) {}

type Room struct {
	sync.Mutex
	RoomId        int64
	AuthorId      int
	ThreadChannel chan byte //每个room对应一个goroutine来执行任务
	clients       []*SocketClient
}

type SocketClient struct {
	UserId           int //进入直播的时候确定鉴权UserId
	AuthorStartIndex int
	AuthorEndIndex   int
	UserMsgIndex     int
	in               <-chan *ChatMsg
	out              chan<- *ChatMsg
	done             <-chan bool
	disconnect       chan<- int
	err              <-chan error
}

type ChatMsg struct {
	Method string      `json:"method"`
	Params interface{} `json:"data"`
}

//注册回调
func OnEmit(method string, callback func(*ChatMsg, *SocketClient, *Room) JSON.Type) {
	onEmitCallback[method] = callback
}
func OnAppend(callback func(*SocketClient, *Room)) {
	onAppend = callback
}

func OnRemove(callback func(int)) {
	onRemove = callback
}

func GetRoom(roomId int64) *Room {
	syncLock.Lock()
	defer syncLock.Unlock()

	room, ok := RoomMap[roomId]
	if ok {
		return room
	}
	rt, err := models.GetRoom(roomId)
	if err != nil {
		glog.Fatalln(err)
	}
	r := &Room{sync.Mutex{}, roomId, rt.UserId, make(chan byte), make([]*SocketClient, 0)}
	RoomMap[roomId] = r
	go r.NewThreadTask()
	return r
}

func AppendClient(userId int, roomId int64, receiver <-chan *ChatMsg, sender chan<- *ChatMsg, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
	client := &SocketClient{userId, 0, 0, 0, receiver, sender, done, disconnect, err}
	r := GetRoom(roomId)
	r.AppendClient(client)

	fmt.Println("waiting for msg...")
	for {
		select {
		case <-client.err:
		// Don't try to do this:
		// client.out <- &Message{"system", "system", "There has been an error with your connection"}
		// The socket connection is already long gone.
		// Use the error for statistics etc
		case msg := <-client.in:
			r.Emit(client, msg)
		case <-client.done:
			r.RemoveClient(client)
			return 200, "OK"
		}
	}
}

func (r *Room) AppendClient(client *SocketClient) {
	r.Lock()
	defer r.Unlock()

	r.clients = append(r.clients, client)
	onAppend(client, r)
}

func (r *Room) RemoveClient(client *SocketClient) {
	r.Lock()
	defer r.Unlock()

	for index, c := range r.clients {
		if c == client {
			r.clients = append(r.clients[:index], r.clients[(index+1):]...)
			onRemove(client.UserId)
		}
	}
}

func (r *Room) Emit(client *SocketClient, msg *ChatMsg) {
	fmt.Println("Emit....")
	syncLock.Lock()
	defer syncLock.Unlock()
	method := msg.Method

	if _, found := onEmitCallback[method]; method != "" && found {
		onEmitCallback[method](msg, client, r)
	} else {
		client.out <- &ChatMsg{method, helper.Error("method undefined")}
	}
}

//启动一个goroutine，用来监听管道执行推送任务
func (r *Room) NewThreadTask() {
	for {
		<-r.ThreadChannel
		fmt.Println("NotifyAllClients")
		for _, c := range r.clients {
			c.out <- &ChatMsg{Method: "hasMessage"}
		}
		<-r.ThreadChannel

		fmt.Println("NotifyAllClients over...")
	}
}

func (r *Room) SendSelf(client *SocketClient, msg *ChatMsg) {
	client.out <- msg
}

package webSocket

import (
	"chatroom/helper"
	"chatroom/service/models"
	"chatroom/utils/JSON"
	"sync"
	"log"
	"chatroom/service/httpGet"
)

var syncLock = sync.Mutex{} //LiveMap
var RoomMap map[int64]*Room = make(map[int64]*Room)
var onEmitCallback = map[string]func(*ChatMsg, *SocketClient, *Room) JSON.Type{}
var onAppend = func(*SocketClient, *Room) {}
var beforeAppend = func(*SocketClient, *Room) {}
var onRemove = func(int, *Room) {}

type Room struct {
	sync.Mutex
	RoomId        int64
	AuthorId      int
	ShutUpUserIds	map[int]int
	ThreadChannel chan bool //每个room对应一个goroutine来执行任务
	clientsMap      map[int][]*SocketClient	//一个用户可能从多个终端登录，有多个socket连接
}

type SocketClient struct {
	UserId          int //进入直播的时候确定鉴权UserId
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
	Pre bool `json:"pre"`	//返回客户端向前还是向后
}

//注册回调
func OnEmit(method string, callback func(*ChatMsg, *SocketClient, *Room) JSON.Type) {
	onEmitCallback[method] = callback
}

func BeforeAppend(callback func(*SocketClient, *Room)) {
	beforeAppend = callback
}

func OnAppend(callback func(*SocketClient, *Room)) {
	onAppend = callback
}

func OnRemove(callback func(int, *Room)) {
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
		log.Fatalln(err)
	}

	dbRoom, err := models.GetRoom(roomId)
	if err != nil {
		log.Fatalln(err)
	}
	m, err := httpGet.GetShutUpList(dbRoom.GetHostId())
	if err != nil {
		log.Fatalln(err)
	}
	r := &Room{sync.Mutex{}, roomId, rt.UserId, m, make(chan bool), make(map[int][]*SocketClient)}
	RoomMap[roomId] = r
	go r.NewThreadTask()
	return r
}

func AppendClient(userId int, roomId int64, receiver <-chan *ChatMsg, sender chan<- *ChatMsg, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
	client := &SocketClient{userId, 0, 0, 0, receiver, sender, done, disconnect, err}
	r := GetRoom(roomId)
	r.AppendClient(client)

	log.Println("waiting for msg...")
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

	cs, _ := r.clientsMap[client.UserId]
	if cs != nil {
		r.clientsMap[client.UserId] = append(r.clientsMap[client.UserId], client)
	} else {
		beforeAppend(client, r)//发送消息
		r.clientsMap[client.UserId] = []*SocketClient{client}
	}
	onAppend(client, r)
}

func (r *Room) RemoveClient(client *SocketClient) {
	r.Lock()
	defer r.Unlock()
	log.Println(*r)
	log.Println(r.clientsMap)

	values, _ := r.clientsMap[client.UserId]
	for index, c := range values {
		if c == client {
			r.clientsMap[client.UserId] = append(values[:index], values[(index+1):]...)

			log.Println(r.clientsMap)
			if len(r.clientsMap[client.UserId]) == 0 {
				delete(r.clientsMap, client.UserId)
				onRemove(client.UserId, r)
			}
			if len(r.clientsMap) == 0 {
				r.ThreadChannel <- false
				delete(RoomMap, r.RoomId)
			}
			break
		}
	}
}

func (r *Room) Emit(client *SocketClient, msg *ChatMsg) {
	log.Println("Emit....")
	syncLock.Lock()
	defer syncLock.Unlock()
	method := msg.Method

	if _, found := onEmitCallback[method]; method != "" && found {
		onEmitCallback[method](msg, client, r)
	} else {
		client.out <- &ChatMsg{method, helper.Error("method undefined"), false}
	}
}

//启动一个goroutine，用来监听管道执行推送任务
func (r *Room) NewThreadTask() {
	for {
		b := <-r.ThreadChannel
		if b {
			log.Println("NotifyAllClients begin...")
			for _, cs := range r.clientsMap {
				for _, c := range cs {
					c.out <- &ChatMsg{Method: "hasMessage"}
				}
			}
			log.Println("NotifyAllClients over...")
		} else {
			break
		}
	}
}

func (r *Room) SendSelf(client *SocketClient, msg *ChatMsg) {
	client.out <- msg
}

func (r *Room) GetUserCount() [2]int {
	cs, _ := r.clientsMap[0]
	if cs != nil {
		return [2]int{len(r.clientsMap) - 1, len(r.clientsMap[0])}
	} else {
		return [2]int{len(r.clientsMap), 0}
	}
}

func AddShutUp(roomId int64, userId int) {
	syncLock.Lock()
	defer syncLock.Unlock()

	room, _ := RoomMap[roomId]
	if room != nil {
		room.ShutUpUserIds[userId] = userId
	}
}

func DelShutUp(roomId int64, userId int) {
	syncLock.Lock()
	defer syncLock.Unlock()

	room, _ := RoomMap[roomId]
	if room != nil {
		delete(room.ShutUpUserIds, userId)
	}
}

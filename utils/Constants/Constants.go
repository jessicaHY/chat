package Constants

const (
	HOST = "http://localhost:8080"
)

//wings
type HttpResult struct {
	Error   error       `json:"error"`
	Code    int         `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}


//用户类型
const (
	User = iota
	Writer
	Staff
)
//内容类型
const (
	IsContent = iota
	IsIn
	IsOut
)
const (
	FIRST_CONTENT_SIZE = 3 //进入聊天室时默认发送几条消息
)

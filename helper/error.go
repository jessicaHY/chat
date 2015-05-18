package helper

import (
	"errors"
)

type ErrorType int

const (
	NoError ErrorType = iota
	DefaultError
	ExistsError
	ParamsError
	RequestError
	DbError
	EmptyError
	OfflineError
	BusyError
	NoLoginError
	NoRightError
	ClosedError //直播间已经关闭
	NeedSubscribeError	//需要购买
	NetworkError
	NoNeedError
	IOError
	DataFormatError

	WingsNoLoginError = 101
	WingsParamError = 102
	WingsNoEnoughMoneyError = 105
	WingsSuccessDbFail = 200 //wings操作成功，但是这边存储数据库失败

)

func GetWingsErrorType(code string) ErrorType {
	switch(code) {
	case "ac-62":
		return WingsNoLoginError
	case "pay-1":
		return WingsNoEnoughMoneyError
	case "ob-5":
		return WingsParamError
	default:
		return DefaultError
	}
}
func NewError(msg string, err ...error) error {
	str := msg
	if len(err) > 0 && err[0] != nil {
		str += ": " + err[0].Error()
	}
	return errors.New(str)
}

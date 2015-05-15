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
	WingsNoEnoughMoneyError = 105
	WingsSuccessDbFail = 200
)

func GetWingsErrorType(code int) ErrorType {
	switch(code) {
	case 1:
		return WingsNoLoginError
	case 5:
		return WingsNoEnoughMoneyError
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

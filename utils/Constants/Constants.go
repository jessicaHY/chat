package Constants

import (
	"net/http"
	"strings"
)

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
type GroupType int8
type SiteType int8
const (
	FIRST_CONTENT_SIZE = 3 //进入聊天室时默认发送几条消息
	STATUS_NORMAL = 0
	STATUS_DELETED = -1

	GROUP_HEIYAN GroupType = 1
	GROUP_RUOCHU GroupType = 2
	GROUP_RUOXIA GroupType = 3

	SITE_WEB SiteType = 1
	SITE_WAP SiteType = 2
	SITE_ANDROID	SiteType = 3
	SITE_IOS	SiteType = 4
)

func GetGroupFromReq(req  *http.Request) GroupType {
	var host = req.URL.Host
	if strings.Contains(host, "ruochu") {
		return GROUP_RUOCHU
	} else if strings.Contains(host, "ruoxia") {
		return GROUP_RUOXIA
	} else {
		return GROUP_HEIYAN
	}
}

func GetSiteFromReq(req *http.Request) SiteType {
	var host = req.URL.Host
	if strings.HasPrefix(host, "m.") {
		return SITE_WAP
	} else if strings.HasPrefix(host, "apk.") {
		return SITE_ANDROID
	} else if strings.HasPrefix(host, "ios.") {
		return SITE_IOS
	} else {
		return SITE_WEB
	}
}

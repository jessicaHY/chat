package httpGet
import (
	"github.com/golang/glog"
	"net/http"
	"strconv"
	"chatroom/utils/JSON"
	"io/ioutil"
	"fmt"
)

type UserType int
type UserInfo struct {
	Id		int  	`json:"id"`	//进入直播的时候到wings获取用户信息
	Name	string	`json:"name"`
	Icon	string	`json:"icon"`
	IsAuthor bool	`json:"isAuthor"`
}

func (info *UserInfo) GetUserInfo(userId int) {
	resp, err := http.Get("http://localhost/ajax/user/info?userId=" + strconv.Itoa(userId))
	if err != nil {
		glog.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatalln(err)
	}
	result := string(body)
	fmt.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		glog.Fatalln(err)
	}
}

//用于增删改直播，判断用户是否是作者
func (info *UserInfo) CheckAuthorRight(cookies []*http.Cookie, bookId int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost/ajax/user/check?bookId=" +  strconv.Itoa(bookId), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	result := string(body)
	fmt.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		fmt.Println(err)
	}
}

func (info *UserInfo) CheckUserConsume(cookies []*http.Cookie, roomId int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost/ajax/user/check?roomId=" +  strconv.Itoa(roomId), nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	result := string(body)
	fmt.Println(result)
	if err = JSON.ParseToStruct(result, info); err != nil {
		fmt.Println(err)
	}
}

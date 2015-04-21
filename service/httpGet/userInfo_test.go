package httpGet

import "testing"
import (
	"net/http"
)

func Test_CheckAuthorRight(t *testing.T) {
	ud := http.Cookie{Name:"ud", Value:"2"}
	sd := http.Cookie{Name:"sd", Value:"5df4cc9ae659657838802"}
	cookies := []*http.Cookie{&ud, &sd}
	CheckAuthorRight(cookies, 1)
}

func Test_GetLoginUserInfo(t *testing.T) {
	ud := http.Cookie{Name:"ud", Value:"2"}
	sd := http.Cookie{Name:"sd", Value:"5df4cc9ae659657838802"}
	cookies := []*http.Cookie{&ud, &sd}
	GetLoginUserInfo(cookies, 1)
}

func Test_GetUserInfo(t *testing.T) {
	GetUserInfo(1)
}

//go test -file userInfo_test.go -test.bench=".*
func Benchmark_CheckAuthorRight(b *testing.B) {
	cookies := []*http.Cookie{}
	for i := 0; i < b.N; i++ {
		CheckAuthorRight(cookies, 1)
	}
}

func Benchmark_GetLoginUserInfo(b *testing.B) {
	cookies := []*http.Cookie{}
	for i := 0; i < b.N; i++ {
		GetLoginUserInfo(cookies, 1)
	}
}

func Benchmark_GetUserInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetUserInfo(1)
	}
}

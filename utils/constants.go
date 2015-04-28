package utils

const (
	HOST = "http://localhost:8080"
)

type HttpResult struct {
	Error   error       `json:"error"`
	Code    int         `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

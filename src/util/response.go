package util

import (
	"encoding/json"
	"log"
)

// Response api请求返回消息结构
type Response struct {
	Status  int64       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	// StatusSuccess 成功状态
	StatusSuccess = 200
	// StatusFailture 失败状态
	StatusFailture = 400
	// StatusException 异常状态
	StatusException = -1
	// MessageSuccess 成功消息
	MessageSuccess = "Success"
	// MessageFailture 失败消息
	MessageFailture = "Failture"
)

// SuccessResponse 成功返回值
func SuccessResponse(data interface{}) *Response {
	return NewResponse(StatusSuccess, MessageSuccess, data)
}

// FailtureResponse 失败返回值
func FailtureResponse(message string, data interface{}) *Response {
	return NewResponse(StatusFailture, message, data)
}

// ExceptionResponse 异常返回值
func ExceptionResponse(message string, data interface{}) *Response {
	return NewResponse(StatusException, message, data)
}

// NewResponse 创建对象
func NewResponse(status int64, message string, data interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// ToByte 转换成字节
func (r Response) ToByte() []byte {
	marshaled, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	return marshaled
}

// ToString 转换成字符串
func (r Response) ToString() string {
	return string(r.ToByte())
}

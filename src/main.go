// Package main provides ...
package main

import (
	"GoStore/src/handler"
	"log"
	"net/http"
)

func main() {
	// 注册路由
	http.HandleFunc("/file/upload", handler.Upload)
	http.HandleFunc("/success", handler.Success)
	http.HandleFunc("/file/query", handler.QueryFile)
	http.HandleFunc("/file/download", handler.Download)
	http.HandleFunc("/file/update", handler.Update)
	http.HandleFunc("/file/destroy", handler.Destroy)
	http.HandleFunc("/user/signup", handler.SignUp)
	http.HandleFunc("/user/signin", handler.SignIn)
	http.HandleFunc("/user/userinfo", handler.TokenInterceptor(handler.UserInfo))

	// 监听请求
	err := http.ListenAndServe(":8090", nil)

	// 异常处理
	if err != nil {
		log.Fatalf("http server exceptions, %s", err.Error())
	}
}

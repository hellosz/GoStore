// Package main provides ...
package main

import (
	"GoStore/src/handler"
	"log"
	"net/http"
)

func main() {
	// 注册路由
	// 文件接口
	http.HandleFunc("/file/upload", handler.Upload)
	http.HandleFunc("/success", handler.Success)
	http.HandleFunc("/file/query", handler.QueryFile)
	http.HandleFunc("/file/query_list", handler.QueryFileList)
	http.HandleFunc("/file/download", handler.Download)
	http.HandleFunc("/file/update", handler.Update)
	http.HandleFunc("/file/destroy", handler.Destroy)

	// 用户接口
	http.HandleFunc("/user/signup", handler.SignUp)
	http.HandleFunc("/user/signin", handler.SignIn)
	http.HandleFunc("/user/userinfo", handler.TokenInterceptor(handler.UserInfo))

	// 文件秒传接口
	http.HandleFunc("/file/rapid_upload", handler.RapidUpload)

	// 文件分片上传
	http.HandleFunc("file/mp_upload/init", handler.InitMultipartUpload)

	// 监听请求
	err := http.ListenAndServe(":8090", nil)

	// 异常处理
	if err != nil {
		log.Fatalf("http server exceptions, %s", err.Error())
	}
}

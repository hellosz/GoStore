package handler

import (
	"GoStore/src/db"
	"GoStore/src/util"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// SignUp 用户注册接口
func SignUp(w http.ResponseWriter, r *http.Request) {
	// Get 请求进入页面
	if r.Method == http.MethodGet {
		// 读取文件
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			log.Print(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	}

	// Post 请求保存数据
	// 解析参数
	r.ParseForm()
	phone := r.Form.Get("phone")
	password := r.Form.Get("password")

	// 参数校验
	if len(phone) < 3 || len(password) < 6 {
		w.Write([]byte("Invalid Parameters"))
		return
	}

	// 加盐处理
	encPassword := util.EncPassword(password)

	// 保存数据
	ok := db.SignUp(phone, encPassword)
	if !ok {
		w.Write([]byte("SignUp Failture"))
		return
	}

	// 返回结果
	w.Write([]byte("Success"))
}

// SignIn 用户登录接口
func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	}

	// 参数解析
	r.ParseForm()
	phone := r.Form.Get("phone")
	password := r.Form.Get("password")
	encPassword := util.EncPassword(password)

	// 验证账号
	ok := db.SignIn(phone, encPassword)
	if !ok {
		return
	}

	// 生成token
	token, err := db.GenerateToken(phone, encPassword)
	if err != nil {
		// 返回失败结果
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}

	// 返回json 成功结果
	w.Write(util.SuccessResponse(struct {
		Token     string `json:"token"`
		Localtion string `json:"localtion"`
		Timestamp int64  `json:"timestamp"`
	}{
		token,
		"http://localhost:8090/success",
		time.Now().Unix(),
	}).ToByte())
}

// UserInfo 用户信息接口
func UserInfo(w http.ResponseWriter, r *http.Request) {
	// 参数解析及验证
	r.ParseForm()
	token := r.Form.Get("token")

	// // Token 验证
	// ok := db.ValidateToken(token)
	// if !ok {
	// 	w.Write(util.FailtureResponse("invalid token", nil).ToByte())
	// 	return
	// }

	// 取用户信息
	user, err := db.UserInfo(token)
	if err != nil {
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}

	// 返回结果
	w.Write(util.SuccessResponse(user).ToByte())
}

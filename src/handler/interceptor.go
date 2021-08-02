package handler

import (
	"GoStore/src/db"
	"net/http"
)

// TokenInterceptor token验证中间件
func TokenInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 验证token
		r.ParseForm()
		token := r.Form.Get("token")
		ok := db.ValidateToken(token)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		// 调用路由处理器
		h(w, r)
	}
}

// JsonResponse 返回请求头设置未 application/json
func JsonResponse(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置请求头
		w.Header().Add("Content-Type", "application/json")

		// 调用路由处理器
		h(w, r)
	}
}

package db

import (
	"GoStore/src/db/conn"
	"database/sql"
	"errors"
	"log"
	"strings"
)

// SignUp 保存用户信息
func SignUp(phone string, password string) bool {
	// 数据库链接
	conn := conn.GetConn()
	stmt, err := conn.Prepare("insert ignore into tbl_user(phone, user_pwd) values(?, ?)")
	if err != nil {
		log.Print(err.Error())
		return false
	}
	defer stmt.Close()

	// 插入数据
	result, err := stmt.Exec(phone, password)
	if err != nil {
		log.Print(err)
	}
	if affectedRows, err := result.RowsAffected(); err == nil && affectedRows > 0 {
		return true
	}

	// 验证结果
	log.Print("dulplicate insert user information")
	return false
}

// SignIn 用户登录
func SignIn(phone string, encPassword string) bool {
	// 查询用户信息
	sql := "select user_pwd from tbl_user where phone = ? limit 1"
	stmt, err := conn.GetConn().Prepare(sql)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	defer stmt.Close()

	// 验证密码是否正确
	var password string
	_ = stmt.QueryRow(phone).Scan(&password)
	// 比较是否一致
	if strings.Compare(password, encPassword) == 0 {
		return true

	}

	log.Print("password not equal", encPassword, password)
	return false
}

// User 用户信息表
type User struct {
	UserName   string         `json:"user_name"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	Profile    sql.NullString `json:"profile"`
	LastActive string         `json:"last_active"`
}

// UserInfo 根据token获取用户信息
func UserInfo(token string) (*User, error) {
	// 关联查询，根据token获取用户信息
	sql := `
	select a.user_name, a.phone, a.email, a.profile, a.last_active
	from tbl_user as a 
	join tbl_user_token as b on a.phone = b.phone
	where b.token = ?
	limit 1
	`
	stmt, err := conn.GetConn().Prepare(sql)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// 解析用户信息
	var user User
	err = stmt.QueryRow(token).Scan(&user.UserName, &user.Phone, &user.Email, &user.Profile, &user.LastActive)
	if err != nil {
		msg := "invalid token, can't find logined user"
		log.Print(err.Error())
		log.Print(msg)
		return nil, errors.New(msg)
	}

	// 返回用户信息
	return &user, nil
}

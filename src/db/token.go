package db

import (
	"GoStore/src/db/conn"
	"GoStore/src/util"
	"log"
)

// GenerateToken 生成token
func GenerateToken(phone string, encPassword string) (string, error) {
	// 生成Token
	token := util.GenerateToken(encPassword)

	// 保存数据库
	sql := "replace into tbl_user_token(phone, token, status) values(?, ?, 1)"
	stmt, err := conn.GetConn().Prepare(sql)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer stmt.Close()
	_, err = stmt.Exec(phone, token)
	if err != nil {
		log.Print(err)
		return "", err
	}

	// 返回token
	return token, nil
}

// ValidateToken 验证token
func ValidateToken(token string) bool {
	// 参数验证
	if len(token) == 0 {
		log.Print("invalid parameter", token)
		return false
	}

	// 数据库验证
	sql := "select id from tbl_user_token where token = ?"
	stmt, err := conn.GetConn().Prepare(sql)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(token)
	if err != nil {
		log.Print(err.Error())
		return false
	}

	ok := rows.Next()
	if !ok {
		log.Print("invalid token")
		return false

	}

	return true
}

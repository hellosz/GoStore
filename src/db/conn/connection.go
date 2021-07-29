package conn

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var connection *sql.DB

// GetConn 获取数据库链接
func GetConn() *sql.DB {
	if connection != nil {
		return connection
	}

	// 创建链接
	connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "123456", "172.16.13.88", 3307, "store"))
	if err != nil {
		log.Fatal(err)
	}

	// 设置初始化参数
	connection.SetConnMaxIdleTime(3 * time.Minute)
	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(10)

	return connection
}

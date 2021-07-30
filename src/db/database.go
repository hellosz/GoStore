package db

import (
	"GoStore/src/db/conn"
	"database/sql"
	"fmt"
	"log"
)

// OnFileMetaUpdateFinished 更新文件数据库
func OnFileMetaUpdateFinished(fileSha1 string, fileName string, fileSize int64, location string) bool {
	// 获取数据库链接
	conn := conn.GetConn()
	stmt, err := conn.Prepare("insert ignore into tbl_file(file_sha1, file_name, file_size, file_addr, status) " +
		"values (?, ?, ?, ?, 1)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	// 插入数据
	_, err = stmt.Exec(fileSha1, fileName, fileSize, location)

	// 验证结果
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// TblFile 数据表对应的结构类型
type TblFile struct {
	FileSha1 string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// QueryFileMeta 文件元信息查询
func QueryFileMeta(filehash string) (*TblFile, error) {
	// 获取数据库连接
	conn := conn.GetConn()

	// 数据查询
	stmt, err := conn.Prepare("select file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1 = ? limit 1")
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer stmt.Close()

	// 数据赋值
	var file TblFile
	err = stmt.QueryRow(filehash).Scan(&file.FileSha1, &file.FileName, &file.FileSize, &file.FileAddr)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	// 返回结果
	return &file, err
}

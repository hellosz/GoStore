package db

import (
	"GoStore/src/db/conn"
	"log"
)

// UpdateUserFile 更新用户文件信息表
func UpdateUserFile(userId int64, fileSha1 string, fileName string, fileSize int64) bool {
	// 数据库链接
	sql := "insert ignore into tbl_user_file(user_id, file_sha1, file_name, file_size, status) values(?, ? , ?, ?, 1)"
	stmt, err := conn.GetConn().Prepare(sql)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	defer stmt.Close()

	// 验证结果，并返回
	_, err = stmt.Exec(userId, fileSha1, fileName, fileSize)
	if err != nil {
		log.Print(err.Error())
		return false
	}

	return true
}

// UserFile 用户文件信息
type UserFile struct {
	ID        int64  `json:"id"`
	FileSha1  string `json:"file_sha1"`
	FileName  string `json:"file_name"`
	FileSize  string `json:"file_size"`
	UpdatedAt string `json:"updated_at"`
}

// QueryFileList 查询文件列表
func QueryFileList(userId int64) ([]UserFile, error) {
	// 查询文件列表
	sql := "select id, file_sha1, file_name, file_size, update_at from tbl_user_file where user_id = ?"
	stmt, err := conn.GetConn().Prepare(sql)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer stmt.Close()

	// 转换数据
	rows, err := stmt.Query(userId)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var list []UserFile
	for rows.Next() {
		var userFile UserFile
		// 解析数据
		err := rows.Scan(&userFile.ID, &userFile.FileSha1, &userFile.FileName, &userFile.FileSize, &userFile.UpdatedAt)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		// 入队列
		list = append(list, userFile)
	}

	return list, nil
}

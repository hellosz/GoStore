package meta

import (
	"GoStore/src/db"
	"log"
)

// FileMeta 文件元信息
type FileMeta struct {
	FileSha1 string `json:"file_sha1"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Location string `json:"location"`
	UploadAt string `json:"upload_at"`
}

// 定义全局变量
var fileMetas map[string]FileMeta

// init 初始化
func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 取元信息
func UpdateFileMeta(meta FileMeta) {
	fileMetas[meta.FileSha1] = meta
}

// UpdateFileMetaDB 数据库更新文件元信息
func UpdateFileMetaDB(meta FileMeta) bool {
	result := db.OnFileMetaUpdateFinished(meta.FileSha1, meta.FileName, meta.FileSize, meta.Location)
	return result

}

// GetFileMeta 获取文件元信息
func GetFileMeta(sha1 string) FileMeta {
	return fileMetas[sha1]
}

// GetFileMetaDB 从数据库获取文件元信息
func GetFileMetaDB(sha1 string) FileMeta {
	tblFile, err := db.QueryFileMeta(sha1)
	if err != nil {
		log.Fatal(err)
	}

	fileMeta := FileMeta{
		FileSha1: tblFile.FileSha1,
		FileName: tblFile.FileName.String,
		FileSize: tblFile.FileSize.Int64,
		Location: tblFile.FileAddr.String,
	}

	return fileMeta
}

// DestoryFileMeta 删除文件元信息
func DestoryFileMeta(sha1 string) {
	delete(fileMetas, sha1)
}

//
func DestroyFileMeta(sha1 string) {

}

package handler

import (
	"GoStore/src/db"
	"GoStore/src/redis"
	"GoStore/src/util"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// RapidUpload 快速上传
func RapidUpload(w http.ResponseWriter, r *http.Request) {
	// 参数解析
	file, _, err := r.FormFile("file")

	r.ParseForm()
	userId, err := strconv.Atoi(r.Form.Get("userId"))
	if err != nil {
		log.Print(err)
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}

	tmpFile, err := ioutil.TempFile("", "tmpFile")
	if err != nil {
		log.Print(err)
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}

	io.Copy(tmpFile, file)
	tmpFile.Seek(0, 0)
	fileSha1 := util.FileSha1(tmpFile)

	// 判断文件是否存在，存在则快速上传
	log.Println(fileSha1)
	fileMeta, err := db.QueryFileMeta(fileSha1)
	if err != nil {
		log.Print(err)
		w.Write(util.FailtureResponse("not rapit upload, please request normal upload uri", nil).ToByte())
		return
	}

	// 保存用户文件关系数据
	ok := db.UpdateUserFile(int64(userId), fileMeta.FileSha1, fileMeta.FileName.String, fileMeta.FileSize.Int64)
	if ok {
		log.Print("rapit upload")
		w.Write(util.SuccessResponse(nil).ToByte())
		return
	}

	// 走普通上err.Error(), nil传
	// http.Redirect(w, r, "/file/upload", http.StatusMovedPermanently)
	log.Println(6)
}

// MultipartUploadFile 分片文件上传结构体
type MultipartUploadFile struct {
	Filename  string
	Filehash  string
	FileSize  int64
	TraceID   string
	Chunks    int64
	ChunkSize int64
}

const (
	// ChunkSize 默认的分块大小
	ChunkSize = 4 * 1024 * 1024
)

// InitMultipartUpload 初始化分片文件上传
func InitMultipartUpload(w http.ResponseWriter, r *http.Request) {
	// 参数初始化
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 创建Redis连接
	redis := redis.NewClient()
	defer redis.Close()

	// 文件信息初始化
	mtFile := MultipartUploadFile{
		FileSize:  int64(filesize),
		Filehash:  filehash,
		TraceID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize: ChunkSize,
		Chunks:    int64(math.Ceil(float64(filesize) / ChunkSize)),
	}

	// 缓存分片信息
	redis.Do("HSET", mtFile.TraceID, "filesize", mtFile.FileSize)
	redis.Do("HSET", mtFile.TraceID, "filehash", mtFile.Filehash)
	redis.Do("HSET", mtFile.TraceID, "chunkSize", mtFile.ChunkSize)
	redis.Do("HSET", mtFile.TraceID, "chunks", mtFile.Chunks)

	// 返回结果
	w.Write(util.SuccessResponse(mtFile).ToByte())
}

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
	"os"
	"path/filepath"
	"strconv"
	"time"

	redisgo "github.com/gomodule/redigo/redis"
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
	Filename  string `json:"filename" redis:"filename"`
	FileHash  string `json:"file_hash" redis:"file_hash"`
	FileSize  int64  `json:"file_size" redis:"file_size"`
	TraceID   string `json:"trace_id" redis:"trace_id"`
	Chunks    int64  `json:"chunks" redis:"chunks"`
	ChunkSize int64  `json:"chunk_size" redis:"chunk_size"`
}

const (
	// ChunkSize 默认的分块大小
	// ChunkSize = 4 * 1024 * 1024
	ChunkSize = 10
)

// InitMultipartUpload 初始化分片文件上传
func InitMultipartUpload(w http.ResponseWriter, r *http.Request) {
	// 参数初始化
	r.ParseForm()
	filename := r.Form.Get("filename")
	filehash := r.Form.Get("filehash")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 创建Redis连接
	redis := redis.NewClient()
	defer redis.Close()

	// 文件信息初始化
	mtFile := MultipartUploadFile{
		Filename:  filename,
		FileSize:  int64(filesize),
		FileHash:  filehash,
		TraceID:   fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize: ChunkSize,
		Chunks:    int64(math.Ceil(float64(filesize) / ChunkSize)),
	}

	// 缓存分片信息
	redis.Do("HSET", mtFile.TraceID, "filename", mtFile.Filename)
	redis.Do("HSET", mtFile.TraceID, "file_size", mtFile.FileSize)
	redis.Do("HSET", mtFile.TraceID, "file_hash", mtFile.FileHash)
	redis.Do("HSET", mtFile.TraceID, "chunk_size", mtFile.ChunkSize)
	redis.Do("HSET", mtFile.TraceID, "chunks", mtFile.Chunks)
	redis.Do("HSET", mtFile.TraceID, "trace_id", mtFile.TraceID)

	// 返回结果
	w.Write(util.SuccessResponse(mtFile).ToByte())
}

const (
	// multipartUploadDir 分片上传目录
	multipartUploadDir = "/var/www/tmp/files"
)

//MultipartUpload 分块文件上传
func MultipartUpload(w http.ResponseWriter, r *http.Request) {
	// 参数解析
	r.ParseForm()
	// username := r.Form.Get("username")
	traceId := r.Form.Get("trace_id")
	chunkIndex := r.Form.Get("chunk_index")

	// 获取Redis链接
	redi := redis.NewClient()

	// 保存本地文件

	fpath := filepath.Join(multipartUploadDir, traceId, chunkIndex)
	// 创建文件目录
	err := os.MkdirAll(filepath.Dir(fpath), 0744)
	if err != nil {
		log.Print(err)
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}

	// 保存文件
	file, err := os.Create(fpath)
	if err != nil {
		log.Print(err)
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}
	defer file.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)

		fmt.Println(string(buf[:n]), "buf", err)
		if err != nil {
			break
		}

		// 将请求内容写入到本地创建的分块文件中
		fmt.Println(string(buf[:n]), "=")
		_, err = file.Write(buf[:n])
		if err != nil {
			log.Print(err)
			w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
			return
		}
	}

	// 保存Redis
	redi.Do("RPUSH", traceId+"_index", chunkIndex)

	// 返回结果
	w.Write(util.SuccessResponse(nil).ToByte())
}

// CompleteMultipartUpload 文件上传完成
func CompleteMultipartUpload(w http.ResponseWriter, r *http.Request) {
	// 参数解析
	r.ParseForm()
	traceId := r.Form.Get("trace_id")

	// 校验是否完成
	// 读取redis中所有存储的信息，然后校验是否完成（校验方式：分块数量是否一致，以及hash值是否一致）
	redi := redis.NewClient()
	values, err := redisgo.Values(redi.Do("HGETALL", traceId))
	if err != nil {
		log.Print(err)
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}

	var mpFile = MultipartUploadFile{}
	redisgo.ScanStruct(values, &mpFile)

	// 读取已经上的分片信息
	values, err = redisgo.Values(redi.Do("LRANGE", traceId+"_index", 0, -1))
	if err != nil {
		log.Print(err)
		w.Write(util.FailtureResponse(err.Error(), nil).ToByte())
		return
	}

	var chunkIndexs []string
	redisgo.ScanSlice(values, &chunkIndexs)
	// 校验所有分块是是否上传完成
	if mpFile.Chunks != int64(len(chunkIndexs)) {
		w.Write(util.FailtureResponse("some chunks doesn't upload yet", nil).ToByte())
		return
	}

	// 读取文件信息 TODO
	// 读取目录下所有文件

	// 文件合并

	// 保存文件信息
	ok := db.OnFileMetaUpdateFinished(mpFile.FileHash, mpFile.Filename, mpFile.FileSize, "")
	if !ok {
		w.Write(util.FailtureResponse("save file meta failed", nil).ToByte())
		return
	}

	// 保存用户文件信息
	db.UpdateUserFile(9, mpFile.FileHash, mpFile.Filename, mpFile.FileSize)
	if !ok {
		w.Write(util.FailtureResponse("save user file failed", nil).ToByte())
		return
	}

	// 返回结果
	w.Write(util.SuccessResponse(nil).ToByte())
}

// ValidMultiparUploadFile
func ValidMultiparUploadFile(traceId string) {

}

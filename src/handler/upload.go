package handler

import (
	"GoStore/src/meta"
	"GoStore/src/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Upload 上传处理方法
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		view, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(view)
	} else if r.Method == "POST" {
		// 获取文件
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Fatal(err)
		}

		// 初始化文件元信息
		fileMeta := meta.FileMeta{
			FileName: fileHeader.Filename,
			Location: filepath.Join("/tmp", fileHeader.Filename),
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 创建本地文件
		localFile, err := os.Create(fileMeta.Location)
		if err != nil {
			log.Fatalf("create local file failed:%s", err.Error())
		}

		// 文件保存到本地
		fileMeta.FileSize, err = io.Copy(localFile, file)
		if err != nil {
			log.Fatalf("copy file to local failed:%s", err.Error())
		}

		// 获取文件的sha1信息
		localFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(localFile)
		ok := meta.UpdateFileMetaDB((fileMeta))
		if !ok {
			log.Fatalf("save filemeta to database failed")

		}
		// meta.UpdateFileMeta(fileMeta)

		// 页面成功显示
		// w.Write([]byte("method post handler"))
		r.Form.Add("action", "upload") // 添加参数
		http.Redirect(w, r, "/success", http.StatusMovedPermanently)
	}
}

// Success 成功页面
func Success(w http.ResponseWriter, r *http.Request) {
	action := r.Form.Get("action")
	fmt.Fprintf(w, "%s success", action)
}

// QueryFile 查询文件
func QueryFile(w http.ResponseWriter, r *http.Request) {
	// 读取参数
	r.ParseForm()
	filehash := r.Form["filehash"][0]

	// 查找本地文件
	// fileMeta := meta.GetFileMeta(filehash)
	fileMeta := meta.GetFileMetaDB(filehash)

	data, err := json.Marshal(fileMeta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 返回结果
	w.Write(data)
}

// Download 下载指定文件
func Download(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	// 查找文件
	fileMeta := meta.GetFileMeta(filehash)
	data, err := ioutil.ReadFile(fileMeta.Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 设置响应头，返回文件流
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}

// Update 更新文件信息
func Update(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	r.ParseForm()
	op := r.Form.Get("op")
	filehash := r.Form.Get("filehash")
	newFilename := r.Form.Get("filename")

	// 参数校验
	if op != "1" { // 不是指定的操作类型，返回报错
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" { // 仅支持POST方法
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 修改文件名称
	fileMeta := meta.GetFileMeta(filehash)
	fileMeta.FileName = newFilename
	meta.UpdateFileMeta(fileMeta)

	// 返回最新的结果
	data, err := json.Marshal(fileMeta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Destroy 删除文件信息
func Destroy(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	// 删除本地文件
	fileMeta := meta.GetFileMeta(filehash)
	os.Remove(fileMeta.Location)

	// 删除内存中的信息
	meta.DestoryFileMeta(filehash)

	// 返回成功提示
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete success"))
}

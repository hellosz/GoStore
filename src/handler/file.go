package handler

import (
	"GoStore/src/db"
	"GoStore/src/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

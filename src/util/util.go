package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Sha1Stream 流对象
type Sha1Stream struct {
	_sha1 hash.Hash
}

// Update 更新操作
func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// 密码需要的盐
const (
	passwordSalt = "^&%JHFDA77423<>:"
	tokenSalt    = "token_salt_^&%^$"
)

// EncPassword 加密密码
func EncPassword(password string) string {
	return Sha1([]byte(password + passwordSalt))
}

// GenerateToken 生成token
func GenerateToken(encPassword string) string {
	// token组成：（密码 + 时间戳 + 盐）+ 时间戳前8位
	timeString := fmt.Sprintf("%x", time.Now().Unix())
	return MD5([]byte(encPassword+timeString+tokenSalt)) + timeString[:8]
}

package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// FileCtlBasic 文件操作
type FileCtlBasic struct{}

// FileCtl 数据转换
var FileCtl *FileCtlBasic

func init() {
	FileCtl = NewFileCtl()
}

// NewFileCtl 初始化 FileCtl
func NewFileCtl() (n *FileCtlBasic) {
	n = new(FileCtlBasic)
	return
}

// Write 写入文件
func (wr *FileCtlBasic) Write(f string, data []byte) {
	err := ioutil.WriteFile(f, []byte(data), 0666)
	if err != nil {
		log.Panic(err)
	}
}

// Read 读取文件
func (wr *FileCtlBasic) Read(f string) []byte {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Panic(err)
	}
	return data
}

// Create 创建文件夹
func (wr *FileCtlBasic) Create(folderPath string) string {
	_, err := os.Stat(folderPath)
	if err == nil {
		return folderPath
	}
	if os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
		return folderPath
	}
	return folderPath
}

// CheckType 检查文件类型
func (wr *FileCtlBasic) CheckType(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}

// CheckExist 检查文件存在
func (wr *FileCtlBasic) CheckExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// LocalPath 获取当前路径
func (wr *FileCtlBasic) LocalPath(mode int) (path string) {
	switch {
	case mode == 1:
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Panic(err)
		}
		path = strings.Replace(dir, "\\", "/", -1)
	case mode == 0:
		dir, _ := os.Getwd()
		path = strings.Replace(dir, "\\", "/", -1)
	}
	return
}

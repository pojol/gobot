package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentDirectory 获取进程的当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetDirectoryFiels(dir string, ext string) []string {

	files := []string{}

	infoLst, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range infoLst {
		if GetFileExt(v.Name()) == ext {
			files = append(files, v.Name())
		}
	}

	return files

}

// GetFileExt 获取文件的扩展名
func GetFileExt(fileName string) string {
	if fileName == "" {
		return ""
	}

	index := strings.LastIndex(fileName, ".")
	if index < 0 {
		return ""
	}

	return string(fileName[index:])

}

// GetFileRealName 获取文件名
func GetFileRealName(fileName string) string {
	if fileName == "" {
		return ""
	}

	idx := strings.LastIndex(fileName, ".")
	if idx < 0 {
		return ""
	}
	return string(fileName[:idx])
}

// Exist 判断文件是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// ReadFile 读取一个文件
func ReadFile(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// SaveFile 保存文件
func SaveFile(dat []byte, fileName string) error {

	// 0666 filemode
	err := ioutil.WriteFile(fileName, dat, 0666)

	return err
}

// WriteJSON 保存json文件
func WriteJSON(path string, jsonByte []byte) {
	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fp.Close()
	_, err = fp.Write(jsonByte)
	if err != nil {
		log.Fatal(err)
	}
}

package models

import (
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"strings"
)

func GetFileList(pathname string) []string {
	var res []string
	basePath := pathname + "/"
	res = GetAllFile(pathname, basePath)
	return res
}

// 递归读取目录下的所有文件
func GetAllFile(pathname string, basePath string) []string {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		logs.Error(err)
	}

	var res []string
	for _, fi := range rd {
		filename := fi.Name()
		if fi.IsDir() {
			tmp := GetAllFile(pathname+"/"+filename, basePath)
			res = append(res, tmp...)
		} else {
			t := pathname + "/" + filename
			t = strings.Replace(t, basePath, "", -1)
			res = append(res, t)
		}
	}
	return res
}

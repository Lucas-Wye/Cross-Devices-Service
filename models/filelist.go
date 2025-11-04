package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego/logs"
)

type File struct {
	Name    string
	Size    int64
	ModTime time.Time
	Path    string // Relative path to base
}

func GetFileList(pathname string) []File {
	var res []File
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		logs.Error(err)
		return res
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			res = append(res, File{
				Name:    fi.Name(),
				Size:    fi.Size(),
				ModTime: fi.ModTime(),
			})
		}
	}
	return res
}

// GetAllFile is deprecated, use GetFileList and walk directories in controller if needed.
func GetAllFile(pathname string, basePath string) []string {
	var res []string
	filepath.Walk(pathname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}
			res = append(res, relPath)
		}
		return nil
	})
	return res
}

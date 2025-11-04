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

// FileNode 表示用于前端树形展示的节点
// 目录节点 IsDir=true，Name 为目录名，RelDir 为该目录相对根的路径（例如 "docs/sub"），Children 包含子节点
// 文件节点 IsDir=false，Name 为文件名，RelDir 为包含它的目录的相对路径（用于下载时 dir 字段）
type FileNode struct {
	Name     string     `json:"name"`
	RelDir   string     `json:"relDir"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children,omitempty"`
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

// GetFileTree 构建以 relRoot（相对 basePath 的路径）为根目录的树
// relRoot 通常是顶层权限目录名（例如用户权限 Path），也可以是更深的相对路径
func GetFileTree(basePath, relRoot string) (*FileNode, error) {
	absRoot := filepath.Join(basePath, relRoot)
	fi, err := os.Stat(absRoot)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return &FileNode{
			Name:   fi.Name(),
			RelDir: filepath.Dir(relRoot),
			IsDir:  false,
		}, nil
	}
	node := &FileNode{
		Name:   filepath.Base(absRoot),
		RelDir: relRoot,
		IsDir:  true,
	}
	entries, err := ioutil.ReadDir(absRoot)
	if err != nil {
		return node, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			childRel := filepath.Join(relRoot, name)
			child, err := GetFileTree(basePath, childRel)
			if err == nil && child != nil {
				node.Children = append(node.Children, *child)
			}
		} else {
			node.Children = append(node.Children, FileNode{
				Name:   name,
				RelDir: relRoot,
				IsDir:  false,
			})
		}
	}
	return node, nil
}

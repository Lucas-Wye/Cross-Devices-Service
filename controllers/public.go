package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PublicFileController struct {
	beego.Controller
}

func (this *PublicFileController) PrivateUpload() {
	overwrite := this.GetString("overwrite")
	dir := this.GetString("dir")
	// fmt.Println("[PublicFileController] dir =", dir, "overwrite =", overwrite)

	// 检查 overwrite 参数
	if overwrite != "true" && overwrite != "false" {
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   "overwrite 参数必须是 'true' 或 'false'",
		}
		this.ServeJSON()
		return
	}

	f, h, err := this.GetFile("uploadfile") // 读取上传文件
	if err != nil {
		// 直接返回 JSON，而不是渲染 HTML
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   "未选择文件或读取文件失败",
			"detail":  err.Error(),
		}
		this.ServeJSON()
		return
	}
	defer f.Close()

	targetDir := DownloadPath
	if dir != "" {
		targetDir = filepath.Join(DownloadPath, dir)
		// fmt.Println("[82] Target directory:", targetDir)
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			os.MkdirAll(targetDir, 0755)
		}
	}

	// 计算目标文件名
	finalName := h.Filename
	destPath := filepath.Join(targetDir, finalName)

	if overwrite == "true" {
		// 允许覆盖，直接使用原名
	} else {
		// 若已存在则自动加后缀(1),(2)...
		if _, statErr := os.Stat(destPath); statErr == nil {
			base := strings.TrimSuffix(finalName, filepath.Ext(finalName))
			ext := filepath.Ext(finalName)
			for i := 1; ; i++ {
				candidate := fmt.Sprintf("%s(%d)%s", base, i, ext)
				candidatePath := filepath.Join(targetDir, candidate)
				if _, err := os.Stat(candidatePath); os.IsNotExist(err) {
					finalName = candidate
					destPath = candidatePath
					break
				}
			}
		}
	}

	if err := this.SaveToFile("uploadfile", destPath); err != nil {
		this.Ctx.Output.SetStatus(500)
		this.Data["json"] = map[string]interface{}{
			"success": false,
			"error":   "保存文件失败",
			"detail":  err.Error(),
		}
		this.ServeJSON()
		return
	}

	// 成功：返回 JSON，不做重定向、不渲染模板
	relPath := finalName
	if dir != "" {
		relPath = filepath.ToSlash(filepath.Join(dir, finalName))
	}
	this.Data["json"] = map[string]interface{}{
		"success":      true,
		"filename":     finalName,
		"dir":          dir,
		"path":         relPath, // 相对 DownloadPath 的路径
		"size":         h.Size,
		"content_type": h.Header.Get("Content-Type"),
		"upload_time":  time.Now().Format("2006-01-02 15:04:05"),
	}
	this.ServeJSON()
	return
}

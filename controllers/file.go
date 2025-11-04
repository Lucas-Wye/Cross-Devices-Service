package controllers

import (
	models "CrossDevicesService/models"

	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var DownloadPath = models.GetLocalDirPath()

const CopyPastePath = "/copy"
const ServiceName = "Cross Devices Service"
const CopyPasteFile = "./logs/copy_paste.txt"

type FileController struct {
	beego.Controller
	user *models.User
}

func (this *FileController) Prepare() {
	sessionUser := this.GetSession("user")
	if sessionUser == nil {
		this.Redirect("/login", 302)
		this.StopRun()
	}

	// It's better to store just the username in session and refetch user data
	// to ensure permissions are always up-to-date.
	user, ok := models.GetUser(sessionUser.(*models.User).Username)
	if !ok {
		this.Redirect("/login", 302)
		this.StopRun()
	}
	this.user = user
}

func (c *FileController) GetUpload() {
	// 读取上传成功的提示参数
	if c.GetString("success") == "1" {
		c.Data["success"] = true
		c.Data["filename"] = c.GetString("filename")
	}
	c.Data["user"] = c.user
	// 如果之前在同一请求中设置了错误信息，直接展示
	c.TplName = "upload.html"
}

func (this *FileController) Upload() {
	dir := this.GetString("dir")
	// fmt.Println("[55] dir =", dir)
	hasPermission := false
	for _, p := range this.user.Permissions {
		if (dir == "" && p.Path == "shared") || (dir == p.Path) {
			if p.Write {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		this.Data["error"] = "没有上传权限"
		this.TplName = "upload.html"
		return
	}

	f, h, err := this.GetFile("uploadfile") // 读取上传文件
	if err != nil {
		this.Data["error"] = "未选择文件或读取文件失败"
		this.TplName = "upload.html"
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

	// 计算目标文件名，若已存在则自动加后缀(1),(2)...
	finalName := h.Filename
	destPath := filepath.Join(targetDir, finalName)
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

	if err := this.SaveToFile("uploadfile", destPath); err != nil {
		this.Data["error"] = "保存文件失败: " + err.Error()
		this.TplName = "upload.html"
		return
	}

	// 成功后重定向到上传页并附带成功提示（避免重复提交）
	redirectURL := "/upload?success=1&filename=" + url.QueryEscape(finalName)
	if dir != "" {
		redirectURL += "&dir=" + url.QueryEscape(dir)
	}
	this.Redirect(redirectURL, 302)
}

func (this *FileController) GetList() {
	var trees []models.FileNode
	for _, p := range this.user.Permissions {
		if p.Read {
			dirPath := filepath.Join(DownloadPath, p.Path)
			if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
				if root, err := models.GetFileTree(DownloadPath, p.Path); err == nil && root != nil {
					trees = append(trees, *root)
				}
			}
		}
	}
	this.Data["trees"] = trees
	this.TplName = "download.html"
}

func (this *FileController) Download() {
	dir := this.GetString("dir")
	filename := this.GetString("filename")
	hasPermission := false
	for _, p := range this.user.Permissions {
		if ((dir == "" && p.Path == "shared") || (dir == p.Path)) && p.Read {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		this.Abort("403")
		return
	}

	targetPath := filepath.Join(DownloadPath, dir, filename)
	this.Ctx.Output.Download(targetPath)
}

func (c *FileController) Paste() {
	// c.Data["input"] = models.ReadFromCopyPaste()
	c.Data["input"] = models.ReadFromFile(CopyPasteFile)
	c.TplName = "copy_paste.html"
}

func (c *FileController) Copy() {
	input := c.GetString("input")
	if input != "" {
		// models.Write2CopyPaste(input)
		models.WriteToFile(CopyPasteFile, "\n"+input)
	}
	c.Redirect(CopyPastePath, 302)
}

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
	userRole string
}

func (this *FileController) Prepare() {
	a := NewBasicAuthenticator(ServiceName, Secret)
	if username := a.CheckAuth(this.Ctx.Request); username == "" {
		a.RequireAuth(this.Ctx.ResponseWriter, this.Ctx.Request)
	} else {
		if username == models.GetAdminUsername() {
			this.userRole = "admin"
		} else {
			this.userRole = "normal"
		}
	}
}

func (c *FileController) GetUpload() {
	// 读取上传成功的提示参数
	if c.GetString("success") == "1" {
		c.Data["success"] = true
		c.Data["filename"] = c.GetString("filename")
	}
	// 如果之前在同一请求中设置了错误信息，直接展示
	c.TplName = "upload.html"
}

func (this *FileController) Upload() {
	f, h, err := this.GetFile("uploadfile") // 读取上传文件
	if err != nil {
		this.Data["error"] = "未选择文件或读取文件失败"
		this.TplName = "upload.html"
		return
	}
	defer f.Close()

	// 计算目标文件名，若已存在则自动加后缀(1),(2)...
	finalName := h.Filename
	destPath := filepath.Join(DownloadPath, finalName)
	if _, statErr := os.Stat(destPath); statErr == nil {
		base := strings.TrimSuffix(finalName, filepath.Ext(finalName))
		ext := filepath.Ext(finalName)
		for i := 1; ; i++ {
			candidate := fmt.Sprintf("%s(%d)%s", base, i, ext)
			candidatePath := filepath.Join(DownloadPath, candidate)
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
	this.Redirect("/upload?success=1&filename="+url.QueryEscape(finalName), 302)
}

func (this *FileController) GetList() {
	if this.userRole != "admin" {
		this.Abort("403")
		return
	}
	this.Data["list"] = models.GetFileList(DownloadPath)
	this.TplName = "download.html"
}

func (this *FileController) Download() {
	if this.userRole != "admin" {
		this.Abort("403")
		return
	}
	filename := this.GetString("filename")
	this.Ctx.Output.Download(DownloadPath + "/" + filename)
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

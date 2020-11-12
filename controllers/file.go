package controllers

import (
	models "CrossDevicesService/models"

	auth "github.com/abbot/go-http-auth"
	"github.com/astaxie/beego"
)

var DownloadPath = models.GetLocalDirPath()

const CopyPastePath = "/copy"
const ServiceName = "Cross Devices Service"

type FileController struct {
	beego.Controller
}

func (this *FileController) Prepare() {
	a := auth.NewBasicAuthenticator(ServiceName, Secret)
	if username := a.CheckAuth(this.Ctx.Request); username == "" {
		a.RequireAuth(this.Ctx.ResponseWriter, this.Ctx.Request)
	}
}

func (c *FileController) GetUpload() {
	c.TplName = "upload.html"
}

func (this *FileController) Upload() {
	f, h, _ := this.GetFile("uploadfile")   // get the files
	path := DownloadPath + "/" + h.Filename // save path
	f.Close()
	this.SaveToFile("uploadfile", path)
	this.Redirect("/upload", 302)
}

func (this *FileController) GetList() {
	this.Data["list"] = models.GetFileList(DownloadPath)
	this.TplName = "download.html"
}

func (this *FileController) Download() {
	filename := this.GetString("filename")
	this.Ctx.Output.Download(DownloadPath + "/" + filename)
}

func (c *FileController) Paste() {
	c.Data["input"] = models.ReadFromCopyPaste()
	c.TplName = "copy_paste.html"
}

func (c *FileController) Copy() {
	input := c.GetString("input")
	if input != "" {
		models.Write2CopyPaste(input)
	}
	c.Redirect(CopyPastePath, 302)
}

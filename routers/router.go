package routers

import (
	"FileService/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	// files
	beego.Router("/upload", &controllers.FileController{}, "get:GetUpload")
	beego.Router("/upload", &controllers.FileController{}, "post:Upload")
	beego.Router("/download", &controllers.FileController{}, "get:GetList")
	beego.Router("/download", &controllers.FileController{}, "post:Download")

	// words
	beego.Router(controllers.CopyPastePath, &controllers.FileController{}, "get:Paste")
	beego.Router(controllers.CopyPastePath, &controllers.FileController{}, "post:Copy")
}

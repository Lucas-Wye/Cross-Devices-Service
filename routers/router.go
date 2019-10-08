package routers

import (
	"upper/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")

	// 文件
	beego.Router("/upload", &controllers.FileController{}, "get:GetUpload")
	beego.Router("/upload", &controllers.FileController{}, "post:Upload")
	beego.Router("/download", &controllers.FileController{}, "get:GetList")
	beego.Router("/download", &controllers.FileController{}, "post:Download")

	// 上传文件目录
	beego.SetStaticPath("./upload", "upload")

}

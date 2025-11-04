package routers

import (
	"CrossDevicesService/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/login", &controllers.LoginController{}, "get:Get;post:Post")
	beego.Router("/logout", &controllers.MainController{}, "get:Logout")
	// files
	beego.Router("/upload", &controllers.FileController{}, "get:GetUpload")
	beego.Router("/upload", &controllers.FileController{}, "post:Upload")
	beego.Router("/download", &controllers.FileController{}, "get:GetList")
	beego.Router("/download", &controllers.FileController{}, "post:Download")

	// words
	beego.Router(controllers.CopyPastePath, &controllers.FileController{}, "get:Paste")
	beego.Router(controllers.CopyPastePath, &controllers.FileController{}, "post:Copy")

	// admin
	beego.Router("/admin", &controllers.AdminController{}, "get:Get")
	beego.Router("/admin/add_user", &controllers.AdminController{}, "post:AddUser")
	beego.Router("/admin/delete_user", &controllers.AdminController{}, "post:DeleteUser")
	beego.Router("/admin/update_permissions", &controllers.AdminController{}, "post:UpdatePermissions")
}

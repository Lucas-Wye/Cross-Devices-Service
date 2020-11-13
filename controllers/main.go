package controllers

import (
	models "CrossDevicesService/models"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func Secret(user, realm string) string {
	if user == models.GetLocalUsername() {
		return models.GetLocalPassword()
	}
	return ""
}

func (this *MainController) Prepare() {
	a := NewBasicAuthenticator(ServiceName, Secret)
	if a.CheckAuth(this.Ctx.Request) == "" {
		a.RequireAuth(this.Ctx.ResponseWriter, this.Ctx.Request)
	}
}

func (this *MainController) Get() {
	this.TplName = "index.html"
}

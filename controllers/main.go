package controllers

import (
	models "CrossDevicesService/models"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func Secret(username, realm string) string {
	if user, ok := models.GetUser(username); ok {
		return user.PasswordHash
	}
	return ""
}

func (this *MainController) Prepare() {
	user := this.GetSession("user")
	if user == nil {
		this.Redirect("/login", 302)
		return
	}
}

func (this *MainController) Get() {
	this.TplName = "index.html"
}

func (c *MainController) Logout() {
	c.DelSession("user")
	c.Redirect("/login", 302)
}

package controllers

import (
	auth "github.com/abbot/go-http-auth"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func Secret(user, realm string) string {
	if user == "Wye" {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

func (this *MainController) Prepare() {
	a := auth.NewBasicAuthenticator("example.com", Secret)
	if username := a.CheckAuth(this.Ctx.Request); username == "" {
		a.RequireAuth(this.Ctx.ResponseWriter, this.Ctx.Request)
	}
}

func (this *MainController) Get() {
	this.TplName = "index.html"
}

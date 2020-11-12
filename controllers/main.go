package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func Secret(user, realm string) string {
	if user == "Wye" {
		// password is "hello"
		// return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
		return "$2a$10$zVeDUQ6CdmzQK55iojloiecJEoHz2qW7AMvIb19JXQ/kRfRFe7s.O"
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

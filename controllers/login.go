package controllers

import (
	"CrossDevicesService/models"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	username := c.GetString("username")
	password := c.GetString("password")

	user, ok := models.GetUser(username)
	if !ok {
		c.Data["error"] = "用户名或密码错误"
		c.TplName = "login.html"
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		c.Data["error"] = "用户名或密码错误"
		c.TplName = "login.html"
		return
	}

	c.SetSession("user", user)
	c.Redirect("/", 302)
}

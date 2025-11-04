package controllers

import (
	"CrossDevicesService/models"
	// "fmt"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/bcrypt"
)

type AdminController struct {
	beego.Controller
	user *models.User
}

func (this *AdminController) Prepare() {
	sessionUser := this.GetSession("user")
	if sessionUser == nil {
		this.Redirect("/login", 302)
		this.StopRun()
	}

	user, ok := models.GetUser(sessionUser.(*models.User).Username)
	if !ok || user.Role != "admin" {
		this.Redirect("/", 302) // Redirect non-admins to home
		this.StopRun()
	}
	this.user = user
}

func (this *AdminController) Get() {
	config := models.GetUserConfig()
	this.Data["Users"] = config.Users
	this.TplName = "admin.html"
}

func (this *AdminController) AddUser() {
	username := this.GetString("username")
	password := this.GetString("password")
	role := this.GetString("role")

	if username == "" || password == "" || role == "" {
		this.Redirect("/admin?error=missing_fields", 302)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		this.Redirect("/admin?error=bcrypt_failed", 302)
		return
	}

	newUser := models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		Permissions:  []models.Permission{},
	}

	config := models.GetUserConfig()
	config.Users = append(config.Users, newUser)
	if err := models.SaveUserConfig(config); err != nil {
		this.Redirect("/admin?error=save_failed", 302)
		return
	}

	this.Redirect("/admin", 302)
}

func (this *AdminController) DeleteUser() {
	username := this.GetString("username")
	if username == "" || username == "admin" {
		this.Redirect("/admin?error=cannot_delete", 302)
		return
	}

	config := models.GetUserConfig()
	var updatedUsers []models.User
	for _, u := range config.Users {
		if u.Username != username {
			updatedUsers = append(updatedUsers, u)
		}
	}
	config.Users = updatedUsers
	if err := models.SaveUserConfig(config); err != nil {
		this.Redirect("/admin?error=save_failed", 302)
		return
	}

	this.Redirect("/admin", 302)
}

func (this *AdminController) UpdatePermissions() {
	username := this.GetString("username")
	if username == "" {
		this.Redirect("/admin?error=missing_username", 302)
		return
	}

	config := models.GetUserConfig()
	var userToUpdate *models.User
	for i := range config.Users {
		if config.Users[i].Username == username {
			userToUpdate = &config.Users[i]
			break
		}
	}

	if userToUpdate == nil {
		this.Redirect("/admin?error=user_not_found", 302)
		return
	}

	// Clear existing permissions
	userToUpdate.Permissions = []models.Permission{}

	// Get permissions from form
	paths := this.GetStrings("path[]")
	reads := this.GetStrings("read[]")
	writes := this.GetStrings("write[]")
	// fmt.Println("Updating permissions for user:", username, "path =", paths, "read =", reads, "write =", writes)

	readMap := make(map[string]bool)
	for _, r := range reads {
		readMap[r] = true
	}
	writeMap := make(map[string]bool)
	for _, w := range writes {
		writeMap[w] = true
	}

	for _, path := range paths {
		if path != "" {
			perm := models.Permission{
				Path:  path,
				Read:  readMap[path],
				Write: writeMap[path],
			}
			userToUpdate.Permissions = append(userToUpdate.Permissions, perm)
		}
	}

	if err := models.SaveUserConfig(config); err != nil {
		this.Redirect("/admin?error=save_failed", 302)
		return
	}

	this.Redirect("/admin", 302)
}

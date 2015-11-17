package controllers

import (
	"github.com/astaxie/beego"

	"beegoblog/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Get() {

	this.Data["IsUser"] = true
	this.TplNames = "user.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

}

func (this *UserController) Load() {
	var users []*(models.User)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		users, err = models.GetAllUsersRedis(false)
	default:
		users, err = models.GetAllUsers(false)
	}

	if err != nil {
		beego.Error(err)
	}

	this.Data["json"] = &users
	this.ServeJson()
}

func (this *UserController) Post() {

	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	// 解析表单
	uname := this.Input().Get("uname")
	pwd := this.Input().Get("pwd")
	isView := this.Input().Get("isView")
	isAdd := this.Input().Get("isAdd")
	var err error
	if isAdd == "true" {
		switch beego.AppConfig.String("database") {
		case "redis":
			err = models.AddUserRedis(uname, pwd)
		default:
			err = models.AddUser(uname, pwd)
		}

	} else if isView == "true" {
		switch beego.AppConfig.String("database") {
		case "redis":
			err = models.ModifyUserRedis(uname, pwd)
		default:
			err = models.ModifyUser(uname, pwd)
		}

	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/user", 302)
}

func (this *UserController) Add() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	this.TplNames = "user_add.html"
	this.Data["IsLogin"] = true
}

func (this *UserController) Delete() {

	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		err = models.DeleteUserRedis(this.Input().Get("uname"))
	default:
		err = models.DeleteUser(this.Input().Get("uname"))
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/user", 302)
}

func (this *UserController) Modify() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	this.Data["uname"] = this.Input().Get("uname")
	this.TplNames = "user_modify.html"
	this.Data["IsLogin"] = true

}

func (this *UserController) LoadModify() {

	uname := this.Input().Get("uname")
	var user *(models.User)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		user, err = models.GetUserRedis(uname)
	default:
		user, err = models.GetUser(uname)
	}

	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	this.Data["json"] = &user
	this.ServeJson()

}

func (this *UserController) View() {
	uname := this.Input().Get("uname")
	this.Data["uname"] = uname

	this.TplNames = "user_view.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

}

func (this *UserController) LoadView() {
	uname := this.Input().Get("uname")
	var user *(models.User)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		user, err = models.GetUserRedis(uname)
	default:
		user, err = models.GetUser(uname)
	}

	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	data := &struct {
		User *models.User
	}{
		User: user,
	}

	this.Data["json"] = data
	this.ServeJson()
}

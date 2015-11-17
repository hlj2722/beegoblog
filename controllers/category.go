package controllers

import (
	"beegoblog/models"
	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (this *CategoryController) Get() {
	this.Data["IsCategory"] = true
	this.TplNames = "category.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func (this *CategoryController) Load() {
	var categories []*(models.Category)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		categories, err = models.GetAllCategoriesRedis(true)
	default:
		categories, err = models.GetAllCategories(true)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Data["json"] = &categories
	this.ServeJson()
}

func (this *CategoryController) Post() {

	name := this.Input().Get("name")
	if len(name) == 0 {
		return
	}

	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		err = models.AddCategoryRedis(name)
	default:
		err = models.AddCategory(name)
	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/category", 302)

}

func (this *CategoryController) Delete() {

	id := this.Input().Get("id")
	if len(id) == 0 {
		return
	}

	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		err = models.DeleteCategoryRedis(id)
	default:
		err = models.DeleteCategory(id)
	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/category", 302)
}

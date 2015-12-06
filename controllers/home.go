package controllers

import (
	"github.com/hopehook/beegoblog/models"

	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	//this.Data 是一个用来存储输出数据的 map，可以赋值任意类型的值
	this.Data["IsHome"] = true
	//this.TplNames 就是需要渲染的模板，这里指定了 index.tpl，如果用户不设置该参数，那么默认会去到模板目录的 Controller/<方法名>.tpl 查找
	this.TplNames = "home.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	this.Data["category"] = this.Input().Get("category")
	this.Data["lable"] = this.Input().Get("lable")
}

func (this *HomeController) Load() {
	category := this.Input().Get("category")
	lable := this.Input().Get("lable")

	var topics []*(models.Topic)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		topics, err = models.GetAllTopicsRedis(category, lable, true)
	default:
		topics, err = models.GetAllTopics(category, lable, true)
	}

	if err != nil {
		beego.Error(err)
	}

	var categories []*(models.Category)
	switch beego.AppConfig.String("database") {
	case "redis":
		categories, err = models.GetAllCategoriesRedis(false)
	default:
		categories, err = models.GetAllCategories(false)
	}

	if err != nil {
		beego.Error(err)
	}

	data := &struct {
		Topics     []*models.Topic
		Categories []*models.Category
	}{
		Topics:     topics,
		Categories: categories,
	}
	this.Data["json"] = data
	this.ServeJson()
}

package controllers

import (
	"beegoblog/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["IsHome"] = true  //this.Data 是一个用来存储输出数据的 map，可以赋值任意类型的值
	this.TplNames = "home.html" //this.TplNames 就是需要渲染的模板，这里指定了 index.tpl，如果用户不设置该参数，那么默认会去到模板目录的 Controller/<方法名>.tpl 查找
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	this.Data["category"] = this.Input().Get("category")
	this.Data["lable"] = this.Input().Get("lable")
}

func (this *HomeController) Load() {
	category := this.Input().Get("category")
	lable := this.Input().Get("lable")

	//日志
	log := logs.NewLogger(10000)
	log.SetLogger("console", `{"level":1}`)
	log.Alert(category)
	log.Alert(lable)

	topics, err := models.GetAllTopics(category, lable, true)
	if err != nil {
		beego.Error(err)
	}

	//TODO 过滤掉 TopicCount<=0 的分类记录
	categories, err := models.GetAllCategories(false)
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

package controllers

import (
	"path"
	"strings"

	"github.com/astaxie/beego"

	"github.com/hopehook/beegoblog/models"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get() {

	this.Data["IsTopic"] = true
	this.TplNames = "topic.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func (this *TopicController) Load() {

	var topics []*(models.Topic)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		topics, err = models.GetAllTopicsRedis("", "", false)
	default:
		topics, err = models.GetAllTopics("", "", false)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Data["json"] = &topics
	this.ServeJson()
}

func (this *TopicController) Post() {

	var author string = ""
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	} else {
		ck, err := this.Ctx.Request.Cookie("uname")
		if err == nil {
			author = ck.Value
		}

	}

	// 解析表单
	tid := this.Input().Get("tid")
	title := this.Input().Get("title")
	content := this.Input().Get("content")
	category := this.Input().Get("category")
	lable := this.Input().Get("lable")

	// 获取附件
	_, fh, err := this.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}

	var attachment string
	if fh != nil {
		// 保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	if len(tid) == 0 {
		//新增文章的时候加入作者
		switch beego.AppConfig.String("database") {
		case "redis":
			err = models.AddTopicRedis(title, category, lable, content, attachment, author)
		default:
			err = models.AddTopic(title, category, lable, content, attachment, author)
		}

	} else {
		//如果其他人修改文章，作者不变
		switch beego.AppConfig.String("database") {
		case "redis":
			err = models.ModifyTopicRedis(tid, title, category, lable, content, attachment)
		default:
			err = models.ModifyTopic(tid, title, category, lable, content, attachment)
		}

	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic", 302)
}

func (this *TopicController) Add() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	this.TplNames = "topic_add.html"
	this.Data["IsLogin"] = true
}

func (this *TopicController) Delete() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		err = models.DeleteTopicRedis(this.Input().Get("tid"))
	default:
		err = models.DeleteTopic(this.Input().Get("tid"))
	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic", 302)
}

func (this *TopicController) Modify() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	this.Data["tid"] = this.Input().Get("tid")
	this.TplNames = "topic_modify.html"
	this.Data["IsLogin"] = true

}

func (this *TopicController) LoadModify() {

	tid := this.Input().Get("tid")

	var topic *(models.Topic)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		topic, err = models.GetTopicRedis(tid)
	default:
		topic, err = models.GetTopic(tid)
	}

	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	this.Data["json"] = &topic
	this.ServeJson()

}

func (this *TopicController) View() {
	tid := this.Input().Get("tid")
	this.Data["tid"] = tid

	this.TplNames = "topic_view.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

}

func (this *TopicController) LoadView() {
	tid := this.Input().Get("tid")

	var topic *(models.Topic)
	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		topic, err = models.GetTopicRedis(tid)
	default:
		topic, err = models.GetTopic(tid)
	}

	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	var replies []*(models.Reply)
	switch beego.AppConfig.String("database") {
	case "redis":
		replies, err = models.GetAllRepliesRedis(tid)
	default:
		replies, err = models.GetAllReplies(tid)
	}

	if err != nil {
		beego.Error(err)
		return
	}

	data := &struct {
		Topic   *models.Topic
		Lables  []string
		Replies []*models.Reply
	}{
		Topic:   topic,
		Lables:  strings.Split(topic.Lables, " "),
		Replies: replies,
	}

	this.Data["json"] = data
	this.ServeJson()
}

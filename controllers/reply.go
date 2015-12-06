package controllers

import (
	"github.com/hopehook/beegoblog/models"

	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

func (this *ReplyController) Add() {
	tid := this.Input().Get("tid")

	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		err = models.AddReplyRedis(tid,
			this.Input().Get("nickname"), this.Input().Get("content"))
	default:
		err = models.AddReply(tid,
			this.Input().Get("nickname"), this.Input().Get("content"))
	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic/view?tid="+tid, 302)
}

func (this *ReplyController) Delete() {

	tid := this.Input().Get("tid")
	beego.Alert(tid)
	beego.Alert(tid)

	if !checkAccount(this.Ctx) {
		return
	}

	var err error
	switch beego.AppConfig.String("database") {
	case "redis":
		err = models.DeleteReplyRedis(this.Input().Get("rid"))
	default:
		err = models.DeleteReply(this.Input().Get("rid"))
	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic/view?tid="+tid, 302)
}

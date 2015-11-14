package main

import (
	_ "beegoblog/routers"
	
	"os"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)


func main() {
	// 开启 ORM 调试模式
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)

	// 附件处理
	os.Mkdir("attachment", os.ModePerm)
	
	/*
	Go 语言的默认模板采用了 {{ 和 }} 作为左右标签，但是我们有时候在开发中可能界面是采用了 AngularJS 开发，
	他的模板也是这个标签，故而引起了冲突。在 beego 中你可以通过配置文件或者直接设置配置变量修改：
	*/
	beego.TemplateLeft = "<<<"
    beego.TemplateRight = ">>>"
	
	
	// 启动 beego
	beego.Run()


}

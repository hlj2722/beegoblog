package routers

import (
	"beegoblog/controllers"
	"beegoblog/models"
	"github.com/astaxie/beego"
)

func init() {
	// 注册数据库
	models.RegisterDB()

	/*
	  注册 beego 固定路由
	*/
	/*
	  如下所示的路由就是我们最常用的路由方式，一个固定的路由，一个控制器，然后根据用户请求方法不同请求控制器中对应的方法，典型的 RESTful 方式。
	*/
	/*
	  除了前缀两个/:controller/:method的匹配之外，剩下的 url beego会帮你自动化解析为参数，保存在 this.Ctx.Input.Params 当中：
	  object/blog/2013/09/12  调用 ObjectController 中的 Blog 方法，参数如下：map[0:2013 1:09 2:12]
	*/
	/*
		现在已经可以通过自动识别出来下面类似的所有url，都会把请求分发到 controller 的 simple 方法：

		/controller/simple
		/controller/simple.html
		/controller/simple.json
		/controller/simple.rss

		可以通过 this.Ctx.Input.Param(":ext") 获取后缀名。
	*/

	beego.Router("/", &controllers.HomeController{})
	beego.Router("/home", &controllers.HomeController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/reply", &controllers.ReplyController{})
	beego.Router("/login", &controllers.LoginController{})

	beego.Router("/load", &controllers.HomeController{}, "get:Load")
	beego.Router("/home/load", &controllers.HomeController{}, "get:Load")
	beego.Router("/category/load", &controllers.CategoryController{}, "get:Load")

	beego.Router("/category/delete", &controllers.CategoryController{}, "get:Delete")
	beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
	beego.Router("/reply/delete", &controllers.ReplyController{}, "get:Delete")
	beego.Router("/attachment/:all", &controllers.AttachController{})

	/*
		智能路由
	*/
	/*
		有了这个AutoRouter，便不需要像以前那样逐一注册了，访问/user/add 调用UserController的Add方法，访问/page/about调用PageController的About方法。
	*/

	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/topic/load", &controllers.TopicController{}, "get:Load")
	beego.Router("/user/load", &controllers.UserController{}, "get:Load")
	beego.Router("/topic/loadModify", &controllers.TopicController{}, "get:LoadModify")
	beego.Router("/topic/loadView", &controllers.TopicController{}, "get:LoadView")
	beego.Router("/user/loadModify", &controllers.UserController{}, "get:LoadModify")
	beego.Router("/user/loadView", &controllers.UserController{}, "get:LoadView")

}

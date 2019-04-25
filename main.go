package main

import (
	"github.com/bluesky1024/goMblog/web/middleware"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
)

func main() {
	//初始化各种基础组件
	initServ()

	app := iris.New()
	app.Use(recover.New())

	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html"). //设置html基础框架
		Reload(true)                  //设置开发过程中html能动态刷新
	app.RegisterView(tmpl)

	//访问静态资源时的路径
	//例： localhost:8080/public/css/site.css  =>  ./web/public/css/site.css
	app.StaticWeb("/public", "./web/public")

	//出现任何错误均显示error页
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	//全局中间件
	app.Use(middleware.CheckLogin)

	//默认页面
	app.Get("/", RequestMain)

	//注册路由
	mvc.Configure(app.Party("/user"), user)
	mvc.Configure(app.Party("/personal"), personal)
	mvc.Configure(app.Party("/mblog"), mblog)
	mvc.Configure(app.Party("/relation"), relation)
	mvc.Configure(app.Party("/feed"), feed)
	//	mvc.Configure(app.Party("/public"), public)
	//	mvc.Configure(app.Party("/search"), search)

	//启动
	app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

	defer func() {
		resourceRecycle()
	}()
}

func resourceRecycle() {
	//服务释放
	if relationSrv != nil {
		relationSrv.ReleaseSrv()
	}
}

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

	//出现任何错误均显示error页
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	//全局中间件
	app.Use(middleware.CheckLogin)
	//app.Use(cache.StaticCache(24 * time.Hour))

	//默认页面
	app.Get("/", RequestMain)

	//注册路由
	mvc.Configure(app.Party("/public"), staticFile)

	mvc.Configure(app.Party("/user"), user)
	//mvc.Configure(app.Party("/personal"), personal)
	//mvc.Configure(app.Party("/mblog"), mblog)
	//mvc.Configure(app.Party("/relation"), relation)
	//mvc.Configure(app.Party("/feed"), feed)
	mvc.Configure(app.Party("/chat"), chat)
	//	mvc.Configure(app.Party("/public"), public)
	//	mvc.Configure(app.Party("/search"), search)

	//启动
	err := app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		resourceRecycle()
	}()
}

func resourceRecycle() {
	//服务释放
	if relationSrv != nil {
		relationSrv.ReleaseSrv()
	}
	if mblogSrv != nil {
		mblogSrv.ReleaseSrv()
	}
}

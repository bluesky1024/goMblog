package main

import (
	"github.com/bluesky1024/goMblog/web/controllers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func RequestMain(ctx iris.Context) {
	ctx.Redirect("/personal/profile")
}

func user(app *mvc.Application) {
	//	//设置中间件
	//	app.Router.Use(middleware.BasicAuth)

	//注册服务
	app.Register(userSrv)
	app.Register(SessManager.Start)

	//绑定控制器
	app.Handle(new(controllers.UserController))
}

func personal(app *mvc.Application) {
	//注册服务
	app.Register(userSrv)
	app.Register(mblogSrv)
	app.Register(relationSrv)
	app.Register(SessManager.Start)

	app.Handle(new(controllers.PersonalController))
}

func mblog(app *mvc.Application) {
	//注册服务
	app.Register(userSrv)
	app.Register(mblogSrv)
	app.Register(SessManager.Start)

	app.Handle(new(controllers.MblogController))
}

func relation(app *mvc.Application) {
	//注册服务
	app.Register(userSrv)
	app.Register(relationSrv)
	app.Register(SessManager.Start)

	app.Handle(new(controllers.RelationController))
}

func feed(app *mvc.Application) {
	app.Register(feedSrv)
	app.Register(userSrv)
	app.Register(mblogSrv)
	app.Register(relationSrv)
	app.Register(SessManager.Start)

	app.Handle(new(controllers.FeedController))
}

//func public(app *mvc.Application) {
//	app.Handle(new(controllers.PublicController))
//}

//func search(app *mvc.Application) {
//	app.Handle(new(controllers.SearchController))
//}

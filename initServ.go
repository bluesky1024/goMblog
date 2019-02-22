package main

import (
	relationServ "github.com/bluesky1024/goMblog/services/relation"
	userServ "github.com/bluesky1024/goMblog/services/user"
	mblogServ "github.com/bluesky1024/goMblog/services/mblog"

	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/web/middleware"
	"github.com/kataras/iris/sessions"
	"time"
)

var(
	SessManager *sessions.Sessions

	userSrv userServ.UserServicer
	mblogSrv  mblogServ.MblogServicer
	relationSrv relationServ.RelationServicer
)

func initServ() {
	initBasicServ()
	initGlobalSession()
	initIdGen()

	initMiddleware()
}

func initBasicServ(){
	userSrv = userServ.NewUserServicer()
	mblogSrv = mblogServ.NewMblogServicer()
	relationSrv = relationServ.NewRelationServicer()
}

func initIdGen() {
	idGen.InitMidPool(10)
	idGen.InitUidPool(3)
}

func initGlobalSession(){
	//设置全局session
	SessManager = sessions.New(sessions.Config{
		Cookie:  "my_session",
		Expires: 24 * time.Hour,
	})
}

func initMiddleware(){
	middleware.RegisterGlobalSession(SessManager)
	middleware.RegisterUserServer(userSrv)
}
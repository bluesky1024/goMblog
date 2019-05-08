package main

import (
	feedServ "github.com/bluesky1024/goMblog/services/feed"
	mblogServ "github.com/bluesky1024/goMblog/services/mblog"
	relationServ "github.com/bluesky1024/goMblog/services/relation"
	userServ "github.com/bluesky1024/goMblog/services/user"

	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/web/middleware"
	"github.com/kataras/iris/sessions"
	"time"
)

var (
	SessManager *sessions.Sessions

	userSrv     userServ.UserServicer
	mblogSrv    mblogServ.MblogServicer
	relationSrv relationServ.RelationServicer
	feedSrv     feedServ.FeedServicer
)

func initServ() {
	initBasicServ()
	initGlobalSession()
	//initIdGen()

	initMiddleware()
}

func initBasicServ() {
	var err error

	userSrv, err = userServ.NewUserServicer()
	if err != nil {
		panic(err.Error())
	}

	mblogSrv, err = mblogServ.NewMblogServicer()
	if err != nil {
		panic(err.Error())
	}

	relationSrv, err = relationServ.NewRelationServicer()
	if err != nil {
		panic(err.Error())
	}

	feedSrv, err = feedServ.NewFeedServicer()
	if err != nil {
		panic(err.Error())
	}
}

func initIdGen() {
	idGen.InitMidPool(10)
	idGen.InitUidPool(3)
}

func initGlobalSession() {
	//设置全局session
	SessManager = sessions.New(sessions.Config{
		Cookie:  "my_session",
		Expires: 24 * time.Hour,
	})
}

func initMiddleware() {
	middleware.RegisterGlobalSession(SessManager)
	middleware.RegisterUserServer(userSrv)
}

package main

import(
	relationServ "github.com/bluesky1024/goMblog/services/relation"
	userServ "github.com/bluesky1024/goMblog/services/user"
	mblogServ "github.com/bluesky1024/goMblog/services/mblog"
)

var(
	userSrv userServ.UserServicer
	mblogSrv  mblogServ.MblogServicer
	relationSrv relationServ.RelationServicer
)

func initServ(){
	initBasicServ()
}

func initBasicServ(){
	userSrv = userServ.NewUserServicer()
	mblogSrv = mblogServ.NewMblogServicer()
	relationSrv = relationServ.NewRelationServicer()
}
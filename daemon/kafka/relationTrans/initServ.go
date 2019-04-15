package main

import (
	mblogServ "github.com/bluesky1024/goMblog/services/mblog"
	relationServ "github.com/bluesky1024/goMblog/services/relation"
	userServ "github.com/bluesky1024/goMblog/services/user"
)

var (
	userSrv     userServ.UserServicer
	mblogSrv    mblogServ.MblogServicer
	relationSrv relationServ.RelationServicer
)

func initServ() {
	initBasicServ()
}

func initBasicServ() {
	userSrv = userServ.NewUserServicer()
	mblogSrv = mblogServ.NewMblogServicer()
	relationSrv = relationServ.NewRelationServicer()
}

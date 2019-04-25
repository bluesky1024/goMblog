package main

import (
	relationServ "github.com/bluesky1024/goMblog/services/relation"
)

var (
	relationSrv relationServ.RelationServicer
)

func initServ() {
	initBasicServ()
}

func initBasicServ() {
	var err error
	relationSrv, err = relationServ.NewRelationServicer()
	if err != nil {
		panic(err.Error())
	}
}

func resourceRecycle() {
	//服务释放
	if relationSrv != nil {
		relationSrv.ReleaseSrv()
	}
}

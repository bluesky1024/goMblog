package main

import (
	feedServ "github.com/bluesky1024/goMblog/services/feed"
)

var (
	feedSrv feedServ.FeedServicer
)

func initServ() {
	initBasicServ()
}

func initBasicServ() {
	var err error
	feedSrv, err = feedServ.NewFeedServicer()
	if err != nil {
		panic(err.Error())
	}
}

func resourceRecycle() {
	//服务释放
	//if userSrv != nil {
	//
	//}
}

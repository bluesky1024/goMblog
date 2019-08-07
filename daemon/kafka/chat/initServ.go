package main

import (
	"github.com/bluesky1024/goMblog/services/chat"
)

var (
	chatSrv chatService.ChatServicer
)

func initServ() {
	initBasicServ()
}

func initBasicServ() {
	var err error
	chatSrv, err = chatService.NewChatServicer()
	if err != nil {
		panic(err.Error())
	}
}

func resourceRecycle() {
	//服务释放
	if chatSrv != nil {
		chatSrv.ReleaseSrv()
	}
}

package main

import (
	userServ "github.com/bluesky1024/goMblog/services/user"
)

var (
	userSrv userServ.UserServicer
)

func initServ() {
	initBasicServ()
}

func initBasicServ() {
	var err error
	userSrv, err = userServ.NewUserServicer()
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

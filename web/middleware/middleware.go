package middleware

import (
	"github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris/sessions"
)

var (
	UserServer  userService.UserServicer
	SessManager *sessions.Sessions
	logType     = "middleware"
)

func RegisterGlobalSession(s *sessions.Sessions) {
	SessManager = s
}

func RegisterUserServer(u userService.UserServicer){
	UserServer = u
}
package controllers

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

//统一返回参数
type ResParams struct {
	Code int
	Msg  string
	Data interface{}
}

//登录相关
const UserIDKey = "UserID"

func GetCurrentUserID(s *sessions.Session) int64 {
	userID := s.GetInt64Default(UserIDKey, 0)
	return userID
}

func IsLoggedIn(s *sessions.Session) bool {
	return GetCurrentUserID(s) > 0
}

func Logout(s *sessions.Session) {
	s.Destroy()
}

func SetSessionUserId(s *sessions.Session,uid int64){
	s.Set(UserIDKey, uid)
}

func GenViewResponse(ctx iris.Context,viewPath string,data iris.Map) mvc.Result{
	//data中统一封装当前登录信息
	curUid := ctx.Values().Get("CurUid").(int64)
	curUserInfo := ctx.Values().Get("CurUserInfo").(dm.User)
	curNickName := curUserInfo.NickName

	if data == nil{
		data = make(iris.Map)
	}
	data["CurUid"] = curUid
	data["CurNickName"] = curNickName

	return mvc.View{
		Name: viewPath,
		Data: data,
	}
}

func GenErrorView(ctx iris.Context,data iris.Map) mvc.Result{
	return GenViewResponse(ctx,"shared/error.html",data)
}
package middleware

import (
	"strings"

	"github.com/bluesky1024/goMblog/config"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/bluesky1024/goMblog/web/controllers"
	"github.com/kataras/iris"
)

func CheckLogin(ctx iris.Context) {
	curSession := SessManager.Start(ctx)
	userID := curSession.GetInt64Default(controllers.UserIDKey, 0)

	if isNeedLogin(ctx) {
		if userID <= 0 {
			notLoginReturn(ctx)
			return
		}
	}

	var CurUserInfo = dm.User{}
	if userID > 0 {
		CurUserInfo, _ = UserServer.GetByUid(userID)
	}

	ctx.Values().Set("CurUid", userID)
	ctx.Values().Set("CurUserInfo", CurUserInfo)
	ctx.Next()
}

func isNeedLogin(ctx iris.Context) bool {
	curMethod := ctx.Method()
	curPath := ctx.Path()

	checkMap := conf.InitConfig("visitCtrl")
	if _, ok := checkMap[curPath]; ok {
		methodSlice := strings.Split(checkMap[curPath], "/")
		for _, v := range methodSlice {
			if v == curMethod {
				return true
			}
		}
	}
	return false
}

func notLoginReturn(ctx iris.Context) {
	ctx.StopExecution()
	//Ajax
	if ctx.IsAjax() {
		_, err := ctx.JSON(controllers.ResParams{
			Code: 1001,
			Msg:  "please login first",
			Data: nil,
		})
		if err != nil {
			logger.Err(logType, err.Error())
		}
		return
	}

	//Page
	ctx.Redirect("/user/login")
}

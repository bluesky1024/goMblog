// file: controllers/user_controller.go

package controllers

import (
	mblogSrv "github.com/bluesky1024/goMblog/services/mblog"
	userSrv "github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

// POST  			/mblog/send
// GET  			/mblog/[mid-encode]
// POST 			/mblog/delete
// POST 			/mblog/like
// POST 			/mblog/transmit
// POST 			/mblog/comment
// POST 			/mblog/comment/list
type MblogController struct {
	Ctx iris.Context

	UserSrv  userSrv.UserServicer
	MblogSrv mblogSrv.MblogServicer

	Session *sessions.Session
}

var (
	mblogSendStaticView = mvc.View{
		Name: "mblog/send.html",
		Data: iris.Map{"Title": "微博发布"},
	}
)

// GetRegister handles GET: http://localhost:8080/mblog/send.
func (c *MblogController) GetSend() mvc.Result {
	data := iris.Map{
		"Title": "微博发布",
	}
	return GenViewResponse(c.Ctx, "mblog/send", data)
}

// PostSend handles Post: http://localhost:8080/mblog/send.
func (c *MblogController) PostSend() interface{} {
	if !IsLoggedIn(c.Session) {
		return ResParams{
			Code: 1001,
			Msg:  "please login first",
			Data: nil,
		}
	}
	var (
		content     = c.Ctx.FormValue("content")
		readAble, _ = c.Ctx.PostValueInt("readAble")
	)
	curUid := GetCurrentUserID(c.Session)

	mblog, err := c.MblogSrv.Create(curUid, content, int8(readAble), 0, 0)
	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	return ResParams{
		Code: 1000,
		Msg:  "发布成功",
		Data: mblog,
	}
}

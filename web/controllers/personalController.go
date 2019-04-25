// file: controllers/user_controller.go

package controllers

import (
	mblogSrv "github.com/bluesky1024/goMblog/services/mblog"
	relationSrv "github.com/bluesky1024/goMblog/services/relation"
	userSrv "github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"strconv"
)

// GET  			/personal/profile/[uid]
// GET  			/personal/baseinfo
// POST				/personal/baseinfo/edit
type PersonalController struct {
	Ctx iris.Context

	UserSrv     userSrv.UserServicer
	MblogSrv    mblogSrv.MblogServicer
	RelationSrv relationSrv.RelationServicer

	Session *sessions.Session
}

// GetProfle handles GET: http://localhost:8080/personal/profile.
func (c *PersonalController) GetProfile() {
	//判断是否登录
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	if CurUid == 0 {
		c.Ctx.Redirect("/user/login")
	}

	c.Ctx.Redirect("/personal/profile/" + strconv.FormatInt(CurUid, 10))
}

// GetProfle handles GET: http://localhost:8080/personal/profile/[uid].
func (c *PersonalController) GetProfileBy(uid int64) interface{} {
	//获取基础数据
	mblogs := c.MblogSrv.GetByUid(uid, 1, 20)
	userInfo, _ := c.UserSrv.GetByUid(uid)

	//获取uid与当前登录账户的关系
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	relationStatus := c.RelationSrv.CheckFollow(curUid, uid)

	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	var title string
	if CurUid == uid {
		title = "personal profile"
	} else {
		title = userInfo.NickName + "'s profile"
	}

	data := iris.Map{
		"Title":          title,
		"Mblogs":         mblogs,
		"UserInfo":       userInfo,
		"RelationStatus": relationStatus,
	}
	return GenViewResponse(c.Ctx, "personal/profile.html", data)
}

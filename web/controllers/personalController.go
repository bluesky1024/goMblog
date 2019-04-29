// file: controllers/user_controller.go

package controllers

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
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
	pageStr := c.Ctx.FormValue("page")
	page, err := strconv.Atoi(pageStr)
	pageSize := 5
	if err != nil || page < 0 {
		page = 1
	}

	userInfo, _ := c.UserSrv.GetByUid(uid)

	//获取uid与当前登录账户的关系
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	relationStatus := c.RelationSrv.CheckFollow(curUid, uid)

	//获取基础数据
	//关系还没做好，目前暂定为朋友关系来获取微博
	readAble := []int8{dm.MblogReadAblePublic, dm.MblogReadAbleFriend}
	mblogs, cnt := c.MblogSrv.GetNormalByUid(uid, page, pageSize, readAble, 0, 0)

	var title string
	if curUid == uid {
		title = "personal profile"
	} else {
		title = userInfo.NickName + "'s profile"
	}

	//分页参数
	uidStr := strconv.FormatInt(uid, 10)
	pageParam := GenPageView(int64(page), int64(pageSize), cnt, "/personal/profile/"+uidStr)

	data := iris.Map{
		"Title":          title,
		"Mblogs":         mblogs,
		"UserInfo":       userInfo,
		"RelationStatus": relationStatus,
		"PageParam":      pageParam,
	}
	return GenViewResponse(c.Ctx, "personal/profile.html", data)
}

/*测试JSONP跨域请求*/
func (c *PersonalController) GetTest() interface{} {
	type tempStruct struct {
		Data string
	}
	return tempStruct{
		Data: "https://www.baidu.com/sugrec?pre=1&p=3&ie=utf-8&prod=pc&from=wise_web&wd=lol&req=2&bs=lol&pbs=lol&pwd=lol",
	}
}

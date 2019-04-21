package controllers

import (
	feedSrv "github.com/bluesky1024/goMblog/services/feed"
	mblogSrv "github.com/bluesky1024/goMblog/services/mblog"
	relationSrv "github.com/bluesky1024/goMblog/services/relation"
	userSrv "github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"strconv"
)

type FeedController struct {
	Ctx iris.Context

	FeedSrv     feedSrv.FeedServicer
	UserSrv     userSrv.UserServicer
	MblogSrv    mblogSrv.MblogServicer
	RelationSrv relationSrv.RelationServicer

	Session *sessions.Session
}

// GetProfle handles GET: http://localhost:8080/feed/.
func (c *FeedController) Get() interface{} {
	//判断是否登录
	CurUid := c.Ctx.Values().Get("CurUid").(int64)

	if CurUid == 0 {
		c.Ctx.Redirect("/user/login")
	}

	//用户基础数据
	userInfo, _ := c.UserSrv.GetByUid(CurUid)

	//分组数据
	groupsInfo, _ := c.RelationSrv.GetGroupsByUid(CurUid)

	data := iris.Map{
		"Title":      "my feed",
		"GroupsInfo": groupsInfo,
		"UserInfo":   userInfo,
	}
	return GenViewResponse(c.Ctx, "feed/feed.html", data)
}

// PostMoreFeed handle POST: http://localhost:8080/feed/more/.
func (c *FeedController) PostMore() interface{} {
	//判断是否登录
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	CurUid = 2317487850917888
	if CurUid == 0 {
		return ResParams{
			Code: 1001,
			Msg:  "please login first",
		}
	}
	var (
		groupIdStr = c.Ctx.FormValue("groupId")
		lastMidStr = c.Ctx.FormValue("lastMid")
		//pageSizeStr = c.Ctx.FormValue("pageSize")
	)

	pageSizeStr := "20"

	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	lastMid, _ := strconv.ParseInt(lastMidStr, 10, 64)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	//获取feed数据
	mids, err := c.FeedSrv.GetFeedMoreByMid(CurUid, groupId, lastMid, pageSize)
	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  "拉取feed失败",
		}
	}

	//根据mid获取具体微博信息（此处时间上有问题，需要重新排序）
	mblogsInfo := c.MblogSrv.GetMultiByMids(mids)

	//根据uid获取用户信息
	uids := make([]int64, len(mblogsInfo))
	ind := 0
	for _, mblog := range mblogsInfo {
		uids[ind] = mblog.Uid
		ind++
	}
	usersInfo, err := c.UserSrv.GetMultiByUids(uids)

	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  "拉取feed失败",
		}
	}

	data := iris.Map{
		"MblogsInfo": mblogsInfo,
		"UsersInfo":  usersInfo,
	}

	return ResParams{
		Code: 1000,
		Msg:  "拉取成功",
		Data: data,
	}
}

// PostNewerFeed handle POST: http://localhost:8080/feed/newer/.
func (c *FeedController) PostNewer() interface{} {
	//判断是否登录
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	CurUid = 2317487850917888
	if CurUid == 0 {
		return ResParams{
			Code: 1001,
			Msg:  "please login first",
		}
	}
	var (
		groupIdStr  = c.Ctx.FormValue("groupId")
		firstMidStr = c.Ctx.FormValue("firstMid")
		//pageSizeStr = c.Ctx.FormValue("pageSize")
	)

	pageSizeStr := "20"

	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	firstMid, _ := strconv.ParseInt(firstMidStr, 10, 64)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	//获取feed数据
	mids, err := c.FeedSrv.GetFeedNewerByMid(CurUid, groupId, firstMid, pageSize)
	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  "拉取feed失败",
		}
	}

	//根据mid获取具体微博信息（此处时间上有问题，需要重新排序）
	mblogsInfo := c.MblogSrv.GetMultiByMids(mids)

	//根据uid获取用户信息
	uids := make([]int64, len(mblogsInfo))
	ind := 0
	for _, mblog := range mblogsInfo {
		uids[ind] = mblog.Uid
		ind++
	}
	usersInfo, err := c.UserSrv.GetMultiByUids(uids)

	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  "拉取feed失败",
		}
	}

	data := iris.Map{
		"MblogsInfo": mblogsInfo,
		"UsersInfo":  usersInfo,
	}

	return ResParams{
		Code: 1000,
		Msg:  "拉取成功",
		Data: data,
	}
}

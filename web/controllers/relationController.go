package controllers

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	relationServ "github.com/bluesky1024/goMblog/services/relation"
	userServ "github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
)

// GET  			/relation/follows/[uid]
// GET  			/relation/fans/[uid]
// POST				/relation/follow
// POST 			/relation/unfollow
type RelationController struct {
	Ctx iris.Context

	RelationServ relationServ.RelationServicer
	UserServ     userServ.UserServicer

	Session *sessions.Session
}

func (c *RelationController) GetFollows() mvc.Result {
	data := iris.Map{
		"Title": "我的关注",
	}

	curUid := c.Ctx.Values().Get("CurUid").(int64)

	follows, cnt := c.RelationServ.GetFollowsByUid(curUid, 1, 20)

	followUids := make([]int64, 20)
	for ind, follow := range follows {
		followUids[ind] = follow.FollowUid
	}
	followUserInfo, err := c.UserServ.GetMultiByUids(followUids)
	if err != nil {
		GenErrorView(c.Ctx, data)
	}

	data["Follows"] = follows
	data["Cnt"] = cnt
	data["UserInfos"] = followUserInfo
	return GenViewResponse(c.Ctx, "relation/follows.html", data)
}

func (c *RelationController) GetFans() mvc.Result {
	data := iris.Map{
		"Title": "用户注册",
	}
	return GenViewResponse(c.Ctx, "user/register.html", data)
}

// GetFollows handles GET: http://localhost:8080/relation/follows.
func (c *RelationController) GetFollowsBy(uid int64) mvc.Result {
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	if uid == curUid {
		c.Ctx.Redirect("relation/follows")
	}
	data := iris.Map{
		"Title": "用户注册",
	}
	return GenViewResponse(c.Ctx, "user/register.html", data)
}

// GetFans handles GET: http://localhost:8080/relation/follows.
func (c *RelationController) GetFansBy(uid int64) mvc.Result {
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	if uid == curUid {
		c.Ctx.Redirect("relation/fans")
	}
	data := iris.Map{
		"Title": "用户注册",
	}
	return GenViewResponse(c.Ctx, "user/register.html", data)
}

// PostFollow handles POST: http://localhost:8080/relation/follow.
func (c *RelationController) PostFollow() interface{} {
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	if curUid == 0 {
		return ResParams{
			Code: 1001,
			Msg:  "please login first",
		}
	}
	uidStr := c.Ctx.FormValue("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	if _, found := c.UserServ.GetByUid(uid); !found {
		return ResParams{
			Code: 1001,
			Msg:  "invalid follow object",
		}
	}
	res := c.RelationServ.Follow(curUid, uid)
	if !res {
		return ResParams{
			Code: 1001,
			Msg:  "follow failed",
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "follow success",
	}
}

// PostUnfollow handles POST: http://localhost:8080/relation/unfollow.
func (c *RelationController) PostUnfollow() interface{} {
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	if curUid == 0 {
		return ResParams{
			Code: 1001,
			Msg:  "please login first",
		}
	}
	uidStr := c.Ctx.FormValue("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	if _, found := c.UserServ.GetByUid(uid); !found {
		return ResParams{
			Code: 1001,
			Msg:  "invalid follow object",
		}
	}
	res := c.RelationServ.UnFollow(curUid, uid)
	if !res {
		return ResParams{
			Code: 1001,
			Msg:  "follow failed",
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "follow success",
	}
}

func (c *RelationController) GetGroups() interface{} {
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	groups, _ := c.RelationServ.GetGroupsByUid(curUid)
	data := iris.Map{
		"Title":  "分组管理",
		"Groups": groups,
	}
	return GenViewResponse(c.Ctx, "relation/groups.html", data)
}

func (c *RelationController) PostAddGroup() interface{} {
	groupName := c.Ctx.FormValue("groupName")
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	res := c.RelationServ.AddGroup(curUid, groupName)
	if !res {
		return ResParams{
			Code: 1001,
			Msg:  "add group failed",
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "add group success",
	}
}

func (c *RelationController) PostDelGroup() interface{} {
	groupIdStr := c.Ctx.FormValue("groupId")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	res := c.RelationServ.DelGroup(curUid, groupId)
	if !res {
		return ResParams{
			Code: 1001,
			Msg:  "del group failed",
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "del group success",
	}
}

func (c *RelationController) PostUpdateGroup() interface{} {
	var (
		groupIdStr = c.Ctx.FormValue("groupId")
		groupName  = c.Ctx.FormValue("groupName")
	)
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	curUid := c.Ctx.Values().Get("CurUid").(int64)
	group := dm.FollowGroup{
		Id:        groupId,
		Uid:       curUid,
		GroupName: groupName,
		Status:    dm.GroupStatusNormal,
	}
	res := c.RelationServ.UpdateGroup(group)

	if !res {
		return ResParams{
			Code: 1001,
			Msg:  "update group failed",
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "update group success",
	}

}

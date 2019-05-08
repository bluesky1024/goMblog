package controllers

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"math"
)

//统一返回参数
type ResParams struct {
	Code int
	Msg  string
	Data interface{}
}

//统一分页样式参数
type PageParams struct {
	CurPage  int64
	Cnt      int64
	BaseUrl  string
	PageList []int64
	First    bool
	Last     bool
	Forward  int64
	Back     int64
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

func SetSessionUserId(s *sessions.Session, uid int64) {
	s.Set(UserIDKey, uid)
}

func GenViewResponse(ctx iris.Context, viewPath string, data iris.Map) mvc.Result {
	//data中统一封装当前登录信息
	curUid := ctx.Values().Get("CurUid").(int64)
	curUserInfo := ctx.Values().Get("CurUserInfo").(dm.User)
	curNickName := curUserInfo.NickName

	if data == nil {
		data = make(iris.Map)
	}
	data["CurUid"] = curUid
	data["CurNickName"] = curNickName
	data["CurUserInfo"] = curUserInfo

	return mvc.View{
		Name: viewPath,
		Data: data,
	}
}

func GenErrorView(ctx iris.Context, data iris.Map) mvc.Result {
	if data == nil {
		data = make(iris.Map)
	}
	data["Title"] = "error page"
	return GenViewResponse(ctx, "shared/error.html", data)
}

func GenPageView(curPage int64, pageSize int64, cnt int64, baseUrl string) (pageParam PageParams) {
	var startListInd int64
	var endListInd int64
	pageCnt := int64(math.Ceil(float64(cnt) / float64(pageSize)))
	pageParam.Cnt = pageCnt
	pageParam.BaseUrl = baseUrl
	pageParam.CurPage = curPage
	switch {
	case curPage > 4:
		startListInd = curPage - 2
		pageParam.First = true
		pageParam.Back = curPage - 3
		break
	case curPage == 4:
		startListInd = curPage - 2
		pageParam.First = true
		pageParam.Back = 0
		break
	default:
		startListInd = 1
		pageParam.First = false
		pageParam.Back = 0
		break
	}
	switch {
	case curPage+2 < pageCnt-1:
		endListInd = curPage + 2
		pageParam.Last = true
		pageParam.Forward = curPage + 3
		break
	case curPage+2 == pageCnt-1:
		endListInd = curPage + 2
		pageParam.Last = true
		pageParam.Forward = 0
		break
	default:
		endListInd = pageCnt
		pageParam.Last = false
		pageParam.Forward = 0
		break
	}
	pageParam.PageList = make([]int64, endListInd-startListInd+1)
	for ind, _ := range pageParam.PageList {
		pageParam.PageList[ind] = startListInd + int64(ind)
	}
	return pageParam
}

// file: controllers/user_controller.go

package controllers

import (
	userServ "github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

// GET  			/user/register
// POST 			/user/register
// GET 				/user/login
// POST 			/user/login
// GET 				/user/me
// All HTTP Methods /user/logout
type UserController struct {
	Ctx iris.Context

	UserServ userServ.UserServicer

	Session *sessions.Session
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	data := iris.Map{
		"Title": "用户注册",
	}
	return GenViewResponse(c.Ctx, "user/register.html", data)
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() interface{} {
	var (
		nickName  = c.Ctx.FormValue("nickName")
		password  = c.Ctx.FormValue("password")
		telephone = c.Ctx.FormValue("telephone")
		email     = c.Ctx.FormValue("email")
	)

	user, err := c.UserServ.Create(nickName, password, telephone, email)

	if err != nil {
		return ResParams{
			Code: 10001,
			Msg:  err.Error(),
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "注册成功",
		Data: user,
	}
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	data := iris.Map{
		"Title": "用户登录",
	}
	return GenViewResponse(c.Ctx, "user/login.html", data)
}

// PostLogin handles Post: http://localhost:8080/user/login.
func (c *UserController) PostLogin() interface{} {
	if IsLoggedIn(c.Session) {
		Logout(c.Session)
	}
	var (
		nickName = c.Ctx.FormValue("nickName")
		password = c.Ctx.FormValue("password")
	)

	user, found := c.UserServ.GetByNicknameAndPassword(nickName, password)
	if !found {
		return ResParams{
			Code: 1001,
			Msg:  "登录验证失败",
			Data: nil,
		}
	}

	//登录成功，种session
	SetSessionUserId(c.Session, user.Uid)

	return ResParams{
		Code: 1000,
		Msg:  "登录成功",
		Data: user,
	}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserController) AnyLogout() {
	if IsLoggedIn(c.Session) {
		Logout(c.Session)
	}

	c.Ctx.Redirect("/user/login")
}

//测试主页
type Test struct {
	Demo1 int
	Demo2 int
}

func (c *UserController) GetMe() interface{} {
	test := Test{
		Demo1: 123,
		Demo2: 234,
	}
	return test
}

//func (c *UserController) GetMe() interface{} {
//	return GenViewResponse(c.Ctx,"user/me.html",nil)
//}

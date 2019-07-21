package controllers

import (
	"fmt"
	userSrv "github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/websocket"
	"time"
)

// GET				/chat/get
type ChatController struct {
	Ctx iris.Context

	//长连接用的websocket
	WebsocketConn websocket.Connection

	//基础服务
	UserSrv userSrv.UserServicer

	Session *sessions.Session
}

// GetProfle handles GET: http://localhost:8080/chat/videoview.
func (c *ChatController) GetVideoviewBy(mid int64) interface{} {
	//判断是否登录
	CurUid := c.Ctx.Values().Get("CurUid").(int64)

	//if CurUid == 0 {
	//	c.Ctx.Redirect("/user/login")
	//}

	//用户基础数据
	userInfo, _ := c.UserSrv.GetByUid(CurUid)
	data := iris.Map{
		"Title":    "chat room",
		"UserInfo": userInfo,
		"RoomNO":   mid,
	}
	return GenViewResponse(c.Ctx, "chat/videoView.html", data)
}

func (c *ChatController) OnBarrageWebsocketConnect(roomName interface{}) {
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	fmt.Println(CurUid, "on connect", roomName)

	c.WebsocketConn.Join(roomName.(string))

	//启动定时任务，定时从这个roomName的kafka队列中拉数据发送到该连接的前端
	go func() {
		for {
			fmt.Println(c.WebsocketConn.ID())
			if !c.WebsocketConn.Server().IsConnected(c.WebsocketConn.ID()) {
				fmt.Println("server disconnect")
				//break
			}

			fmt.Println("for exist")
			fmt.Println(c.WebsocketConn.Emit("clientNewMsg", "roomName1"))
			time.Sleep(5 * time.Second)
			fmt.Println(c.WebsocketConn.Emit("clientNewMsg", "roomName1"))
			time.Sleep(5 * time.Second)
			fmt.Println(c.WebsocketConn.Emit("clientNewMsg", "roomName1"))
			time.Sleep(5 * time.Second)
		}

		fmt.Println("start_to_send_public")
		c.WebsocketConn.To(websocket.All).Emit("clientNewMsg", "this is a public message for all")
	}()
}
func (c *ChatController) OnBarrageWebsocketDisconnect(roomName string) {
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	fmt.Println(CurUid, "disconnect", roomName)
}
func (c *ChatController) OnBarrageWebsocketGetNewMessage(message string) {
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	fmt.Println(CurUid, "getNewMessage", message)
}

func (c *ChatController) GetBarrageWebsocketBy(mid int64) {
	c.WebsocketConn.OnLeave(c.OnBarrageWebsocketDisconnect)

	//mvc模式没找到OnConnect钩子函数，只能通过自定义事件，建立连接之后前端再次发送connect来响应
	c.WebsocketConn.On("connected", c.OnBarrageWebsocketConnect)
	c.WebsocketConn.On("newMsg", c.OnBarrageWebsocketGetNewMessage)

	//call it after all event callbacks registration.
	c.WebsocketConn.Wait()
}
package controllers

import (
	"encoding/json"
	"fmt"
	chatSrv "github.com/bluesky1024/goMblog/services/chat"
	userSrv "github.com/bluesky1024/goMblog/services/user"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/websocket"
	"math/rand"
	"strconv"
	"time"
)

// GET				/chat/get
type ChatController struct {
	Ctx iris.Context

	//长连接用的websocket
	WebsocketConn websocket.Connection

	//基础服务
	UserSrv userSrv.UserServicer
	ChatSrv chatSrv.ChatServicer

	Session *sessions.Session
}

func (c *ChatController) GetNoRoom() interface{} {
	data := iris.Map{
		"Title": "chat error",
	}
	return GenViewResponse(c.Ctx, "chat/noRoom.html", data)
}

// GetProfle handles GET: http://localhost:8080/chat/videoview.
func (c *ChatController) GetVideoviewBy(roomId int64) interface{} {
	//判断房间是否存在
	roomConfig, err := c.ChatSrv.GetRoomConfigByRoomId(roomId)
	if err != nil {
		c.Ctx.Redirect("/chat/no/room")
		return nil
	}

	//判断是否登录
	CurUid := c.Ctx.Values().Get("CurUid").(int64)

	//用户基础数据
	userInfo, _ := c.UserSrv.GetByUid(CurUid)
	data := iris.Map{
		"Title":    "chat room-" + roomConfig.RoomName,
		"UserInfo": userInfo,
		"RoomNO":   strconv.FormatInt(roomId, 10),
	}
	return GenViewResponse(c.Ctx, "chat/videoView.html", data)

}

type BarrageInfoView struct {
	UserName  string
	VideoTime int64
	Message   string
}

func (c *ChatController) OnBarrageWebsocketConnect(roomIdStr string) {
	roomId, _ := strconv.ParseInt(roomIdStr, 10, 64)
	//判断房间是否存在
	_, err := c.ChatSrv.GetRoomConfigByRoomId(roomId)
	if err != nil {
		disConnectErr := c.WebsocketConn.Disconnect()
		if disConnectErr != nil {
			fmt.Println(disConnectErr)
		}
		return
	}
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	//后续需要根据curUid来进行hash，如果未登录，需要分配一个随机的int64，此处允许重复
	if CurUid == 0 {
		rand.Seed(time.Now().UnixNano())
		CurUid = rand.Int63n(10000000-1) + 1
	}

	//启动定时任务，定时从这个roomId的redis池中拉数据发送到该连接的前端
	go func(roomId int64, visitorUid int64, server *websocket.Server) {
		for {
			if !server.IsConnected(c.WebsocketConn.ID()) {
				fmt.Println("server disconnect")
				break
			}

			barrageInfos, err := c.ChatSrv.GetBarrageByRoomId(visitorUid, roomId)
			if err != nil {
				//处理方式待补充
				fmt.Println(err)
			}

			//获取弹幕作者名
			uids := make([]int64, len(barrageInfos))
			for i, barrageInfo := range barrageInfos {
				uids[i] = barrageInfo.Uid
			}
			userInfos, err := c.UserSrv.GetMultiByUids(uids)
			if err != nil {
				//处理方式待补充
				fmt.Println(err)
			}

			//拼凑输出内容
			viewInfos := make([]BarrageInfoView, len(barrageInfos))
			for i, barrageInfo := range barrageInfos {
				viewInfos[i] = BarrageInfoView{
					UserName:  userInfos[barrageInfo.Uid].NickName,
					VideoTime: barrageInfo.VideoTime,
					Message:   barrageInfo.Message,
				}
			}
			if len(viewInfos) != 0 {
				err = c.WebsocketConn.Emit("clientNewMsg", viewInfos)
				if err != nil {
					fmt.Println(err)
				}
			}

			time.Sleep(5 * time.Second)
		}
	}(roomId, CurUid, c.WebsocketConn.Server())
}
func (c *ChatController) OnBarrageWebsocketDisconnect(roomId string) {
	CurUid := c.Ctx.Values().Get("CurUid").(int64)
	fmt.Println(CurUid, "disconnect", roomId)
}
func (c *ChatController) OnBarrageWebsocketGetNewMessage(message string) {
	CurUid := c.Ctx.Values().Get("CurUid").(int64)

	type Barrage struct {
		RoomId    string `json:"roomNO"`
		Message   string `json:"message"`
		videoTime int64  `json:"videoTime"`
	}
	var newMessage Barrage
	err := json.Unmarshal([]byte(message), &newMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	roomId, _ := strconv.ParseInt(newMessage.RoomId, 10, 64)
	err = c.ChatSrv.SendBarrage(CurUid, roomId, newMessage.Message, newMessage.videoTime)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *ChatController) GetBarrageWebsocketBy(roomId int64) {
	c.WebsocketConn.OnLeave(c.OnBarrageWebsocketDisconnect)

	//mvc模式没找到OnConnect钩子函数，只能通过自定义事件，建立连接之后前端再次发送connect来响应
	c.WebsocketConn.On("connected", c.OnBarrageWebsocketConnect)
	c.WebsocketConn.On("newMsg", c.OnBarrageWebsocketGetNewMessage)

	//call it after all event callbacks registration.
	c.WebsocketConn.Wait()
}

//注册房间
func (c *ChatController) PostRoomRegister() interface{} {
	var (
		roomName       = c.Ctx.FormValue("roomName")
		roomIdStr      = c.Ctx.FormValue("roomId")
		redisSetCntStr = c.Ctx.FormValue("redisSetCnt")
	)
	curUid := GetCurrentUserID(c.Session)
	if curUid == 0 {
		return ResParams{
			Code: 1001,
			Msg:  "please login first",
			Data: nil,
		}
	}
	roomId, _ := strconv.ParseInt(roomIdStr, 10, 64)
	redisSetCnt, _ := strconv.Atoi(redisSetCntStr)

	err := c.ChatSrv.AddRoom(roomName, roomId, curUid, redisSetCnt)
	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	return ResParams{
		Code: 1000,
		Msg:  "注册成功",
		Data: nil,
	}
}

//房间开播(通知后台开启kafka和redis资源的处理协程)
func (c *ChatController) PostRoomStart() interface{} {
	curUid := GetCurrentUserID(c.Session)
	if curUid == 0 {
		return ResParams{
			Code: 1001,
			Msg:  "未登陆",
			Data: nil,
		}
	}
	err := c.ChatSrv.StartRoom(curUid)
	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  "直播间打开失败",
			Data: nil,
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "启动成功",
		Data: nil,
	}
}

//房间停播(通知后台释放相关kafka和redis资源)
func (c *ChatController) PostRoomStop() interface{} {
	curUid := GetCurrentUserID(c.Session)
	if curUid == 0 {
		return ResParams{
			Code: 1001,
			Msg:  "未登陆",
			Data: nil,
		}
	}
	err := c.ChatSrv.StopRoom(curUid)
	if err != nil {
		return ResParams{
			Code: 1001,
			Msg:  "直播间关闭失败",
			Data: nil,
		}
	}
	return ResParams{
		Code: 1000,
		Msg:  "关闭成功",
		Data: nil,
	}
}

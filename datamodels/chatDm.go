package datamodels

import (
	"time"
)

const (
	//房间是否有效状态类型--Status
	RoomStatusNormal int8 = 1
	RoomStatusDelete int8 = 2

	//房间是否开播状态类型--WorkStatus
	WorkStatusInvalid  int8 = 0 //未开播
	WorkStatusStarting int8 = 1 //启动中
	WorkStatusOn       int8 = 2 //直播中
	WorkStatusStoping  int8 = 3 //关闭中
)

type ChatRoomConfigure struct {
	Id           int64
	RoomId       int64
	RoomName     string
	RoomOwnerUid int64
	RedisSetCnt  int
	Status       int8
	WorkStatus   int8
	CreateTime   time.Time `xorm:"created"`
	UpdateTime   time.Time `xorm:"updated"`
}

type BarrageInfo struct {
	Uid        int64
	message    string
	CreateTime time.Time
	VideoTime  time.Time
}

package datamodels

import "time"

//follow_info&fans_info&follow_group表中状态位数值定义
const (
	//Status 状态类型
	FollowStatusNormal int8 = 1
	FollowStatusDelete int8 = 0

	//IsFriend 状态类型
	IsFriendFalse int8 = 0
	IsFriendTrue  int8 = 1

	//follow_group 状态类型
	GroupStatusNormal int8 = 1
	GroupStatusDelete int8 = 0
)

//uidA 与 uidB 的关系
var (
	RelationSelf       int8 = 0 //本尊
	RelationNone       int8 = 1 //没关系
	RelationFollow     int8 = 2 //A关注B
	RelationFan        int8 = 3 //B关注A
	RelationCorelation int8 = 4 //AB互关
)

type FollowInfo struct {
	Id         int32
	Uid        int64
	FollowUid  int64
	Status     int8
	IsFriend   int8
	GroupId    int64
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}

type FanInfo struct {
	Id         int32
	Uid        int64
	FanUid     int64
	Status     int8
	IsFriend   int8
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}

type FollowGroup struct {
	Id         int64
	Uid        int64
	GroupName  string
	Status     int8
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}

type FollowMsg struct {
	MsgId     int64
	Uid       int64
	FollowUid int64
	GroupId   int64
	Status    int8
}

type GroupMsg struct {
	MsgId     int64
	Uid       int64
	FollowUid int64
	GroupId   int64
	InOrOut   bool
}

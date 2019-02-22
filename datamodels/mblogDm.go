package datamodels

import "time"

//mblog_info&uid_to_mblog表中状态位数值定义
const (
	//Status 状态类型
	MblogStatusNormal int8 = 1
	MblogStatusDelete int8 = 2
	MblogStatusShield int8 = 3

	//ReadAble 状态类型
	MblogReadAblePublic int8 = 1
	MblogReadAbleFriend int8 = 2
	MblogReadAblePerson int8 = 3
)

type MblogInfo struct {
	Id         int32
	Mid        int64
	Uid        int64
	Content    string
	OriginMid  int64
	OriginUid  int64
	TransCnt   int32
	LikesCnt   int32
	CommentCnt int32
	Status     int8
	ReadAble   int8
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}

type UidToMblog struct {
	Id         int32
	Uid        int64
	Mid        int64
	Status     int8
	ReadAble   int8
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}

type MblogSendKafkaStruct struct {
	Uid       int64
	FollowUid int64
	Status    int8
}
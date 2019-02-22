package datamodels

import "time"

type MblogComment struct{
	Id int32
	Mid int64
	Uid int64
	CommentUid int64
	content string
	LikesCnt int32
	Status int8
	CreateTime   time.Time `xorm:"created"`
	UpdateTime   time.Time `xorm:"updated"`
}

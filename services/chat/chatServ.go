package chatService

import (
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/services/mblogGrpc"
	"github.com/bluesky1024/goMblog/services/userGrpc"
)

type ChatServicer interface {
	//room相关
	EnterRoom(connectId int64, uid int64, roomId int64) error
	LeaveRoom(connectId int64, uid int64, roomId int64) error

	//message相关
	PostNewMessageIntoRoom(connectId int64, uid int64, roomId int64) error
	GetMessageFromRoom(connectId int64, uid int64, roomId int64)

	ReleaseSrv() error
}

type chatService struct {
	//弹幕消息队列
	kafkaProducer sarama.AsyncProducer

	userSrv  userGrpc.UserServicer
	mblogSrv mblogGrpc.MblogServicer
}

func NewChatServicer() (ChatServicer, error) {

	return &chatService{}, nil
}

func (s *chatService) ReleaseSrv() (err error) {
	return nil
}

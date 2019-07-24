package chatService

import (
	"errors"
	"github.com/Shopify/sarama"
	ds "github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/chat"
	"github.com/bluesky1024/goMblog/services/userGrpc"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "chatService"

type ChatServicer interface {
	//后台管理
	AddRoom(roomName string, roomId int64, roomOwnerUid int64, redisSetCnt int) error
	RemoveRoom(roomId int64) error

	//主播操作
	StartRoom(uid int64) error
	StopRoom(uid int64) error

	////观众操作
	//EnterRoom(connectId int64, uid int64, roomId int64) error
	//LeaveRoom(connectId int64, uid int64, roomId int64) error
	//PostNewMessageIntoRoom(connectId int64, uid int64, roomId int64) error
	//GetMessageFromRoom(connectId int64, uid int64, roomId int64)

	ReleaseSrv() error
}

type chatService struct {
	//弹幕消息队列
	kafkaProducer sarama.AsyncProducer
	dbRepo        *chatDbRepo.ChatDbRepository

	userSrv userGrpc.UserServicer
}

func NewChatServicer() (ChatServicer, error) {

	//user服务mysql仓库初始化
	chatSourceM, err := ds.LoadChats(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	chatSourceR, err := ds.LoadChats(false)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	chatRepo := chatDbRepo.NewChatRepository(chatSourceM, chatSourceR)

	//kafka生产者初始化
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")
	kafkaConfig := make(map[string]string)
	kafkaConfig["host"] = "0.0.0.0"
	kafkaConfig["port"] = "9092"
	kafkaProducer, err := newKafkaProducer([]string{kafkaConfig["host"] + ":" + kafkaConfig["port"]})
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}

	//userGrpc服务
	userSrv := userGrpc.NewUserGrpcServicer()
	if userSrv == nil {
		return nil, errors.New("user grpc server invalid")
	}

	return &chatService{
		dbRepo:        chatRepo,
		kafkaProducer: kafkaProducer,
		userSrv:       userSrv,
	}, nil
}

func newKafkaProducer(addr []string) (producer sarama.AsyncProducer, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_2_0_0
	producer, err = sarama.NewAsyncProducer(addr, config)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	return producer, nil
}

func (s *chatService) ReleaseSrv() (err error) {
	if s.kafkaProducer != nil {
		err = s.kafkaProducer.Close()
	}
	return err
}

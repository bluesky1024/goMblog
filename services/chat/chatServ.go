package chatService

import (
	"errors"
	"github.com/Shopify/sarama"
	dm "github.com/bluesky1024/goMblog/datamodels"
	ds "github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/chat"
	"github.com/bluesky1024/goMblog/repositories/redisRepo/chat"
	"github.com/bluesky1024/goMblog/services/userGrpc"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "chatService"

type ChatServicer interface {
	//获取房间配置信息/判断房间是否存在
	GetRoomConfigByRoomId(roomId int64) (info dm.ChatRoomConfigure, err error)

	//后台管理
	AddRoom(roomName string, roomId int64, roomOwnerUid int64, redisSetCnt int) error
	RemoveRoom(roomId int64) error

	//主播操作
	StartRoom(uid int64) error
	StopRoom(uid int64) error

	////观众操作
	GetBarrageByRoomId(uid int64, roomId int64) (barrages []dm.ChatBarrageInfo, err error)
	SendBarrage(uid int64, roomId int64, message string, videoTime int64) error

	//队列操作
	HandleRoomStartMsg(msg dm.RoomStatusSwitchMsg) (err error)
	HandleRoomStopMsg(msg dm.RoomStatusSwitchMsg) (err error)
	HandleNewBarrageToRoomMsg(msg dm.ChatBarrageInfo) (err error)

	ReleaseSrv() error
}

type chatService struct {
	//弹幕消息队列
	kafkaProducer sarama.AsyncProducer
	dbRepo        *chatDbRepo.ChatDbRepository
	rdRepo        *chatRdRepo.ChatRbRepository

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

	//redis服务初始化
	rdSource, err := redisSource.LoadChatRdSour()
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	chatRd := chatRdRepo.NewChatRdRepo(rdSource)
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
		rdRepo:        chatRd,
		kafkaProducer: kafkaProducer,
		userSrv:       userSrv,
	}, nil
}

func (s *chatService) ReleaseSrv() (err error) {
	if s.kafkaProducer != nil {
		err = s.kafkaProducer.Close()
	}
	return err
}

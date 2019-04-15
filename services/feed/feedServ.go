package feedService

import (
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"github.com/bluesky1024/goMblog/repositories/redisRepo/feed"
	"github.com/bluesky1024/goMblog/services/mblogGrpc"
	"github.com/bluesky1024/goMblog/services/relationGrpc"
	"github.com/bluesky1024/goMblog/services/userGrpc"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "relationService"

type FeedServicer interface {
	GetFeedByUidAndGroupId(uid int64, groupId int64, page int, pageSize int) (mids []int64)

	/*kafka关注取关分组管理补充操作*/
	//HandleFollowMsg(msg dm.FollowKafkaStruct) (err error)
	//HandleUnFollowMsg(msg dm.FollowKafkaStruct) (err error)
	////HandleGroupAddUidMsg()
	////HandleGroupRemUidMsg()
	////HandleGroupDelMsg()
	//
	//ReleaseSrv() error
}

type feedService struct {
	feedRdRepo  *feedRdRepo.FeedRbRepository
	userSrv     userGrpc.UserServicer
	mblogSrv    mblogGrpc.MblogServicer
	relationSrv relationGrpc.RelationServicer
}

func NewFeedServicer() FeedServicer {
	feedRdSourM, err := redisSource.LoadFeedRdSour(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}
	feedRdSourS, err := redisSource.LoadFeedRdSour(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}
	feedRepo := feedRdRepo.NewFeedRdRepo(feedRdSourM, feedRdSourS)

	userSrv := userGrpc.NewUserGrpcServicer()
	mblogSrv := mblogGrpc.NewMblogServicer()

	return &feedService{
		feedRdRepo: feedRepo,
		userSrv:    userSrv,
		mblogSrv:   mblogSrv,
	}
}

func (f *feedService) GetFeedByUidAndGroupId(uid int64, groupId int64, page int, pageSize int) (mids []int64) {
	return f.feedRdRepo.GetFeeds(uid, groupId, page, pageSize)
}

//func (f *feedService) HandleMblogSendMsg(msg dm.MblogSendKafkaStruct) (err error) {
//
//}

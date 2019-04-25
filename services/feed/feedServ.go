package feedService

import (
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"github.com/bluesky1024/goMblog/repositories/redisRepo/feed"
	"github.com/bluesky1024/goMblog/services/mblogGrpc"
	"github.com/bluesky1024/goMblog/services/relationGrpc"
	"github.com/bluesky1024/goMblog/services/userGrpc"
	"github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
	"time"
)

var logType = "relationService"

type FeedServicer interface {
	//按page顺序获取feed信息
	GetFeedByUidAndGroupId(uid int64, groupId int64, page int, pageSize int) (mids []int64, err error)
	//获取在某条微博之前（更旧）的size条微博
	GetFeedMoreByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error)
	//获取在某条微博之后（更新）的size条微博
	GetFeedNewerByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error)

	/*kafka关注取关分组管理补充操作*/
	//HandleFollowMsg(msg dm.FollowMsg) (err error)
	//HandleUnFollowMsg(msg dm.FollowMsg) (err error)
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

func (f *feedService) GetFeedByUidAndGroupId(uid int64, groupId int64, page int, pageSize int) (mids []int64, err error) {
	return f.feedRdRepo.GetFeeds(uid, groupId, page, pageSize)
}

func (f *feedService) GetFeedMoreByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error) {
	var timeBefore int64
	if Mid == 0 {
		timeBefore = time.Now().UnixNano()
		timeBefore = timeBefore / 1e6
	} else {
		timeBefore = idGenerate.GetDetailTimeById(Mid)
	}
	return f.feedRdRepo.GetByTimeBefore(uid, groupId, timeBefore, size)
}

func (f *feedService) GetFeedNewerByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error) {
	timeAfter := idGenerate.GetDetailTimeById(Mid)
	return f.feedRdRepo.GetByTimeAfter(uid, groupId, timeAfter, size)
}

//func (f *feedService) HandleMblogSendMsg(msg dm.MblogSendKafkaStruct) (err error) {
//
//}

package feedService

//个人feed中保存100天内的mid数据
//每天脚本淘汰过期feed数据

import (
	"errors"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"github.com/bluesky1024/goMblog/repositories/redisRepo/feed"
	"github.com/bluesky1024/goMblog/services/mblogGrpc"
	"github.com/bluesky1024/goMblog/services/relationGrpc"
	"github.com/bluesky1024/goMblog/services/userGrpc"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "feedService"

type FeedServicer interface {
	//按page顺序获取feed信息
	GetFeedByUidAndGroupId(uid int64, groupId int64, page int, pageSize int) (mids []int64, err error)
	//获取在某条微博之前（更旧）的size条微博
	GetFeedMoreByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error)
	//获取在某条微博之后（更新）的size条微博
	GetFeedNewerByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error)

	/*kafka消息处理*/
	HandleFollowMsg(msg dm.FollowMsg) (err error)
	HandleUnFollowMsg(msg dm.FollowMsg) (err error)
	HandleMblogNewMsg(msg dm.MblogNewMsg) (err error)
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

func NewFeedServicer() (s FeedServicer, err error) {
	feedRdSourM, err := redisSource.LoadFeedRdSour(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	feedRdSourS, err := redisSource.LoadFeedRdSour(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	feedRepo := feedRdRepo.NewFeedRdRepo(feedRdSourM, feedRdSourS)

	userSrv := userGrpc.NewUserGrpcServicer()
	if userSrv == nil {
		return nil, errors.New("user grpc server invalid")
	}
	mblogSrv := mblogGrpc.NewMblogServicer()
	if mblogSrv == nil {
		return nil, errors.New("mblog grpc server invalid")
	}

	return &feedService{
		feedRdRepo: feedRepo,
		userSrv:    userSrv,
		mblogSrv:   mblogSrv,
	}, nil
}

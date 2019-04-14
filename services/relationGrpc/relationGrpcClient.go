package relationGrpc

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	pb "github.com/bluesky1024/goMblog/services/relationGrpc/relationProto"
	"github.com/bluesky1024/goMblog/tools/logger"
	"google.golang.org/grpc"
)

var logType = "relationGrpcClient"

type RelationServicer interface {
	GetFollowsByUid(uid int64, page int, pageSize int) (follows []dm.FollowInfo, cnt int64)
	GetFansByUid(uid int64, page int, pageSize int) (fans []dm.FanInfo, cnt int64)

	/*分组信息*/
	GetGroupsByUid(uid int64) (groups []dm.FollowGroup, cnt int64)
}

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50054"
)

func NewRelationGrpcServicer() RelationServicer {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}

	c := pb.NewRelationClient(conn)

	return &relationService{
		client: c,
	}
}

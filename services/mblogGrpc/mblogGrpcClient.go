package mblogGrpc

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	pb "github.com/bluesky1024/goMblog/services/mblogGrpc/mblogProto"
	"github.com/bluesky1024/goMblog/tools/logger"
	"google.golang.org/grpc"
)

var logType = "mblogGrpcClient"

type MblogServicer interface {
	Create(uid int64, content string, readAble int8, originUid int64, originMid int64) (mblog dm.MblogInfo, err error)
	GetMultiByMids(mid []int64) (mblogs map[int64]dm.MblogInfo)
	GetNormalByUid(uid int64, page int, pageSize int, readAble []int8, startTime int64, endTime int64) (mblogs []dm.MblogInfo, cnt int64)
}

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50053"
)

func NewMblogServicer() MblogServicer {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}

	c := pb.NewMblogClient(conn)

	return &mblogService{
		client: c,
	}
}

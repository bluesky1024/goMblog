package userGrpc

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	pb "github.com/bluesky1024/goMblog/services/userGrpc/userProto"
	"github.com/bluesky1024/goMblog/tools/logger"
	"google.golang.org/grpc"
)

var logType = "userGrpcClient"

type UserServicer interface {
	Create(nickname string, password string, telephone string, email string) (user dm.User, err error)
	GetByNicknameAndPassword(nickname string, password string) (user dm.User, err error)
	GetByUid(uid int64) (user dm.User, err error)
	GetMultiByUids(uids []int64) (users []dm.User, err error)
}

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func NewUserGrpcServicer() UserServicer {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))

	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}

	c := pb.NewMblogUserClient(conn)

	return &userService{
		client: c,
	}
}

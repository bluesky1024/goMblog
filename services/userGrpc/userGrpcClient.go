package userGrpc

import (
	pb "github.com/bluesky1024/goMblog/grpcServer/user/userProto"
	"github.com/bluesky1024/goMblog/tools/logger"
	"google.golang.org/grpc"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main(){
	var opts []grpc.DialOption

	opts = append(opts,grpc.WithInsecure())

	opts = append(opts,grpc.WithPerRPCCredentials(new(customCredential)))

	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		logger.Err("userGrpcClient",err.Error())
	}

	c := pb.NewMblogUserClient(conn)
}
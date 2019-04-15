package mblogGrpc

import (
	"context"
	dm "github.com/bluesky1024/goMblog/datamodels"
	pb "github.com/bluesky1024/goMblog/services/mblogGrpc/mblogProto"
	"google.golang.org/grpc"
)

type relationService struct {
	client pb.MblogClient
}

func (s *mblogService) Create(uid int64, content string, readAble int8, originUid int64, originMid int64) (mblog dm.MblogInfo, err error) {

	return mblog, err
}

func (s *mblogService) GetByMid(mid int64) (mblog dm.MblogInfo, found bool) {

	return mblog, found
}

func (s *mblogService) GetByUid(uid int64, page int, pageSize int) (mblogs map[int64]dm.MblogInfo) {

	return mblogs
}

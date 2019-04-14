package main

import (
	"context"
	pb "github.com/bluesky1024/goMblog/grpcServer/relation/relationProto"
	"github.com/bluesky1024/goMblog/services/relation"
)

func NewRelationService() *RelationService {
	return &RelationService{
		serv: relationService.NewRelationServicer(),
	}
}

type RelationService struct {
	serv relationService.RelationServicer
}

func (s *RelationService) GetFollowsByUid(ctx context.Context, req *pb.UidReq) (res *pb.MultiFollows, err error) {
	return res, err
}

func (s *RelationService) GetFansByUid(ctx context.Context, req *pb.UidReq) (res *pb.MultiFans, err error) {
	return res, err
}

func (s *RelationService) GetGroupsByUid(ctx context.Context, req *pb.UidReq) (res *pb.MultiFollowGroups, err error) {
	return res, err
}
